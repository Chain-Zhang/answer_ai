package main

import (
	"fmt"
	"answer_ai/ai"
	"os"
)

func main(){
	for {
		var cmd string
		fmt.Printf("> ")
		fmt.Scan(&cmd)
		switch cmd{
		case "1":
			ai.Start()
		case "2":
            ai.ExeCommand("cmd", []string{"/c", "adb", "devices"})
		case "exit":
			os.Exit(1)
		}
	}
}