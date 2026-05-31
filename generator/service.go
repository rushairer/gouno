package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:                   "service [name]",
	Short:                 "Generate service",
	Aliases:               []string{"s"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(cmd, args, "service", defaultServicePath)
	},
}

var defaultServicePath = filepath.Join("internal", "service")

func init() {
	serviceCmd.Flags().StringP("path", "p", defaultServicePath, "path to service")
	serviceCmd.Flags().BoolP("force", "f", false, "force overwrite")
	serviceCmd.Flags().String("template-set", "", "template set name")
}
