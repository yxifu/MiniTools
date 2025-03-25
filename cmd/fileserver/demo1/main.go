package main

import (
	fileserver "MiniTools/internal/fileServer"
	"fmt"
	"log"
	"net/http"
)

func main() {
	dir := "./"
	pathPrefix := "/ab/"
	port := 8801
	fs := fileserver.FileServer(http.Dir(dir), true, true, true)
	http.Handle(pathPrefix, http.StripPrefix(pathPrefix, fs))
	//http.Handle(pathPrefix, fs)

	// 启动HTTP服务器
	log.Printf("Starting server on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
