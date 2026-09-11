package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"go.yaml.in/yaml/v3"
)

const (
	ManifestSchema      = "gouno.dev/codegen/v1"
	DefaultManifestPath = ".gouno/codegen.yaml"
)

var (
	ErrManifestNotFound = errors.New("gouno codegen manifest not found")
	commandNamePattern  = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
)

type Manifest struct {
	Schema     string          `yaml:"schema"`
	Command    CommandSpec     `yaml:"command"`
	Generators []GeneratorSpec `yaml:"generators"`
}

type CommandSpec struct {
	Use     string   `yaml:"use"`
	Short   string   `yaml:"short"`
	Aliases []string `yaml:"aliases"`
}

type GeneratorSpec struct {
	Name    string         `yaml:"name"`
	Short   string         `yaml:"short"`
	Aliases []string       `yaml:"aliases"`
	Args    []ArgumentSpec `yaml:"args"`
	Flags   []FlagSpec     `yaml:"flags"`
	Outputs []OutputSpec   `yaml:"outputs"`
	Compose []string       `yaml:"compose"`
}

type ArgumentSpec struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
}

type FlagSpec struct {
	Name        string `yaml:"name"`
	Shorthand   string `yaml:"shorthand"`
	Type        string `yaml:"type"`
	Default     any    `yaml:"default"`
	Description string `yaml:"description"`
}

type OutputSpec struct {
	Template string `yaml:"template"`
	Path     string `yaml:"path"`
}

func FindManifest(startDir string) (manifestPath string, projectRoot string, err error) {
	if startDir == "" {
		startDir, err = os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("get working directory: %w", err)
		}
	}

	current, err := filepath.Abs(startDir)
	if err != nil {
		return "", "", fmt.Errorf("resolve start directory: %w", err)
	}

	info, err := os.Stat(current)
	if err != nil {
		return "", "", fmt.Errorf("stat start directory: %w", err)
	}
	if !info.IsDir() {
		current = filepath.Dir(current)
	}

	for {
		candidate := filepath.Join(current, DefaultManifestPath)
		info, statErr := os.Stat(candidate)
		switch {
		case statErr == nil && !info.IsDir():
			return candidate, current, nil
		case statErr == nil:
			return "", "", fmt.Errorf("codegen manifest path is a directory: %s", candidate)
		case !errors.Is(statErr, os.ErrNotExist):
			return "", "", fmt.Errorf("stat codegen manifest: %w", statErr)
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", "", ErrManifestNotFound
		}
		current = parent
	}
}

func LoadManifest(path string) (*Manifest, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read codegen manifest: %w", err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("parse codegen manifest: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (m *Manifest) Validate() error {
	if m.Schema != ManifestSchema {
		return fmt.Errorf("unsupported codegen schema %q, want %q", m.Schema, ManifestSchema)
	}
	if err := validateCommandName("command.use", m.Command.Use); err != nil {
		return err
	}
	if len(m.Generators) == 0 {
		return fmt.Errorf("codegen manifest must define at least one generator")
	}

	generators := make(map[string]*GeneratorSpec, len(m.Generators))
	for i := range m.Generators {
		spec := &m.Generators[i]
		if err := validateCommandName("generator.name", spec.Name); err != nil {
			return err
		}
		if _, exists := generators[spec.Name]; exists {
			return fmt.Errorf("duplicate generator %q", spec.Name)
		}
		generators[spec.Name] = spec

		seenArgs := map[string]struct{}{}
		optionalSeen := false
		for _, arg := range spec.Args {
			if err := validateCommandName("argument.name", arg.Name); err != nil {
				return fmt.Errorf("generator %q: %w", spec.Name, err)
			}
			if _, exists := seenArgs[arg.Name]; exists {
				return fmt.Errorf("generator %q has duplicate argument %q", spec.Name, arg.Name)
			}
			seenArgs[arg.Name] = struct{}{}
			if !arg.Required {
				optionalSeen = true
			} else if optionalSeen {
				return fmt.Errorf("generator %q has required argument %q after an optional argument", spec.Name, arg.Name)
			}
		}

		seenFlags := map[string]struct{}{}
		seenShorthands := map[string]struct{}{}
		for _, flag := range spec.Flags {
			if err := validateCommandName("flag.name", flag.Name); err != nil {
				return fmt.Errorf("generator %q: %w", spec.Name, err)
			}
			if _, exists := seenFlags[flag.Name]; exists {
				return fmt.Errorf("generator %q has duplicate flag %q", spec.Name, flag.Name)
			}
			seenFlags[flag.Name] = struct{}{}
			if flag.Shorthand != "" {
				if len(flag.Shorthand) != 1 {
					return fmt.Errorf("generator %q flag %q shorthand must be one character", spec.Name, flag.Name)
				}
				if _, exists := seenShorthands[flag.Shorthand]; exists {
					return fmt.Errorf("generator %q has duplicate shorthand %q", spec.Name, flag.Shorthand)
				}
				seenShorthands[flag.Shorthand] = struct{}{}
			}
			switch flag.Type {
			case "string":
				if flag.Default != nil {
					if _, ok := flag.Default.(string); !ok {
						return fmt.Errorf("generator %q flag %q default must be a string", spec.Name, flag.Name)
					}
				}
			case "bool":
				if flag.Default != nil {
					if _, ok := flag.Default.(bool); !ok {
						return fmt.Errorf("generator %q flag %q default must be a bool", spec.Name, flag.Name)
					}
				}
			default:
				return fmt.Errorf("generator %q flag %q has unsupported type %q", spec.Name, flag.Name, flag.Type)
			}
			if flag.Name == "force" && flag.Type != "bool" {
				return fmt.Errorf("generator %q reserved flag %q must be bool", spec.Name, flag.Name)
			}
		}

		if len(spec.Outputs) == 0 && len(spec.Compose) == 0 {
			return fmt.Errorf("generator %q must define outputs or compose", spec.Name)
		}
		for _, output := range spec.Outputs {
			if output.Template == "" || output.Path == "" {
				return fmt.Errorf("generator %q output requires template and path", spec.Name)
			}
		}
	}

	for _, spec := range m.Generators {
		parentArgs := map[string]struct{}{}
		for _, arg := range spec.Args {
			parentArgs[arg.Name] = struct{}{}
		}
		for _, childName := range spec.Compose {
			child, ok := generators[childName]
			if !ok {
				return fmt.Errorf("generator %q composes unknown generator %q", spec.Name, childName)
			}
			for _, arg := range child.Args {
				if _, ok := parentArgs[arg.Name]; !ok {
					return fmt.Errorf("generator %q cannot compose %q because argument %q is not available", spec.Name, childName, arg.Name)
				}
			}
		}
	}

	visiting := map[string]bool{}
	visited := map[string]bool{}
	var visit func(string) error
	visit = func(name string) error {
		if visiting[name] {
			return fmt.Errorf("codegen composition cycle detected at %q", name)
		}
		if visited[name] {
			return nil
		}
		visiting[name] = true
		for _, child := range generators[name].Compose {
			if err := visit(child); err != nil {
				return err
			}
		}
		visiting[name] = false
		visited[name] = true
		return nil
	}
	for name := range generators {
		if err := visit(name); err != nil {
			return err
		}
	}

	return nil
}

func validateCommandName(field, value string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", field)
	}
	if !commandNamePattern.MatchString(value) {
		return fmt.Errorf("%s %q must match %s", field, value, commandNamePattern.String())
	}
	return nil
}
