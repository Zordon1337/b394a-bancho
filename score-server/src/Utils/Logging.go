package Utils

import (
	"fmt"
	"time"
)

func LogErr(message string, args ...interface{}) {
	fmt.Printf("\033[31m["+time.Now().String()+"] [ERROR] "+message+"\n\033[0m\n", args...)
}
func LogInfo(message string, args ...interface{}) {
	fmt.Printf("\033[34m["+time.Now().String()+"] [INFO] "+message+"\n\033[0m\n", args...)
}
func LogWarning(message string, args ...interface{}) {
	fmt.Printf("\033[33m["+time.Now().String()+"] [WARNING] "+message+"\n\033[0m\n", args...)
}
