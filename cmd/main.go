package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/estevam31/slide-maker/pkg/robots"
	"github.com/joho/godotenv"
)

// State shared between the robots
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	Start()
}

// Start the orchestrator
func Start() {
	state := &robots.State{}

	robotUserInput := &robots.UserInputRobot{}
	robotImage := &robots.ImageRobot{}
	robotText := &robots.TextRobot{}
	robotSlides := &robots.SlidesRobot{}

	robotUserInput.Start(state)
	robotText.Start(state)
	robotImage.Start(state)
	robotSlides.Start(state)

	json, _ := json.MarshalIndent(state, "", " ")
	_ = ioutil.WriteFile("result.json", json, 0644)
}
