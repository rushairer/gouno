package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var repositoryCmd = &cobra.Command{
	Use:                   "repository [name]",
	Short:                 "Generate repository",
	Aliases:               []string{"r"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFile(cmd, args, "repository", defaultRepositoryPath)
	},
}

var defaultRepositoryPath = filepath.Join("internal", "repository")

func init() {
	repositoryCmd.Flags().StringP("path", "p", defaultRepositoryPath, "path to repository")
	repositoryCmd.Flags().BoolP("force", "f", false, "force overwrite")
	repositoryCmd.Flags().String("template-set", "", "template set name")
}
