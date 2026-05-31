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
	Run: func(cmd *cobra.Command, args []string) {
		generateFile(cmd, args, "repository", defaultRepositoryPath, repositoryTemplate)
	},
}

var defaultRepositoryPath = filepath.Join("internal", "repository")

func init() {
	repositoryCmd.Flags().StringP("path", "p", defaultRepositoryPath, "path to repository")
	repositoryCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

const repositoryTemplate = `package repository

import "context"

type %sRepository struct {
}

func New%sRepository() *%sRepository {
	return &%sRepository{}
}

func (r *%sRepository) Foo(ctx context.Context) (bar string, err error) {
	return
}`
