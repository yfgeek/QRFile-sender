package main

import (
	"QRFileSender/src/core"
	"QRFileSender/src/utils"
	"time"
)


const TmpFilePath = "./out" //
const FileChunkSize = 1024

func main() {
	defer utils.Timer(time.Now())
	qs := core.NewQRFileSender("./1.zip",FileChunkSize,TmpFilePath)
	qs.Start()
}
