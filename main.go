package main

import (
	"fmt"
	"os"
)

func main() {
	msg, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer msg.Close()
	bytes := make([]byte, 8)
	for {
		_, err = msg.Read(bytes)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", bytes)
	}
	os.Exit(0)
}
