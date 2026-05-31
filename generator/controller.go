package generator

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:                   "controller [name]",
	Short:                 "Generate controller",
	Aliases:               []string{"c"},
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		generateFile(cmd, args, "controller", defaultControllerPath, controllerTemplate)
	},
}

var defaultControllerPath = filepath.Join("controller")

func init() {
	controllerCmd.Flags().StringP("path", "p", defaultControllerPath, "path to controller")
	controllerCmd.Flags().BoolP("force", "f", false, "force overwrite")
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

func (c *%sController) Foo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gouno.NewSuccessResponse("bar"))
}`
