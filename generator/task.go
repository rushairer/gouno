package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:                   "task [name]",
	Short:                 "Generate task",
	Aliases:               []string{"t"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(cmd, args, "task", defaultTaskPath)
	},
}

var defaultTaskPath = filepath.Join("internal", "task")

func init() {
	taskCmd.Flags().StringP("path", "p", defaultTaskPath, "path to task")
	taskCmd.Flags().BoolP("force", "f", false, "force overwrite")
	taskCmd.Flags().String("template-set", "", "template set name")
}
