package pdf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// CutA3ToA4 使用pdfcpu cut命令将A3的PDF切割成A4大小，每页切成两页
func CutA3ToA4(inputFilePath, outDir string, outFileName string) error {
	selectedPages := []string{}
	//cut := new(model.CUT)
	//conf := model.NewDefaultConfiguration()
	//conf.
	// 创建一个新的 *model.Cut 实例
	cut := &model.Cut{Vert: []float64{0.5}}
	// 调用 api.CutFile 函数，传入正确的参数
	err := api.CutFile(inputFilePath, outDir, "", selectedPages, cut, nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	pageCount, err := api.PageCountFile(inputFilePath)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	fmt.Printf("pageCount: %v\n", pageCount)

	outFile := strings.TrimSuffix(filepath.Base(inputFilePath), ".pdf")
	AddFiless := []string{}
	DeleteFiles := []string{}

	for i := 1; i <= pageCount; i++ {
		outFilePage := filepath.Join(outDir, fmt.Sprintf("%s_page_%d.pdf", filepath.Base(outFile), i))
		fmt.Printf("outFilePage: %v\n", outFilePage)
		api.SplitFile(outFilePage, outDir, 1, nil)

		pageCount2, _ := api.PageCountFile(outFilePage)
		for j := 1; j <= pageCount2; j++ {
			outFilePage2 := filepath.Join(outDir, fmt.Sprintf("%s_page_%d_%d.pdf", filepath.Base(outFile), i, j))
			if j == 1 {
				DeleteFiles = append(DeleteFiles, outFilePage2)
				continue
			}
			AddFiless = append(AddFiless, outFilePage2)
		}
		DeleteFiles = append(DeleteFiles, outFilePage)
	}

	api.MergeAppendFile(AddFiless, filepath.Join(outDir, outFileName), false, nil)
	DeleteFiles = append(DeleteFiles, AddFiless...)
	fmt.Printf("AddFiless: %v\n", AddFiless)

	for _, f := range DeleteFiles {
		os.Remove(f)
	}

	fmt.Printf("输出: %v\n", filepath.Join(outDir, outFileName))
	return nil
}
