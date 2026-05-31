package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rushairer/gouno/utility"
	"github.com/spf13/cobra"
)

// generateFile 是所有代码生成器的公共逻辑：
// 1. 将名称转为驼峰命名
// 2. 在目标目录下创建文件
// 3. 若文件已存在且未指定 --force，提示用户确认覆盖
func generateFile(cmd *cobra.Command, args []string, typeName, defaultPath, tmpl string) error {
	name := args[0]
	structName := utility.ToCamelCase(name)

	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	path := defaultPath
	if flag := cmd.Flag("path"); flag != nil {
		path = flag.Value.String()
	}
	dir := filepath.Join(projectRoot, path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", typeName, err)
		}
		cmd.Printf("Created directory: %s\n", dir)
	}

	filePath := filepath.Join(dir, fmt.Sprintf("%s.go", name))
	if force, _ := cmd.Flags().GetBool("force"); !force {
		if _, err := os.Stat(filePath); err == nil {
			var confirm string
			displayName := string(typeName[0]-32) + typeName[1:]
			fmt.Printf("%s file already exists: %s, do you want to overwrite it? (y/n) ", displayName, filePath)
			fmt.Scanln(&confirm)
			if confirm != "y" {
				return fmt.Errorf("%s file not overwritten: %s", typeName, filePath)
			}
		}
	}

	content := fmt.Sprintf(tmpl, structName, structName, structName, structName, structName)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create %s file: %w", typeName, err)
	}
	cmd.Printf("Created %s file: %s\n", typeName, filePath)
	return nil
}
