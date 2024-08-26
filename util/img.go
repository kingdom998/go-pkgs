package util

import (
	"encoding/base64"
	"log"
	"os"
)

type Image string

func (i *Image) B64encode() (string, error) {
	// 读取PNG图片文件
	imageData, err := os.ReadFile(string(*i))
	if err != nil {
		log.Fatalf("Error reading image file: %v", err)
		return "", err
	}

	// 将图片内容转换为Base64编码
	base64Encoded := base64.StdEncoding.EncodeToString(imageData)
	return base64Encoded, nil
}

func (i *Image) B64encode2File(outputFile string) (err error) {
	base64Encoded, err := i.B64encode()
	if err != nil {
		return
	}

	err = os.WriteFile(outputFile, []byte(base64Encoded), 0644)
	if err != nil {
		log.Fatalf("Error writing base64 output file: %v", err)
	}

	return
}

func (i *Image) Decode() (string, error) {
	// 解码Base64编码为字节数据
	decodedData, err := base64.StdEncoding.DecodeString(string(*i))
	if err != nil {
		log.Fatalf("Error decoding base64 data: %v", err)
	}
	return string(decodedData), nil
}

func (i *Image) Decode2File(filePath string) error {
	decodedData, err := i.Decode()
	if err != nil {
		log.Fatalf("Error decoding base64 data: %v", err)
	}

	err = os.WriteFile(filePath, []byte(decodedData), 0644)
	if err != nil {
		log.Fatalf("Error writing output image file: %v", err)
	}

	return nil
}
