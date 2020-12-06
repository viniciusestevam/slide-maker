package robots

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// UserInputRobot handles user's input
type UserInputRobot struct{}

// Start UserInputRobot
func (robot *UserInputRobot) Start(state *State) error {
	searchTerm, prefix, err := robot.requestUserInput(state)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	state.SearchTerm = searchTerm
	state.Prefix = prefix

	logrus.WithFields(logrus.Fields{
		"searchTerm": state.SearchTerm,
		"prefix":     state.Prefix,
	}).Info("⌨ [USER_INPUT] => Successfully requested user's input")
	return nil
}

func (robot *UserInputRobot) requestUserInput(state *State) (string, string, error) {
	searchTerm, err := robot.askForSearchTerm()
	if err != nil {
		return "", "", err
	}

	prefix, err := robot.askForPrefix()
	if err != nil {
		return "", "", err
	}
	return searchTerm, prefix, nil
}

func (robot *UserInputRobot) askForSearchTerm() (string, error) {
	fmt.Println("⌨ Type a Wikipedia search term:")
	return robot.readline()
}

func (robot *UserInputRobot) askForPrefix() (string, error) {
	prefixes := [3]string{"Who is", "What is", "The history of"}
	fmt.Println("\n⌨ Select one option:")
	for i := 0; i < len(prefixes); i++ {
		fmt.Printf("[%d] - %s\n", i+1, prefixes[i])
	}

	selected, err := robot.readline()
	if err != nil {
		return "", err
	}

	prefixIndex, err := strconv.Atoi(selected)
	prefixIndex = prefixIndex - 1
	if prefixIndex > len(prefixes) || prefixIndex < 0 || err != nil {
		fmt.Println("✖ Invalid option, please try again.")
		robot.askForPrefix()
	}

	return prefixes[prefixIndex], nil
}

func (robot *UserInputRobot) readline() (string, error) {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", ErrReadline
}
