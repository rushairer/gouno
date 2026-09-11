package generator

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// GeneratorCmd is kept as a compatibility proxy for projects created by older
// gouno-template releases. New templates should call AttachProjectCommand so
// the root command only exposes codegen when the current template provides it.
var GeneratorCmd = &cobra.Command{
	Use:                "generator",
	Short:              "Generate code from the current project's template",
	Aliases:            []string{"gen"},
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectCmd, err := LoadProjectCommand("")
		if err != nil {
			return err
		}
		if projectCmd == nil {
			return ErrManifestNotFound
		}
		projectCmd.SetArgs(args)
		projectCmd.SetOut(cmd.OutOrStdout())
		projectCmd.SetErr(cmd.ErrOrStderr())
		return projectCmd.Execute()
	},
}

func LoadProjectCommand(startDir string) (*cobra.Command, error) {
	manifestPath, projectRoot, err := FindManifest(startDir)
	if err != nil {
		if errors.Is(err, ErrManifestNotFound) {
			return nil, nil
		}
		return nil, err
	}
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	return buildCommand(projectRoot, manifest), nil
}

func AttachProjectCommand(root *cobra.Command, startDir string) (bool, error) {
	cmd, err := LoadProjectCommand(startDir)
	if err != nil {
		return false, err
	}
	if cmd == nil {
		return false, nil
	}
	root.AddCommand(cmd)
	return true, nil
}

func buildCommand(projectRoot string, manifest *Manifest) *cobra.Command {
	root := &cobra.Command{
		Use:     manifest.Command.Use,
		Short:   manifest.Command.Short,
		Aliases: append([]string(nil), manifest.Command.Aliases...),
	}

	r := &runner{projectRoot: projectRoot, manifest: manifest}
	for i := range manifest.Generators {
		root.AddCommand(r.buildGeneratorCommand(&manifest.Generators[i]))
	}
	return root
}

type runner struct {
	projectRoot string
	manifest    *Manifest
}

type invocation struct {
	args  map[string]string
	flags map[string]any
}

func (r *runner) buildGeneratorCommand(spec *GeneratorSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   generatorUse(spec),
		Short:                 spec.Short,
		Aliases:               append([]string(nil), spec.Aliases...),
		DisableFlagsInUseLine: true,
	}

	for _, flag := range spec.Flags {
		switch flag.Type {
		case "string":
			defaultValue, _ := flag.Default.(string)
			cmd.Flags().StringP(flag.Name, flag.Shorthand, defaultValue, flag.Description)
		case "bool":
			defaultValue, _ := flag.Default.(bool)
			cmd.Flags().BoolP(flag.Name, flag.Shorthand, defaultValue, flag.Description)
		}
	}

	minArgs := 0
	for _, arg := range spec.Args {
		if arg.Required {
			minArgs++
		}
	}
	cmd.Args = func(_ *cobra.Command, args []string) error {
		if len(args) < minArgs || len(args) > len(spec.Args) {
			if minArgs == len(spec.Args) {
				return fmt.Errorf("accepts %d arg(s), received %d", len(spec.Args), len(args))
			}
			return fmt.Errorf("accepts between %d and %d arg(s), received %d", minArgs, len(spec.Args), len(args))
		}
		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, values []string) error {
		inv, err := invocationFromCommand(spec, cmd, values)
		if err != nil {
			return err
		}
		return r.execute(spec, inv, cmd.OutOrStdout())
	}
	return cmd
}

func generatorUse(spec *GeneratorSpec) string {
	parts := []string{spec.Name}
	for _, arg := range spec.Args {
		if arg.Required {
			parts = append(parts, "["+arg.Name+"]")
		} else {
			parts = append(parts, "["+arg.Name+"?]")
		}
	}
	return strings.Join(parts, " ")
}

func invocationFromCommand(spec *GeneratorSpec, cmd *cobra.Command, values []string) (invocation, error) {
	inv := invocation{args: make(map[string]string, len(spec.Args)), flags: make(map[string]any, len(spec.Flags))}
	for i, arg := range spec.Args {
		if i < len(values) {
			inv.args[arg.Name] = values[i]
		}
	}
	for _, flag := range spec.Flags {
		switch flag.Type {
		case "string":
			value, err := cmd.Flags().GetString(flag.Name)
			if err != nil {
				return invocation{}, err
			}
			inv.flags[flag.Name] = value
		case "bool":
			value, err := cmd.Flags().GetBool(flag.Name)
			if err != nil {
				return invocation{}, err
			}
			inv.flags[flag.Name] = value
		default:
			return invocation{}, fmt.Errorf("unsupported flag type %q", flag.Type)
		}
	}
	return inv, nil
}

func (r *runner) execute(spec *GeneratorSpec, inv invocation, out io.Writer) error {
	for _, output := range spec.Outputs {
		if err := r.renderOutput(spec, output, inv, out); err != nil {
			return err
		}
	}
	for _, childName := range spec.Compose {
		child := r.findGenerator(childName)
		if child == nil {
			return fmt.Errorf("generator %q composes unknown generator %q", spec.Name, childName)
		}
		childInv := invocation{args: make(map[string]string, len(child.Args)), flags: make(map[string]any, len(child.Flags))}
		for _, arg := range child.Args {
			childInv.args[arg.Name] = inv.args[arg.Name]
		}
		for _, flag := range child.Flags {
			if inherited, ok := inv.flags[flag.Name]; ok {
				childInv.flags[flag.Name] = inherited
				continue
			}
			switch flag.Type {
			case "string":
				value, _ := flag.Default.(string)
				childInv.flags[flag.Name] = value
			case "bool":
				value, _ := flag.Default.(bool)
				childInv.flags[flag.Name] = value
			}
		}
		if err := r.execute(child, childInv, out); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) findGenerator(name string) *GeneratorSpec {
	for i := range r.manifest.Generators {
		if r.manifest.Generators[i].Name == name {
			return &r.manifest.Generators[i]
		}
	}
	return nil
}

func stringify(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprint(v)
	}
}
