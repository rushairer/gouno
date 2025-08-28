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

var domainCmd = &cobra.Command{
	Use:                   "domain [name]",
	Short:                 "Generate domain",
	Aliases:               []string{"d"},
	DisableFlagsInUseLine: true,
	Run:                   runDomain,
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
	domainStructName := cases.Title(language.English).String(domainName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	domainsDir := filepath.Join(projectRoot, "internal", "domain")
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

func (c *%s) Get(ctx context.Context) (result string, err error) {
	return
}`
