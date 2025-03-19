package main

import (
	fileserver "MiniTools/internal/fileServer"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd.Execute()
}

// var outputPath string
// var crateOutputDir bool
var rootCmd = &cobra.Command{
	Use:   "fileserver",
	Short: "启动文件服务器，只有查看和下载功能",
	Long:  `启动文件服务器，只有查看和下载功能`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("Print: " + strings.Join(args, " "))
		fileserver.Start(port, dir, pathPrefix)
	},
}

var port int
var dir string
var pathPrefix string

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "服务端口")
	rootCmd.Flags().StringVarP(&dir, "dir", "d", ".", "目录")
	rootCmd.Flags().StringVarP(&pathPrefix, "path", "", "/", "路径前缀")
}
