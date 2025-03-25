package fileserver

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const sniffLen = 512

var unixEpochTime = time.Unix(0, 0)

// For parsing this time format, see [ParseTime].
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// GODEBUG=httpservecontentkeepheaders=1 restores the pre-1.23 behavior of not deleting
// Cache-Control, Content-Encoding, Etag, or Last-Modified headers on ServeContent errors.
var httpservecontentkeepheaders_del = errors.New("httpservecontentkeepheaders")

type yxifuFileHandler struct {
	root     http.FileSystem
	download bool
	upload   bool
	delete   bool
}

func FileServer(root http.FileSystem, download, upload, delete bool) http.Handler {
	return &yxifuFileHandler{root, download, upload, delete}
}

func (f *yxifuFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	//w.Write([]byte(upath))
	//w.Write([]byte("\n"))
	//w.Write([]byte(r.Method))
	//w.Write([]byte("\n"))
	if r.Method == http.MethodGet {
		// 打开，下载
		f.Get(w, r, upath)

	} else if r.Method == http.MethodPut {
		// 上传
		f.Put(w, r, upath)
	} else if r.Method == http.MethodPost {
		// 重命名

	} else if r.Method == http.MethodDelete {
		// 删除
	}

	//serveFile(w, r, f.root, path.Clean(upath), true)
}

func (f *yxifuFileHandler) Get(w http.ResponseWriter, r *http.Request, upath string) {

	file, err := f.root.Open(upath)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return
	}

	if stat.IsDir() {
		//w.Write([]byte("dir\n"))
		dirList(w, r, file, f.download, f.upload) //显示文件列表
	} else {
		if !f.download {
			Error(w, "禁止下载", 200)
			return
		}
		//w.Write([]byte("file"))
		openFile(w, r, file, stat)
	}

}

func (f *yxifuFileHandler) Put(w http.ResponseWriter, r *http.Request, upath string) {
}

type anyDirs interface {
	len() int
	name(i int) string
	isDir(i int) bool
}

func openFile(w http.ResponseWriter, r *http.Request, f http.File, fs fs.FileInfo) {

	w.Header().Set("Last-Modified", fs.ModTime().UTC().Format(TimeFormat))
	w.Header().Set("Date", time.Now().UTC().Format(TimeFormat))

	ctype := mime.TypeByExtension(filepath.Ext(fs.Name()))
	if ctype == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [sniffLen]byte
		n, _ := io.ReadFull(f, buf[:])
		ctype = http.DetectContentType(buf[:n])
		_, err := f.Seek(0, io.SeekStart) // rewind to output whole file
		if err != nil {
			Error(w, "seeker can't seek", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", ctype)
	// 设置响应头，告知客户端这是一个附件（即下载），并提供默认文件名
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fs.Name()))
	//w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", string(fs.Size()))

	// 将文件内容复制到HTTP响应流
	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, "Error sending file.", 500)
		return
	}
}
func dirList(w http.ResponseWriter, r *http.Request, f http.File, download, upload bool) {
	// Prefer to use ReClient
	// because the former doesn't require calling
	// Stat on every entry of a directory on Unix.
	var err error

	list, err := f.Readdir(-1)

	if err != nil {
		//logf(r, "http: error reading directory: %v", err)
		//Error(w, "Error reading directory", StatusInternalServerError)
		return
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!doctype html>\n")
	fmt.Fprintf(w, "<meta name=\"viewport\" content=\"width=device-width\">\n")

	if upload {
		//显示上传代码
		fmt.Fprintf(w, "<form action=\"\" method=\"Put\" enctype=\"multipart/form-data\">")
		fmt.Fprintf(w, "<label for=\"fileUpload\">选择文件:</label>")
		fmt.Fprintf(w, "<input type=\"file\" id=\"fileUpload\" name=\"fileUpload\">")
		fmt.Fprintf(w, "<button type=\"submit\">上传</button>")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "<hr/>")
	}
	fmt.Fprintf(w, "<pre>\n")
	for i, n := 0, len(list); i < n; i++ {
		name := list[i].Name()
		if list[i].IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		if download || list[i].IsDir() {
			url := url.URL{Path: name}
			//fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), name)
		} else {
			fmt.Fprintf(w, "%s\n", name)
		}
	}
	fmt.Fprintf(w, "</pre>\n")
}

// Helper handlers

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
// The error message should be plain text.
//
// Error deletes the Content-Length header,
// sets Content-Type to “text/plain; charset=utf-8”,
// and sets X-Content-Type-Options to “nosniff”.
// This configures the header properly for the error message,
// in case the caller had set it up expecting a successful output.
func Error(w http.ResponseWriter, error string, code int) {
	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del("Content-Length")

	// There might be content type already set, but we reset it to
	// text/plain for the error message.
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
