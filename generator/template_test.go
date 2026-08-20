package generator

import (
	"strings"
	"testing"
)

func TestLoadTemplateBuiltin(t *testing.T) {
	types := []string{"domain", "repository", "service", "controller", "task"}
	for _, typeName := range types {
		t.Run(typeName, func(t *testing.T) {
			content, err := loadTemplate(typeName)
			if err != nil {
				t.Fatalf("loadTemplate(%q) failed: %v", typeName, err)
			}
			if content == "" {
				t.Fatalf("template content for %q is empty", typeName)
			}
			if !strings.Contains(content, "package "+typeName) {
				t.Errorf("template for %q should contain 'package %s'", typeName, typeName)
			}
		})
	}
}

func TestLoadTemplateUnknownType(t *testing.T) {
	_, err := loadTemplate("bogus")
	if err == nil || !strings.Contains(err.Error(), "unknown template type") {
		t.Fatalf("expected 'unknown template type' error, got: %v", err)
	}
}

