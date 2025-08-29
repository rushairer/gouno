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

var controllerCmd = &cobra.Command{
	Use:                   "controller [name]",
	Short:                 "Generate controller",
	Aliases:               []string{"c"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run:                   runController,
}

var defaultControllerPath = filepath.Join("controller")

func init() {
	controllerCmd.Flags().StringP("path", "p", defaultControllerPath, "path to controller")
	controllerCmd.Flags().BoolP("force", "f", false, "force overwrite")
}

func runController(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("controller name required")
		os.Exit(1)
	}
	controllerName := args[0]
	if controllerName == "" {
		fmt.Println("controller name is empty")
		os.Exit(1)
	}
	controllerFileName := fmt.Sprintf("%s.go", controllerName)
	controllerStructName := cases.Title(language.English).String(controllerName)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	path := defaultControllerPath
	if flag := cmd.Flag("path"); flag != nil {
		path = flag.Value.String()
	}
	controllersDir := filepath.Join(projectRoot, path)
	if _, innerErr := os.Stat(controllersDir); os.IsNotExist(innerErr) {
		err = os.MkdirAll(controllersDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create controller directory: %v", err)
		}
		log.Printf("Created directory: %s", controllersDir)
	}

	controllerFilePath := filepath.Join(controllersDir, controllerFileName)
	controllerContent := fmt.Sprintf(
		controllerTemplate,
		controllerStructName,
		controllerStructName,
		controllerStructName,
		controllerStructName,
		controllerStructName,
	)
	if force, _ := cmd.Flags().GetBool("force"); !force {
		// Check if the file already exists, 如果存在，询问是否覆盖
		if _, innerErr := os.Stat(controllerFilePath); innerErr == nil {
			var confirm string
			fmt.Printf("Controller file already exists: %s, do you want to overwrite it? (y/n) ", controllerFilePath)
			fmt.Scanln(&confirm)
			if confirm != "y" {
				log.Fatalf("Controller file not overwritten: %s", controllerFilePath)
			}
		}
	}

	err = os.WriteFile(controllerFilePath, []byte(controllerContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create controller file: %v", err)
	}
	log.Printf("Created controller file: %s", controllerFilePath)
}

const controllerTemplate = `package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/rushairer/gouno"
)

type %sController struct {
}

func New%sController() *%sController {
	return &%sController{}
}

func (c *%sController) Get(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gouno.NewSuccessResponse("foo"))
}`
