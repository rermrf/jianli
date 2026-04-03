package main

import (
	"log"

	"jianli/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server bootstrap placeholder on :%s", cfg.Port)
}
