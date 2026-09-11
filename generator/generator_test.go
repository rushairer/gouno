package generator_test

import (
	"bytes"
	"errors"
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

func writeTestProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, ".gouno", "codegen.yaml"), testManifest)
	mustWrite(t, filepath.Join(root, ".gouno", "codegen", "domain.tmpl"), `package domain

type {{ camel (arg "name") }} struct{}
`)
	mustWrite(t, filepath.Join(root, ".gouno", "codegen", "service.tmpl"), `package service

type {{ camel (arg "name") }}Service struct{}
`)
	return root
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
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

func TestLoadProjectCommandAbsent(t *testing.T) {
	cmd, err := generator.LoadProjectCommand(t.TempDir())
	if err != nil {
		t.Fatalf("LoadProjectCommand returned error: %v", err)
	}
	if cmd != nil {
		t.Fatal("expected no codegen command when manifest is absent")
	}
}

func TestAttachProjectCommandShapesRootCLI(t *testing.T) {
	withoutGen := &cobra.Command{Use: "app"}
	attached, err := generator.AttachProjectCommand(withoutGen, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if attached || len(withoutGen.Commands()) != 0 {
		t.Fatal("gen command must not appear when template provides no manifest")
	}

	rootDir := writeTestProject(t)
	withGen := &cobra.Command{Use: "app"}
	attached, err = generator.AttachProjectCommand(withGen, rootDir)
	if err != nil {
		t.Fatal(err)
	}
	if !attached {
		t.Fatal("expected codegen capability to attach")
	}
	cmd, _, err := withGen.Find([]string{"gen"})
	if err != nil || cmd == withGen {
		t.Fatalf("expected template-defined gen command, cmd=%v err=%v", cmd, err)
	}
}

func TestTemplateDefinedGeneratorRendersAndFormats(t *testing.T) {
	rootDir := writeTestProject(t)
	cmd, err := generator.LoadProjectCommand(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	output, err := execute(t, cmd, "service", "foo_bar")
	if err != nil {
		t.Fatalf("execute generator: %v\n%s", err, output)
	}
	generated := filepath.Join(rootDir, "internal", "service", "foo_bar.go")
	content, err := os.ReadFile(generated)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "type FooBarService struct{}") {
		t.Fatalf("unexpected generated content:\n%s", content)
	}
}

func TestTemplateDefinedPathFlagAndOverwritePolicy(t *testing.T) {
	rootDir := writeTestProject(t)
	cmd, err := generator.LoadProjectCommand(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := execute(t, cmd, "service", "user", "--path", "custom/services"); err != nil {
		t.Fatal(err)
	}
	generated := filepath.Join(rootDir, "custom", "services", "user.go")
	mustWrite(t, generated, "sentinel")

	cmd, _ = generator.LoadProjectCommand(rootDir)
	output, err := execute(t, cmd, "service", "user", "--path", "custom/services")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "already exists, skipping") {
		t.Fatalf("expected skip message, got %q", output)
	}
	content, _ := os.ReadFile(generated)
	if string(content) != "sentinel" {
		t.Fatalf("existing file changed without force: %q", content)
	}

	cmd, _ = generator.LoadProjectCommand(rootDir)
	if _, err := execute(t, cmd, "service", "user", "--path", "custom/services", "--force"); err != nil {
		t.Fatal(err)
	}
	content, _ = os.ReadFile(generated)
	if strings.Contains(string(content), "sentinel") {
		t.Fatal("force did not overwrite existing file")
	}
}

func TestComposeGeneratorUsesTemplatePolicy(t *testing.T) {
	rootDir := writeTestProject(t)
	cmd, err := generator.LoadProjectCommand(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := execute(t, cmd, "suite", "account"); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join(rootDir, "internal", "domain", "account.go"),
		filepath.Join(rootDir, "internal", "service", "account.go"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected generated file %s: %v", path, err)
		}
	}
}

func TestFindManifestWalksUpFromNestedDirectory(t *testing.T) {
	rootDir := writeTestProject(t)
	nested := filepath.Join(rootDir, "internal", "service")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	path, foundRoot, err := generator.FindManifest(nested)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(rootDir, ".gouno", "codegen.yaml") || foundRoot != rootDir {
		t.Fatalf("unexpected manifest lookup: path=%s root=%s", path, foundRoot)
	}
}

func TestManifestValidationRejectsCompositionCycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "codegen.yaml")
	mustWrite(t, path, `schema: gouno.dev/codegen/v1
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
		t.Fatalf("expected cycle validation error, got %v", err)
	}
}

func TestOutputTraversalIsRejected(t *testing.T) {
	rootDir := writeTestProject(t)
	manifestPath := filepath.Join(rootDir, ".gouno", "codegen.yaml")
	mustWrite(t, manifestPath, strings.Replace(testManifest, "internal/service", "../outside", 1))
	cmd, err := generator.LoadProjectCommand(rootDir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = execute(t, cmd, "service", "escape")
	if err == nil || !strings.Contains(err.Error(), "outside the project root") {
		t.Fatalf("expected traversal error, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(filepath.Dir(rootDir), "outside", "escape.go")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("unexpected file outside project root: %v", statErr)
	}
}
