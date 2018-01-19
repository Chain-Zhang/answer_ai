package main

import (
	"time"
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
			start := time.Now()
			ai.Start()
			elapsed := time.Since(start)
			fmt.Println("本次答题耗时: ", elapsed)
		case "2":
            ai.ExeCommand("cmd", []string{"/c", "adb", "devices"})
		case "exit":
			os.Exit(1)
		}
	}
}