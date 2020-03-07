package core

import (
	"encoding/base64"
	"github.com/skip2/go-qrcode"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
)

type QRFileSender struct{
	filePath string
	fileOffset int
	tmpFilePath string
}

type FileBlock struct {
	id int
	content  string
	filePath string
}


func NewQRFileSender(filePath string, fileOffset int, tmpFilePath string) *QRFileSender{
	return &QRFileSender{
		filePath:   filePath ,
		fileOffset:  fileOffset,
		tmpFilePath: tmpFilePath,
	}
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

func (q *QRFileSender) cleanTmpFolder() error{
	err := os.RemoveAll(q.tmpFilePath)
	return err
}

func (q *QRFileSender) newTmpFolder() error{
	exist, err := PathExists(q.tmpFilePath)
	if err != nil {
		log.Printf("get dir error![%v]\n", err)
		return err
	}
	if exist {
		log.Printf("has dir![%v]\n", q.tmpFilePath)
	} else {
		log.Printf("no dir![%v]\n", q.tmpFilePath)
		err := os.Mkdir(q.tmpFilePath, os.ModePerm)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
			return err
		} else {
			log.Println("mkdir success!")
		}
	}
	return err
}

func (q *QRFileSender) Start(){
	err := q.cleanTmpFolder()
	if err!=nil{
		log.Fatal(err)
	}
	err = q.newTmpFolder()
	if err!=nil{
		log.Fatal(err)
	}
	err = q.splitFile()
	if err!=nil{
		log.Fatal(err)
	}
}

func  (q *QRFileSender) splitFile() error{
	// using go routine to speed up
	genQRchan := make(chan *FileBlock)

	file, err := os.Open(q.filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	fileInfo, _ := file.Stat()
	var fileSize = fileInfo.Size()
	log.Printf("File size: %d\n", fileSize)
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(q.fileOffset)))
	log.Printf("Splited to %d pieces.\n", totalPartsNum)

	pool := NewPool(100)
	log.Printf("The current routine nums :%d\n", runtime.NumGoroutine())
	for i := uint64(1); i <= totalPartsNum; i++ {
		pool.Add(1)
		partSize := int(math.Min(float64(q.fileOffset), float64(fileSize - int64(i-1) * int64(q.fileOffset))))
		partBuffer := make([]byte, partSize)
		file.Read(partBuffer)
		fileName := "file_" + strconv.FormatUint(i, 10)
		filePath := q.tmpFilePath + "/" + fileName + ".png"
		// base64
		fileBlock := &FileBlock{
			id : int(i),
			content:  base64.StdEncoding.EncodeToString(partBuffer),
			filePath: filePath,
		}
		go func() {
			GenerateQRCode(genQRchan)
			pool.Done()
		}()
		genQRchan <- fileBlock

	}
	pool.Wait()
	return err
}

func GenerateQRCode(c chan *FileBlock) {
	fileBlock := <-c
	fileId := fileBlock.id
	fileContent := fileBlock.content
	filePath := fileBlock.filePath
	s:= &SplitProtocol{
		id:      fileId,
		content: fileContent,
	}
	s.Encoder = Base64Encoder{}
	fileContent = s.String()

	_, err := os.Create(filePath)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = qrcode.WriteFile(fileContent, qrcode.Low, len(fileContent), filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("File %10d done", s.id)
}