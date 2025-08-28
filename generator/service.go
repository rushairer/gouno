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

var serviceCmd = &cobra.Command{
	Use:                   "service [name]",
	Short:                 "Generate service",
	Aliases:               []string{"s"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runService,
}

func runService(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("service name required")
		os.Exit(1)
	}
	serviceName := args[0]
	if serviceName == "" {
		fmt.Println("service name is empty")
		os.Exit(1)
	}
	serviceFileName := fmt.Sprintf("%s.go", serviceName)
	serviceStructName := cases.Title(language.English).String(serviceName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	servicesDir := filepath.Join(projectRoot, "internal", "service")
	if _, innerErr := os.Stat(servicesDir); os.IsNotExist(innerErr) {
		err = os.MkdirAll(servicesDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create service directory: %v", err)
		}
		log.Printf("Created directory: %s", servicesDir)
	}

	serviceFilePath := filepath.Join(servicesDir, serviceFileName)
	serviceContent := fmt.Sprintf(
		serviceTemplate,
		serviceStructName,
		serviceStructName,
		serviceStructName,
		serviceStructName,
		serviceStructName,
	)

	err = os.WriteFile(serviceFilePath, []byte(serviceContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create service file: %v", err)
	}
	log.Printf("Created service file: %s", serviceFilePath)
}

const serviceTemplate = `package service

import "context"

type %sService struct {
}

func New%sService() *%sService {
	return &%sService{}
}

func (c *%sService) Get(ctx context.Context) (result string, err error) {
	return
}`
