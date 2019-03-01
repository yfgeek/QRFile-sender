package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"encoding/base64"
	"github.com/skip2/go-qrcode"
)

const FileChunkSize = 1024
const TmpFilePath = "./out" //

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
func CleanTmpFolder(){
	os.RemoveAll(TmpFilePath)
}

func NewTmpFolder(folderPath string){
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


func SplitFile(fileToBeChunked string){
	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()


	fmt.Print(FileChunkSize)

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(FileChunkSize)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(1); i < totalPartsNum; i++ {

		partSize := int(math.Min(FileChunkSize, float64(fileSize-int64(i*FileChunkSize))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		fileName := "file_" + strconv.FormatUint(i, 10)
		filePath := TmpFilePath + "/" + fileName + ".png"
		_, err := os.Create(filePath)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// base64
		encodeString := base64.StdEncoding.EncodeToString(partBuffer)

		fmt.Println(encodeString)

		qrcode.WriteFile(encodeString,qrcode.Low,1024, filePath)

		fmt.Println("Split to : ", fileName)
	}
}


func main() {
	CleanTmpFolder()
	NewTmpFolder(TmpFilePath)
	SplitFile("./a.txt")

}