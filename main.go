package main

import (
	"encoding/base64"
	"fmt"
	"github.com/skip2/go-qrcode"
	"math"
	"os"
	"strconv"
)

const FileChunkSize = 1024
const TmpFilePath = "./out" //

type FileBlock struct {
	content  string
	filePath string
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func CleanTmpFolder() {
	os.RemoveAll(TmpFilePath)
}

func NewTmpFolder(folderPath string) {
	exist, err := PathExists(folderPath)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("has dir![%v]\n", folderPath)
	} else {
		fmt.Printf("no dir![%v]\n", folderPath)
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}

func SplitFile(fileToBeChunked string) {
	genQRchan := make(chan *FileBlock)
	done := make(chan bool)

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize = fileInfo.Size()

	fmt.Printf("File size: %d\n", fileSize)

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(FileChunkSize)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(1); i <= totalPartsNum; i++ {

		go GenerateQRCode(genQRchan, done)

		partSize := int(math.Min(FileChunkSize, float64(fileSize-int64((i-1)*FileChunkSize))))

		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		fileName := "file_" + strconv.FormatUint(i, 10)
		filePath := TmpFilePath + "/" + fileName + ".png"

		// base64
		fileBlock := &FileBlock{
			content:  base64.StdEncoding.EncodeToString(partBuffer),
			filePath: filePath,
		}

		genQRchan <- fileBlock
	}
	for i := uint64(1); i <= totalPartsNum; i++ {
		<-done
	}

}

func GenerateQRCode(c chan *FileBlock, done chan bool) {
	fileBlock := <-c
	fileContent := fileBlock.content
	filePath := fileBlock.filePath
	_, err := os.Create(filePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Path: ", filePath, "Base64: ", fileContent)

	err = qrcode.WriteFile(fileContent, qrcode.Low, len(fileContent), filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	done <- true

}

func main() {
	CleanTmpFolder()
	NewTmpFolder(TmpFilePath)
	SplitFile("./1.zip")

}
