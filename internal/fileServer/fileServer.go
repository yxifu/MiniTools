package fileserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func Start(port int, dir, pathPrefix string) {
	fmt.Printf("port: %v\n", port)
	fmt.Printf("dir: %v\n", dir)
	fmt.Printf("pathPrefix: %v\n", pathPrefix)
	if port < 1 && port > 65535 {
		log.Fatalf("端口号不合法：%d", port)
		return
	}
	if !isDirExists(dir) {
		log.Fatalf("目录不存在：%s", dir)
		return
	}
	if len(pathPrefix) > 0 && pathPrefix[len(pathPrefix)-1] != '/' {
		pathPrefix += "/"
	}

	fs := http.FileServer(http.Dir(dir))
	http.Handle(pathPrefix, http.StripPrefix(pathPrefix, fs))
	//http.Handle(pathPrefix, fs)

	// 启动HTTP服务器
	log.Printf("Starting server on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// IsDirExists 判断指定路径的文件夹是否存在
func isDirExists(dirPath string) bool {
	// 获取文件信息
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		// 若获取信息出错且错误为文件不存在，则文件夹不存在
		if os.IsNotExist(err) {
			return false
		}
	}
	// 判断是否为目录
	return fileInfo.IsDir()
}
