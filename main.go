package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	msg, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer msg.Close()
	bytes := make([]byte, 8)
	var current_line strings.Builder
	for {
		_, err = msg.Read(bytes)
		current_bytes := strings.Split(string(bytes), "\n")
		if len(current_bytes) != 1 {
			current_line.WriteString(current_bytes[0])
			fmt.Printf("read: %s\n", current_line.String())
			current_line.Reset()
			current_line.WriteString(current_bytes[1])
		} else {
			current_line.WriteString(current_bytes[0])
		}
		if err != nil {
			break
		}

	}
	os.Exit(0)
}
