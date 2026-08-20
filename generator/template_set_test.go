package generator

import (
	"os"
	"path/filepath"
	"strings"
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

func TestLoadTemplateUnknownType(t *testing.T) {
	cmd := &cobra.Command{}
	_, err := loadTemplate(cmd, "default", "bogus")
	if err == nil || !strings.Contains(err.Error(), "unknown template type") {
		t.Fatalf("expected 'unknown template type' error, got: %v", err)
	}
}

func TestLoadTemplateEmptySetName(t *testing.T) {
	// templateSet 为空字符串时也应允许使用内置模板
	cmd := &cobra.Command{}
	content, err := loadTemplate(cmd, "", "domain")
	if err != nil {
		t.Fatalf("loadTemplate with empty set name: %v", err)
	}
	if !contains(content, "package domain") {
		t.Errorf("expected builtin template, got: %s", content)
	}
}

func TestTemplateSetDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	dir, err := templateSetDir()
	if err != nil {
		t.Fatalf("templateSetDir: %v", err)
	}
	want := filepath.Join(homeDir, ".gouno", "templates")
	if dir != want {
		t.Errorf("templateSetDir = %q; want %q", dir, want)
	}
}

func TestLoadProjectConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		chdirToTemp(t)
		if err := os.WriteFile(configFileName, []byte("template-set: my-set\n"), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := loadProjectConfig()
		if cfg == nil || cfg.TemplateSet != "my-set" {
			t.Fatalf("expected template-set my-set, got %+v", cfg)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		chdirToTemp(t)
		if cfg := loadProjectConfig(); cfg != nil {
			t.Fatalf("expected nil config without .gouno.yaml, got %+v", cfg)
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		chdirToTemp(t)
		if err := os.WriteFile(configFileName, []byte("{invalid: [\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if cfg := loadProjectConfig(); cfg != nil {
			t.Fatalf("expected nil config for invalid yaml, got %+v", cfg)
		}
	})

	t.Run("empty template-set", func(t *testing.T) {
		chdirToTemp(t)
		if err := os.WriteFile(configFileName, []byte("template-set: \"\"\n"), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := loadProjectConfig()
		if cfg == nil {
			t.Fatal("expected config to be parsed")
		}
		if cfg.TemplateSet != "" {
			t.Errorf("expected empty template-set, got %q", cfg.TemplateSet)
		}
	})
}

func TestResolveTemplateSetFromConfig(t *testing.T) {
	chdirToTemp(t)
	if err := os.WriteFile(configFileName, []byte("template-set: yaml-set\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := &cobra.Command{}
	cmd.Flags().String("template-set", "", "")
	if got := resolveTemplateSet(cmd); got != "yaml-set" {
		t.Errorf("resolveTemplateSet = %q; want yaml-set", got)
	}
}

func TestResolveTemplateSetFlagOverridesConfig(t *testing.T) {
	chdirToTemp(t)
	if err := os.WriteFile(configFileName, []byte("template-set: yaml-set\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := &cobra.Command{}
	cmd.Flags().String("template-set", "", "")
	if err := cmd.Flags().Set("template-set", "flag-set"); err != nil {
		t.Fatal(err)
	}
	if got := resolveTemplateSet(cmd); got != "flag-set" {
		t.Errorf("resolveTemplateSet = %q; want flag-set (flag should win)", got)
	}
}

func chdirToTemp(t *testing.T) string {
	t.Helper()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restore wd: %v", err)
		}
	})
	return tmpDir
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
