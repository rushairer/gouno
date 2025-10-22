package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rushairer/gouno/utilitiy"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:                   "domain [name]",
	Short:                 "Generate domain",
	Aliases:               []string{"d"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runDomain,
}

var defaultDomainPath = filepath.Join("internal", "domain")

func init() {
	domainCmd.Flags().StringP("path", "p", defaultDomainPath, "path to domain")
	domainCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

func runDomain(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("domain name required")
		os.Exit(1)
	}
	domainName := args[0]
	if domainName == "" {
		fmt.Println("domain name is empty")
		os.Exit(1)
	}
	domainFileName := fmt.Sprintf("%s.go", domainName)
	domainStructName := utilitiy.ToCamelCase(domainName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	path := defaultDomainPath
	if flag := cmd.Flag("path"); flag != nil {
		path = flag.Value.String()
	}
	domainsDir := filepath.Join(projectRoot, path)
	if _, innerErr := os.Stat(domainsDir); os.IsNotExist(innerErr) {
		err = os.MkdirAll(domainsDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create domain directory: %v", err)
		}
		log.Printf("Created directory: %s", domainsDir)
	}

	domainFilePath := filepath.Join(domainsDir, domainFileName)
	domainContent := fmt.Sprintf(
		domainTemplate,
		domainStructName,
		domainStructName,
		domainStructName,
		domainStructName,
		domainStructName,
	)

	if force, _ := cmd.Flags().GetBool("force"); !force {
		// Check if the file already exists, 如果存在，询问是否覆盖
		if _, innerErr := os.Stat(domainFilePath); innerErr == nil {
			var confirm string
			fmt.Printf("Domain file already exists: %s, do you want to overwrite it? (y/n) ", domainFilePath)
			fmt.Scanln(&confirm)
			if confirm != "y" {
				log.Fatalf("Domain file not overwritten: %s", domainFilePath)
			}
		}
	}

	err = os.WriteFile(domainFilePath, []byte(domainContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create domain file: %v", err)
	}
	log.Printf("Created domain file: %s", domainFilePath)
}

const domainTemplate = `package domain

import "context"

type %s struct {
}

func New%s() *%s {
	return &%s{}
}

func (d *%s) Foo(ctx context.Context) (bar string, err error) {
	return
}`
