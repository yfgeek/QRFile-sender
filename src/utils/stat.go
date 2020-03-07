package utils

import (
	"fmt"
	"time"
)

func Timer(start time.Time){
	terminal:=time.Since(start)
	fmt.Println(terminal)
}
