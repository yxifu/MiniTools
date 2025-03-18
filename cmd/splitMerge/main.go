package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	splitMerge "MiniTools/internal/splitMerge"
)

// var outputPath string
// var crateOutputDir bool
var rootCmd = &cobra.Command{
	Use:   "splitMerge",
	Short: "文件分割与合并",
	Long:  `文件分割与合并`,
	//Args:  cobra.MinimumNArgs(1),
	/*Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("splitMerge: " + strings.Join(args, " "))
		//ExecuteCut(args[0], outputPath)
	},*/
}

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "将文件按大小切割成多个文件",
	Long:  `将文件按大小切割成多个文件`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("split: " + strings.Join(args, " "))
		//ExecuteCut(args[0], outputPath)
		inputFile := args[0]
		splitMerge.Split(inputFile, splitSize, splitOutputDir)
	},
	Example: `splitMerge split ./abc.txt -s 10k
splitMerge split ./abc.txt -s 10k -o ./test/`,
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "将分割后的小文件，合并成一个文件",
	Long:  "将分割后的小文件，合并成一个文件",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("merge: " + strings.Join(args, " "))
		//ExecuteCut(args[0], outputPath)
		//inputFilePath := args[0]
		//splitMerge.Split(inputFilePath, splitSize, splitOutputDir)

		inputFile := args[0]
		splitMerge.Merge(inputFile, mergeOutputDir)
	},
	Example: `splitMerge merge ./abc.txt.s001
splitMerge merge ./abc.txt.s001 -o ./test/`,
}

var splitSize string
var splitOutputDir string
var mergeOutputDir string

func init() {
	splitCmd.Flags().StringVarP(&splitSize, "size", "s", "", "每块文件大小（如：10k;20m;30g）")
	splitCmd.Flags().StringVarP(&splitOutputDir, "outputDir", "o", "", "分割后文件输出目录（默认分割文件所在目录）")
	splitCmd.MarkFlagRequired("size")
	rootCmd.AddCommand(splitCmd)

	mergeCmd.Flags().StringVarP(&mergeOutputDir, "outputDir", "o", "", "合并后文件输出目录")
	rootCmd.AddCommand(mergeCmd)
}

func main() {
	rootCmd.Execute()

}
