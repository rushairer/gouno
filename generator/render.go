package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rushairer/gouno/utility"
)

type templateData struct {
	Args  map[string]string
	Flags map[string]any
}

func (r *runner) renderOutput(spec *GeneratorSpec, output OutputSpec, inv invocation, out io.Writer) error {
	funcs := templateFuncs(inv)
	targetRelative, err := renderString("output path", output.Path, templateData{Args: inv.args, Flags: inv.flags}, funcs)
	if err != nil {
		return fmt.Errorf("generator %q: %w", spec.Name, err)
	}
	targetPath, err := safeProjectPath(r.projectRoot, targetRelative)
	if err != nil {
		return fmt.Errorf("generator %q output path: %w", spec.Name, err)
	}

	templatePath, err := safeProjectPath(r.projectRoot, output.Template)
	if err != nil {
		return fmt.Errorf("generator %q template path: %w", spec.Name, err)
	}
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("generator %q read template %q: %w", spec.Name, output.Template, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Option("missingkey=error").Funcs(funcs).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("generator %q parse template %q: %w", spec.Name, output.Template, err)
	}
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, templateData{Args: inv.args, Flags: inv.flags}); err != nil {
		return fmt.Errorf("generator %q render template %q: %w", spec.Name, output.Template, err)
	}

	content := rendered.Bytes()
	if filepath.Ext(targetPath) == ".go" {
		content, err = format.Source(content)
		if err != nil {
			return fmt.Errorf("generator %q format Go output %q: %w", spec.Name, targetRelative, err)
		}
	}

	force, _ := inv.flags["force"].(bool)
	if !force {
		if _, err := os.Stat(targetPath); err == nil {
			_, _ = fmt.Fprintf(out, "%s file already exists, skipping: %s (use --force to overwrite)\n", spec.Name, targetPath)
			return nil
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("generator %q stat output %q: %w", spec.Name, targetRelative, err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("generator %q create output directory: %w", spec.Name, err)
	}
	if err := os.WriteFile(targetPath, content, 0o644); err != nil {
		return fmt.Errorf("generator %q write output %q: %w", spec.Name, targetRelative, err)
	}
	_, _ = fmt.Fprintf(out, "Created %s file: %s\n", spec.Name, targetPath)
	return nil
}

func renderString(name, source string, data templateData, funcs template.FuncMap) (string, error) {
	tmpl, err := template.New(name).Option("missingkey=error").Funcs(funcs).Parse(source)
	if err != nil {
		return "", fmt.Errorf("parse %s template: %w", name, err)
	}
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("render %s template: %w", name, err)
	}
	return rendered.String(), nil
}

func templateFuncs(inv invocation) template.FuncMap {
	return template.FuncMap{
		"arg":   func(name string) string { return inv.args[name] },
		"flag":  func(name string) string { return stringify(inv.flags[name]) },
		"camel": utility.ToCamelCase,
		"snake": utility.ToSnakeCase,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"kebab": func(value string) string { return strings.ReplaceAll(utility.ToSnakeCase(value), "_", "-") },
	}
}

func safeProjectPath(projectRoot, relative string) (string, error) {
	if relative == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if filepath.IsAbs(relative) {
		return "", fmt.Errorf("absolute path %q is not allowed", relative)
	}
	cleaned := filepath.Clean(relative)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q resolves outside the project root", relative)
	}
	resolved := filepath.Join(projectRoot, cleaned)
	rel, err := filepath.Rel(projectRoot, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q resolves outside the project root", relative)
	}
	return resolved, nil
}
