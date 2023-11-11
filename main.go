package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var StaticPath string

func main() {
	// 检查是否存在 static 文件夹
	if StaticPath == "" {
		if _, err := os.Stat("static"); os.IsNotExist(err) {
			fmt.Println("static 文件夹不存在，程序退出。")
			return
		}
	}

	// 设置上传文件的处理函数
	http.HandleFunc("/upload", uploadHandler)

	// 启动Web服务器
	fmt.Println("服务器已启动，监听端口 8080...")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 检查是否存在 static/images 文件夹，不存在则创建
	imagesDir := filepath.Join("static", "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.MkdirAll(imagesDir, os.ModePerm)
		if err != nil {
			http.Error(w, "无法创建目录", http.StatusInternalServerError)
			return
		}
	}

	// 解析表单，获取上传文件
	err := r.ParseMultipartForm(10 << 20) // 设置最大文件大小为 10MB
	if err != nil {
		http.Error(w, "无法解析表单", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "无法获取文件", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 生成新的文件名，以时间戳和年份分类
	timestamp := time.Now().Unix()
	year := time.Now().Year()
	newFileName := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	yearDir := filepath.Join(imagesDir, strconv.Itoa(year))

	// 检查年份文件夹是否存在，不存在则创建
	if _, err := os.Stat(yearDir); os.IsNotExist(err) {
		err := os.Mkdir(yearDir, os.ModePerm)
		if err != nil {
			http.Error(w, "无法创建年份目录", http.StatusInternalServerError)
			return
		}
	}

	// 创建并写入新文件
	newFilePath := filepath.Join(yearDir, newFileName)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		http.Error(w, "无法创建新文件", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "无法写入文件", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "文件上传成功，保存路径：%s", newFilePath)
}

func init() {
	flag.StringVar(&StaticPath, "path", "", "hugo static path")
	flag.Parse()
}
