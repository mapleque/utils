/*
改名工具，支持将当前文件夹的所有文件名改成1-n的序号

使用方法：
	拷贝到要改名的文件所在的文件夹运行之

交叉编译：
GOOS=darwin GOARCH=amd64 go build -o rename-mac-amd64 rename.go
GOOS=windows GOARCH=amd64 go build -o rename-win64.exe rename.go
GOOS=windows GOARCH=386 go build -o rename-win32.exe rename.go

*/
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var confirm string
	// 读取当前路径
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("正在读取%s下的所有文件\n", currentDir)

	// 读取路径下所有文件
	fileList := []string{}
	dirList := []string{}
	filepath.Walk(currentDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			dirList = append(dirList, path)
		} else {
			if !strings.HasPrefix(f.Name(), "rename") && !strings.HasPrefix(f.Name(), ".") {
				fileList = append(fileList, f.Name())
			}
		}
		return nil
	})
	if len(fileList) > 0 {
		fmt.Printf("将要修改以下文件的文件名：\n%s\n", strings.Join(fileList, "\n"))
		fmt.Printf("共找到%d个文件\n", len(fileList))
	} else {
		fmt.Println("操作失败：没有发现任何要修改的文件，请确保文件夹下有文件，且文件名不是.或rename开头！")
		fmt.Println("按任意键退出")
		fmt.Scanln(&confirm)
		os.Exit(2)
	}
	// 等用户确认
	fmt.Printf("是否全部重命名？ [y/n]:")
	fmt.Scanln(&confirm)
	if confirm != "y" {
		fmt.Printf("任务已取消！")
		os.Exit(0)
	}
	// 查看文件夹内有没有命名为result的文件夹
	for _, name := range dirList {
		if name == "result" {
			fmt.Println("操作失败：当前文件夹内含有result文件夹，请改名或删除后重新运行！")
			fmt.Println("按任意键退出")
			fmt.Scanln(&confirm)
			os.Exit(2)
		}
	}

	// 创建result文件夹，并将重命名文件copy进来
	fmt.Printf("创建文件夹%s\n", currentDir+string([]byte{os.PathSeparator})+"result")
	os.Mkdir(currentDir+string([]byte{os.PathSeparator})+"result", os.ModePerm)
	num := 0
	for _, name := range fileList {
		num++
		ext := filepath.Ext(name)
		srcName := currentDir + string([]byte{os.PathSeparator}) + name
		tarName := currentDir + string([]byte{os.PathSeparator}) + "result" + string([]byte{os.PathSeparator}) + strconv.Itoa(num) + ext
		fmt.Printf("copy %s -> %s\n", srcName, tarName)
		if _, err := copyFile(srcName, tarName); err != nil {
			log.Fatal("复制文件失败", err)
		}
	}

}

func copyFile(srcName, dstName string) (int64, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
