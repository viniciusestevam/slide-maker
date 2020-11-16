package main

import (
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
	robotText := &robots.TextRobot{}

	robotUserInput.Start(state)
	robotText.Start(state)
}
