package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var repositoryCmd = &cobra.Command{
	Use:                   "repository [name]",
	Short:                 "Generate repository",
	Aliases:               []string{"r"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runRepository,
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
	repositoryStructName := cases.Title(language.English).String(repositoryName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	repositorysDir := filepath.Join(projectRoot, "internal", "repository")
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
