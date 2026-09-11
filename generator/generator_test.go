package generator_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rushairer/gouno/generator"
	"github.com/spf13/cobra"
)

const testManifest = `schema: gouno.dev/codegen/v1
command:
  use: gen
  short: Generate project code
  aliases: [generator]
generators:
  - name: domain
    aliases: [d]
    short: Generate domain
    args:
      - name: name
        required: true
    flags:
      - name: path
        shorthand: p
        type: string
        default: internal/domain
        description: path to domain
      - name: force
        shorthand: f
        type: bool
        default: false
        description: force overwrite
    outputs:
      - template: .gouno/codegen/domain.tmpl
        path: '{{ flag "path" }}/{{ arg "name" }}.go'
  - name: service
    aliases: [s]
    short: Generate service
    args:
      - name: name
        required: true
    flags:
      - name: path
        shorthand: p
        type: string
        default: internal/service
        description: path to service
      - name: force
        shorthand: f
        type: bool
        default: false
        description: force overwrite
    outputs:
      - template: .gouno/codegen/service.tmpl
        path: '{{ flag "path" }}/{{ arg "name" }}.go'
  - name: suite
    short: Generate domain and service
    args:
      - name: name
        required: true
    flags:
      - name: force
        shorthand: f
        type: bool
        default: false
        description: force overwrite
    compose: [domain, service]
`

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".gouno", "codegen.yaml"), testManifest)
	writeFile(t, filepath.Join(root, ".gouno", "codegen", "domain.tmpl"), "package domain\n\ntype {{ camel (arg \"name\") }} struct{}\n")
	writeFile(t, filepath.Join(root, ".gouno", "codegen", "service.tmpl"), "package service\n\ntype {{ camel (arg \"name\") }}Service struct{}\n")
	return root
}

func execute(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	return output.String(), err
}

func TestAttachProjectCommandIsCapabilityDriven(t *testing.T) {
	root := &cobra.Command{Use: "app"}
	attached, err := generator.AttachProjectCommand(root, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if attached || len(root.Commands()) != 0 {
		t.Fatal("codegen command must be absent without a manifest")
	}

	root = &cobra.Command{Use: "app"}
	attached, err = generator.AttachProjectCommand(root, writeProject(t))
	if err != nil {
		t.Fatal(err)
	}
	if !attached {
		t.Fatal("expected template codegen capability to attach")
	}
	cmd, _, err := root.Find([]string{"gen"})
	if err != nil || cmd == root {
		t.Fatalf("expected template-defined gen command, cmd=%v err=%v", cmd, err)
	}
}

func TestTemplateDefinedGeneratorRendersAndFormats(t *testing.T) {
	project := writeProject(t)
	cmd, err := generator.LoadProjectCommand(project)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := execute(t, cmd, "service", "foo_bar"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(project, "internal", "service", "foo_bar.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "type FooBarService struct{}") {
		t.Fatalf("unexpected generated content:\n%s", content)
	}
}

func TestTemplateDefinedCompositionAndForce(t *testing.T) {
	project := writeProject(t)
	cmd, err := generator.LoadProjectCommand(project)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := execute(t, cmd, "suite", "account"); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"internal/domain/account.go", "internal/service/account.go"} {
		if _, err := os.Stat(filepath.Join(project, rel)); err != nil {
			t.Fatalf("expected %s: %v", rel, err)
		}
	}

	service := filepath.Join(project, "internal", "service", "account.go")
	writeFile(t, service, "sentinel")
	cmd, _ = generator.LoadProjectCommand(project)
	if _, err := execute(t, cmd, "suite", "account", "--force"); err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(service)
	if strings.Contains(string(content), "sentinel") {
		t.Fatal("compose generator did not propagate force")
	}
}

func TestManifestSearchesParents(t *testing.T) {
	project := writeProject(t)
	nested := filepath.Join(project, "internal", "service")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	path, root, err := generator.FindManifest(nested)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(project, ".gouno", "codegen.yaml") || root != project {
		t.Fatalf("unexpected manifest lookup: path=%s root=%s", path, root)
	}
}

func TestManifestRejectsCompositionCycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "codegen.yaml")
	writeFile(t, path, `schema: gouno.dev/codegen/v1
command:
  use: gen
generators:
  - name: a
    compose: [b]
  - name: b
    compose: [a]
`)
	_, err := generator.LoadManifest(path)
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("expected cycle error, got %v", err)
	}
}

func TestOutputTraversalIsRejected(t *testing.T) {
	project := writeProject(t)
	manifest := strings.Replace(testManifest, "internal/service", "../outside", 1)
	writeFile(t, filepath.Join(project, ".gouno", "codegen.yaml"), manifest)
	cmd, err := generator.LoadProjectCommand(project)
	if err != nil {
		t.Fatal(err)
	}
	_, err = execute(t, cmd, "service", "escape")
	if err == nil || !strings.Contains(err.Error(), "outside the project root") {
		t.Fatalf("expected traversal error, got %v", err)
	}
}
