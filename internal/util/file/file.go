package file

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/h2non/filetype"
)

// 判断路径是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 路径存在
	}
	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 其他错误
}

func CreateDirIfNotExist(path string) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}

	if !exists {
		err := os.MkdirAll(path, os.ModePerm) // 创建多级目录
		if err != nil {
			return err
		}
	}

	return nil
}

func RmDirIfExist(path string) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("目录 %s 存在，正在删除...\n", path)
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
		fmt.Printf("目录 %s 删除成功！\n", path)

	}
	return nil
}

func SaveFileToLocal(b []byte, filePath string) error {
	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, bytes.NewReader(b))
	if err != nil {
		return err
	}
	return nil
}

func ZipFile(savePath string, files []string) error {
	// 创建输出 ZIP 文件
	zipFile, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建 ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加文件到 ZIP 压缩包
	for _, file := range files {
		err := addFileToZip(zipWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// 将文件添加到 ZIP 压缩包
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	// 打开要压缩的文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	// 创建文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("error creating file header: %w", err)
	}
	header.Name = filepath.Base(filePath)

	// 创建文件在 ZIP 压缩包中的位置
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("error creating file in zip: %w", err)
	}

	// 将文件内容复制到 ZIP 压缩包
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("error writing file to zip: %w", err)
	}

	return nil
}

func IsVideo(data []byte) bool {
	// 检测文件类型
	kind, err := filetype.Match(data)
	if err != nil {
		log.Fatal(err)
	}

	// 判断文件类型
	switch kind.MIME.Type {
	case "image":
		fmt.Println("The file is an image.")
		return false
	case "video":
		fmt.Println("The file is a video.")
		return true
	default:
		fmt.Printf("The file is of an unknown type: %s/%s\n", kind.MIME.Type, kind.MIME.Subtype)
		return false
	}
}

func GetPicSizeByURL(ctx context.Context, url string) (width int, height int, err error) {
	c := g.Client()
	r, err := c.Get(ctx, url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()
	// 解码图片
	img, _, err := image.DecodeConfig(r.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// 获取图片的长宽
	width = img.Width
	height = img.Height
	return
}

func FileExists(path string) (bool, error) {
	return pathExists(path)
}
