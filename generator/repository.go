package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rushairer/gouno/utilitiy"
	"github.com/spf13/cobra"
)

var repositoryCmd = &cobra.Command{
	Use:                   "repository [name]",
	Short:                 "Generate repository",
	Aliases:               []string{"r"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runRepository,
}

var defaultRepositoryPath = filepath.Join("internal", "repository")

func init() {
	repositoryCmd.Flags().StringP("path", "p", defaultRepositoryPath, "path to repository")
	repositoryCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

func runRepository(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("repository name required")
		os.Exit(1)
	}
	repositoryName := args[0]
	if repositoryName == "" {
		fmt.Println("repository name is empty")
		os.Exit(1)
	}
	repositoryFileName := fmt.Sprintf("%s.go", repositoryName)
	repositoryStructName := utilitiy.ToCamelCase(repositoryName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	path := defaultRepositoryPath
	if flag := cmd.Flag("path"); flag != nil {
		path = flag.Value.String()
	}
	repositorysDir := filepath.Join(projectRoot, path)
	if _, innerErr := os.Stat(repositorysDir); os.IsNotExist(innerErr) {
		err = os.MkdirAll(repositorysDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create repository directory: %v", err)
		}
		log.Printf("Created directory: %s", repositorysDir)
	}

	repositoryFilePath := filepath.Join(repositorysDir, repositoryFileName)
	repositoryContent := fmt.Sprintf(
		repositoryTemplate,
		repositoryStructName,
		repositoryStructName,
		repositoryStructName,
		repositoryStructName,
		repositoryStructName,
	)

	if force, _ := cmd.Flags().GetBool("force"); !force {
		// Check if the file already exists, 如果存在，询问是否覆盖
		if _, innerErr := os.Stat(repositoryFilePath); innerErr == nil {
			var confirm string
			fmt.Printf("Repository file already exists: %s, do you want to overwrite it? (y/n) ", repositoryFilePath)
			fmt.Scanln(&confirm)
			if confirm != "y" {
				log.Fatalf("Repository file not overwritten: %s", repositoryFilePath)
			}
		}
	}

	err = os.WriteFile(repositoryFilePath, []byte(repositoryContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create repository file: %v", err)
	}
	log.Printf("Created repository file: %s", repositoryFilePath)
}

const repositoryTemplate = `package repository

import "context"

type %sRepository struct {
}

func New%sRepository() *%sRepository {
	return &%sRepository{}
}

func (c *%sRepository) Get(ctx context.Context) (result string, err error) {
	return
}`
