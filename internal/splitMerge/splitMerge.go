package splitMerge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Split(filePath, splitSize, splitOutputDir string) bool {
	if !fileExists(filePath) {
		fmt.Printf("filePath:%s not exists", filePath)
		return false
	}
	if !isRegularFile(filePath) {
		fmt.Printf("filePath:%s is not a regular file", filePath)
		return false
	}
	size, ok := validateSplitSize(splitSize)
	if !ok {
		return false
	}
	if splitOutputDir == "" {
		splitOutputDir = filepath.Dir(filePath)
	} else if !fileExists(splitOutputDir) {
		fmt.Printf("splitOutputDir:%s not exists", splitOutputDir)
		return false
	}

	fmt.Println("Being Split")
	SplitFile(filePath, size, splitOutputDir)
	fmt.Println("END Split")
	return true
}
func Merge(firstFilePath, mergeOutputDir string) bool {
	if !fileExists(firstFilePath) {
		fmt.Printf("filePath:%s not exists", firstFilePath)
		return false
	}
	if !isRegularFile(firstFilePath) {
		fmt.Printf("filePath:%s is not a regular file", firstFilePath)
		return false
	}
	if len(firstFilePath) < 5 || firstFilePath[len(firstFilePath)-5:] != ".s001" {
		fmt.Printf("The file must be the first file(*.s001), got: %s\n", firstFilePath)
		return false
	}
	outputFile := ""
	if mergeOutputDir == "" {
		outputFile = removeSuffix(firstFilePath, ".s001")
	} else if !CheckFileExist(mergeOutputDir) {
		fmt.Println("目录不存在：" + outputFile)
		return false
	} else {
		fn := filepath.Base(firstFilePath)
		fn = removeSuffix(fn, ".s001")
		outputFile = filepath.Join(mergeOutputDir, fn)
	}
	mergeFile(firstFilePath, outputFile)
	return true

}

// validateSplitSize 验证 splitSize 格式为 %d[KkMmGg] 并返回字节大小
func validateSplitSize(splitSize string) (int64, bool) {
	// 检查输入是否为空
	if len(splitSize) == 0 {
		fmt.Println("splitSize cannot be empty")
		return 0, false
	}
	// 提取数字部分
	numStr := ""
	for i := 0; i < len(splitSize); i++ {
		if splitSize[i] >= '0' && splitSize[i] <= '9' {
			numStr += string(splitSize[i])
		} else {
			break
		}
	}
	if len(numStr) == 0 {
		fmt.Println("Invalid splitSize format: missing number part")
		return 0, false
	}
	// 将数字部分转换为 int64
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid splitSize format: %v\n", err)
		return 0, false
	}
	// 提取单位部分
	unitStr := splitSize[len(numStr):]
	var multiplier int64
	switch unitStr {
	case "K", "k":
		multiplier = 1024
	case "M", "m":
		multiplier = 1024 * 1024
	case "G", "g":
		multiplier = 1024 * 1024 * 1024
	default:
		if unitStr != "" {
			fmt.Println("Invalid splitSize format: invalid unit")
			return 0, false
		}
		multiplier = 1
	}
	return num * multiplier, true
}

// 添加判断文件是否存在的函数
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsRegularFile 判断给定路径的文件是否为普通文件
func isRegularFile(filePath string) bool {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// 若获取信息出错，直接返回 false
		return false
	}
	// 通过 Mode().IsRegular() 判断是否为普通文件
	return fileInfo.Mode().IsRegular()
}

// SplitFile Split the input file into multiple files according to the specified size.
func SplitFile(inputFilePath string, chunkSize int64, outputDir string) error {
	var outputPrefix string
	outputFileName := filepath.Base(inputFilePath)
	fmt.Printf("outputFileName: %v\n", outputFileName)
	outputPrefix = filepath.Join(outputDir, outputFileName+".s")
	fmt.Printf("outputPrefix: %v\n", outputPrefix)
	//var chunkSize int64 = 1024
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	if fileSize <= chunkSize {
		// If the file size is less than or equal to the chunk size, no splitting is required.。
		return nil
	}

	var count int
	var ts int64 = 0
	for {
		// create output file
		outputFileName := fmt.Sprintf("%s%03d", outputPrefix, count+1)
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		reader := io.LimitReader(inputFile, chunkSize)

		_, err = io.Copy(outputFile, reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		count++
		ts += chunkSize
		if ts >= fileSize {
			break
		}
	}

	return nil
}

func mergeFile(firstFilePath, outputFile string) error {
	//outputFile := removeSuffix(firstFilePath, ".s001")
	// fmt.Printf("outputFile: %v\n", outputFile)
	// fmt.Printf("firstFilePath: %v\n", firstFilePath)
	// if !isEndWithS001(firstFilePath) {
	// 	fmt.Println("The file must be the first file(*.s001)")
	// 	return nil
	// }
	fmt.Println("outputFile:", outputFile)
	if CheckFileExist(outputFile) {
		fmt.Println("文件已经存在！" + outputFile)
		return nil
	}

	// 创建输出文件。
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var count int = 0
	for {
		// 创建输出文件。
		filename := fmt.Sprintf("%s%s%03d", outputFile, ".s", count+1)

		if !CheckFileExist(filename) {
			break
		}

		// 打开每个文件。
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将文件内容复制到输出文件。
		_, err = io.Copy(outFile, file)
		if err != nil {
			return err
		}
		count++
	}

	return nil

}
func isEndWithS001(str string) bool {
	return len(str) >= 5 && str[len(str)-5:] == ".s001"
}

// CheckFileExist 检查文件是否存在。
func CheckFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在
			return false
		}
		// 其他错误
		fmt.Println("Error:", err)
		return false
	}
	// 文件存在
	return true
}
func removeSuffix(s, suffix string) string {
	// 使用 strings.HasSuffix 检查字符串是否以指定的后缀结尾
	if strings.HasSuffix(s, suffix) {
		// 使用切片去掉后缀
		return s[:len(s)-len(suffix)]
	}
	// 如果字符串不以该后缀结尾，则返回原始字符串
	return s
}
