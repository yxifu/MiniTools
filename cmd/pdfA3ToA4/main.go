package main

import (
	"MiniTools/pdf"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var outputPath string
var crateOutputDir bool
var rootCmd = &cobra.Command{
	Use:   "pdfA3ToA4 [A3pdf]",
	Short: "将A3PDF转成A4的",
	Long:  `将A3(横向)PDF切成A4（竖向）`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("Print: " + strings.Join(args, " "))
		ExecuteCut(args[0], outputPath)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputPath, "outputPath", "o", "", "输出pdf输出路径（默认当前目录）")
	rootCmd.PersistentFlags().BoolVarP(&crateOutputDir, "crateOutputDir", "c", false, "如果输出目录不存在则创建")

}

func main() {
	rootCmd.Execute()

}

/*
func test() {
	pdf.CutA3ToA4("./samples/pdfA3ToA4/input_A3.pdf", "./temp", "output_A3.pdf")
}*/

func ExecuteCut(inputPath, outputPath string) {

	outputDir := ""
	outputPdfName := ""
	fmt.Println("输入参数：")
	fmt.Printf("inputPath: %v\n", inputPath)
	fmt.Printf("outputPath: %v\n", outputPath)

	//判断inputPath是否是PDF文件
	if !isPdf(inputPath) {
		log.Fatalln("输入文件不是PDF文件")
	}

	if outputPath == "" {
		//输入目录为输入目录，输出文件名为输入文件名_A4.pdf
		outputDir = "./"
		_, outputPdfName = getDirAndFileName(inputPath)
		outputPdfName = convertToA4FileName(outputPdfName)
	} else if isPdf(outputPath) {
		outputDir, outputPdfName = getDirAndFileName(outputPath)
		//} else if isDir(outputPath) {
		//	outputDir = outputPath
		//	_, outputPdfName = getDirAndFileName(inputPath)
		//	outputPdfName = convertToA4FileName(outputPdfName)
	} else {
		//outputDir, outputPdfName = getDirAndFileName(outputPath)
		outputDir = outputPath
		_, outputPdfName = getDirAndFileName(inputPath)
		outputPdfName = convertToA4FileName(outputPdfName)
	}

	if isDirExists(outputDir) {
		if crateOutputDir {
			createDirIfNotExists(outputDir)
		} else {
			log.Fatalf("输出目录%s不存在\n", outputDir)
		}
	}
	fmt.Printf("outputDir: %v\n", outputDir)
	fmt.Printf("outputPdfName: %v\n", outputPdfName)
	pdf.CutA3ToA4(inputPath, outputDir, outputPdfName)

}

// 假设我们已经实现了 GetDirAndFileName 函数，以下是一个简单的示例实现
// 注意：这里的实现仅为示例，实际使用中需要根据具体需求进行调整
func getDirAndFileName(path string) (string, string) {
	// 这里简单返回路径和文件名，实际需要根据操作系统路径分隔符等进行处理
	// 以下是一个简单的示例，假设路径使用 '/' 分隔
	lastIndex := strings.LastIndex(path, "/")
	if lastIndex == -1 {
		return "", path
	}
	return path[:lastIndex], path[lastIndex+1:]
}

// 定义 isPdf 函数来检查文件是否为 PDF 文件
func isPdf(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".pdf")
}

// convertToA4FileName 将文件名转换为 A4 格式的文件名
func convertToA4FileName(name string) string {
	// 找到文件扩展名的位置
	dotIndex := strings.LastIndex(name, ".")
	if dotIndex == -1 {
		// 如果没有扩展名，直接添加 _A4
		return name + "_A4"
	}
	// 提取文件名和扩展名
	baseName := name[:dotIndex]
	ext := name[dotIndex:]
	// 组合新的文件名
	return baseName + "_A4" + ext
}

// 判断路径是否为目录
func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return fileInfo.IsDir()
}

// 判断目录是否存在
func isDirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// 创建目录
func createDirIfNotExists(path string) error {
	if !isDirExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
