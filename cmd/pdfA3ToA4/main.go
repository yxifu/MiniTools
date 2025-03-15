package main

import (
	"MiniTools/pdf"
	"fmt"
)

func main() {
	fmt.Println("将A3PDF转成A4的")
	pdf.CutA3ToA4("./samples/pdfA3ToA4/input_A3.pdf", "./temp", "output_A3.pdf")

}
