package utils

import (
	"fmt"
	"time"
)

func Log(messages ...interface{}) {
	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s\n", currentTime, fmt.Sprint(messages...))
}
