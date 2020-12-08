package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/estevam31/slide-maker/pkg/robots"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// State shared between the robots
func main() {
	logrus.SetFormatter(&logrus.TextFormatter{})

	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
		return
	}

	Start()
}

// Start the orchestrator
func Start() {
	state := &robots.State{}
	robotUserInput := &robots.UserInputRobot{}
	robotImage := &robots.ImageRobot{}
	robotText := &robots.TextRobot{}

	robotUserInput.Start(state)
	robotText.Start(state)
	robotImage.Start(state)

	json, _ := json.MarshalIndent(state, "", " ")
	_ = ioutil.WriteFile("result.json", json, 0644)
}
