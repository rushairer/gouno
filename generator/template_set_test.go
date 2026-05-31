package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestLoadTemplateBuiltin(t *testing.T) {
	cmd := &cobra.Command{}
	content, err := loadTemplate(cmd, "default", "domain")
	if err != nil {
		t.Fatalf("loadTemplate failed: %v", err)
	}
	if content == "" {
		t.Fatal("template content is empty")
	}
	if !contains(content, "package domain") {
		t.Errorf("template should contain 'package domain'")
	}
}

func TestLoadTemplateLocal(t *testing.T) {
	// 创建临时模板目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	tmplDir := filepath.Join(homeDir, ".gouno", "templates", "test-local")
	os.MkdirAll(tmplDir, 0755)
	defer os.RemoveAll(tmplDir)

	customTmpl := `package domain

type %s struct {
	CustomField string
}
`
	os.WriteFile(filepath.Join(tmplDir, "domain.tmpl"), []byte(customTmpl), 0644)

	cmd := &cobra.Command{}
	content, err := loadTemplate(cmd, "test-local", "domain")
	if err != nil {
		t.Fatalf("loadTemplate failed: %v", err)
	}
	if !contains(content, "CustomField") {
		t.Errorf("should use local template, got: %s", content)
	}
}

func TestLoadTemplateNotFound(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := loadTemplate(cmd, "nonexistent", "domain")
	if err == nil {
		t.Fatal("expected error for nonexistent template set")
	}
}

func TestResolveTemplateSet(t *testing.T) {
	t.Run("from flag", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("template-set", "", "")
		cmd.Flags().Set("template-set", "my-set")
		got := resolveTemplateSet(cmd)
		if got != "my-set" {
			t.Errorf("resolveTemplateSet = %q; want my-set", got)
		}
	})

	t.Run("default fallback", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("template-set", "", "")
		got := resolveTemplateSet(cmd)
		if got != "default" {
			t.Errorf("resolveTemplateSet = %q; want default", got)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsCheck(s, substr))
}

func containsCheck(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
