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
	Run: func(cmd *cobra.Command, args []string) {
		generateFile(cmd, args, "task", defaultTaskPath, taskTemplate)
	},
}

var defaultTaskPath = filepath.Join("internal", "task")

func init() {
	taskCmd.Flags().StringP("path", "p", defaultTaskPath, "path to task")
	taskCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

const taskTemplate = `package task

import "context"

type %sTask struct {
}

func New%sTask() *%sTask {
	return &%sTask{}
}

func (t *%sTask) Run(ctx context.Context) error {
	return nil
}`
