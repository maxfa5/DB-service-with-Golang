package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	file, err := os.Create("loger.txt")
	if err != nil {
		fmt.Println("error in logger")
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))
}
