package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rushairer/gouno/utilitiy"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:                   "service [name]",
	Short:                 "Generate service",
	Aliases:               []string{"s"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runService,
}

var defaultServicePath = filepath.Join("internal", "service")

func init() {
	serviceCmd.Flags().StringP("path", "p", defaultServicePath, "path to service")
	serviceCmd.Flags().BoolP("force", "f", false, "force overwrite")
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
	serviceStructName := utilitiy.ToCamelCase(serviceName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	path := defaultServicePath
	if flag := cmd.Flag("path"); flag != nil {
		path = flag.Value.String()
	}
	servicesDir := filepath.Join(projectRoot, path)
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

	if force, _ := cmd.Flags().GetBool("force"); !force {
		// Check if the file already exists, 如果存在，询问是否覆盖
		if _, innerErr := os.Stat(serviceFilePath); innerErr == nil {
			var confirm string
			fmt.Printf("Service file already exists: %s, do you want to overwrite it? (y/n) ", serviceFilePath)
			fmt.Scanln(&confirm)
			if confirm != "y" {
				log.Fatalf("Service file not overwritten: %s", serviceFilePath)
			}
		}
	}

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
