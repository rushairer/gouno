package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:                   "controller [name]",
	Short:                 "Generate controller",
	Aliases:               []string{"c"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(cmd, args, "controller", defaultControllerPath)
	},
}

var defaultControllerPath = filepath.Join("controller")

func init() {
	controllerCmd.Flags().StringP("path", "p", defaultControllerPath, "path to controller")
	controllerCmd.Flags().BoolP("force", "f", false, "force overwrite")
	controllerCmd.Flags().String("template-set", "", "template set name")
}
