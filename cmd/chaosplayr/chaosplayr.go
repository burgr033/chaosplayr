package main

import (
	"log"

	"github.com/burgr033/chaosplayr/internal/file"
	"github.com/burgr033/chaosplayr/internal/gui"
)

func main() {
	err := file.CreateFavorites()
	if err != nil {
		log.Println("failed to create Favorites File")
	}
	program, err := gui.CreateProgram("https://media.ccc.de/podcast-hq.xml")
	if err != nil {
		log.Fatalf("Error initializing program: %v\n", err)
	}
	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running program: %v\n", err)
	}
}
