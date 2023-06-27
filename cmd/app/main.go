package main

import (
	"forum/internal/app"
	"log"
)

const cfgFilePath = "configs/config.json"

func main() {
	if err := app.Run(cfgFilePath); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
