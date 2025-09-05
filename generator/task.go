package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rushairer/gouno/utilitiy"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:                   "task [name]",
	Short:                 "Generate task",
	Aliases:               []string{"r"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runTask,
}

var defaultTaskPath = filepath.Join("internal", "task")

func init() {
	taskCmd.Flags().StringP("path", "p", defaultTaskPath, "path to task")
	taskCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

func runTask(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("task name required")
		os.Exit(1)
	}
	taskName := args[0]
	if taskName == "" {
		fmt.Println("task name is empty")
		os.Exit(1)
	}
	taskFileName := fmt.Sprintf("%s.go", taskName)
	taskStructName := utilitiy.ToCamelCase(taskName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	path := defaultTaskPath
	if flag := cmd.Flag("path"); flag != nil {
		path = flag.Value.String()
	}
	tasksDir := filepath.Join(projectRoot, path)
	if _, innerErr := os.Stat(tasksDir); os.IsNotExist(innerErr) {
		err = os.MkdirAll(tasksDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create task directory: %v", err)
		}
		log.Printf("Created directory: %s", tasksDir)
	}

	taskFilePath := filepath.Join(tasksDir, taskFileName)
	taskContent := fmt.Sprintf(
		taskTemplate,
		taskStructName,
		taskStructName,
		taskStructName,
		taskStructName,
		taskStructName,
	)

	if force, _ := cmd.Flags().GetBool("force"); !force {
		// Check if the file already exists, 如果存在，询问是否覆盖
		if _, innerErr := os.Stat(taskFilePath); innerErr == nil {
			var confirm string
			fmt.Printf("Task file already exists: %s, do you want to overwrite it? (y/n) ", taskFilePath)
			fmt.Scanln(&confirm)
			if confirm != "y" {
				log.Fatalf("Task file not overwritten: %s", taskFilePath)
			}
		}
	}

	err = os.WriteFile(taskFilePath, []byte(taskContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create task file: %v", err)
	}
	log.Printf("Created task file: %s", taskFilePath)
}

const taskTemplate = `package task

import "context"

type %sTask struct {
}

func New%sTask() *%sTask {
	return &%sTask{}
}

func (c *%sTask) Run(ctx context.Context) error {
	return nil
}`
