package robots

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
)

// UserInputRobot handles user's input
type UserInputRobot struct{}

// Start UserInputRobot
func (robot *UserInputRobot) Start(state *State) {
	searchTerm := robot.askForSearchTerm()
	prefix := robot.askForPrefix()

	state.SearchTerm = searchTerm
	state.Prefix = prefix

	log.Println("[user_input] => Successfully requested user's input")
	log.Printf("[user_input] => Search term: %s\n", searchTerm)
	log.Printf("[user_input] => Prefix: %s\n", prefix)

	log.Println("[user_input] => Done, adiós xD")
}

func (robot *UserInputRobot) askForSearchTerm() string {
	log.Println("Type a Wikipedia search term:")
	searchTerm, err := robot.readline()
	if err != nil {
		log.Fatalf("\n[user_input] => Error asking for search term %v", err)
	}
	return searchTerm
}

func (robot *UserInputRobot) askForPrefix() string {
	log.Println("Select one option:")
	prefixes := [3]string{"Who is", "What is", "The history of"}
	for i, prefix := range prefixes {
		log.Printf("[%d] - %s\n", i+1, prefix)
	}

	selected, err := robot.readline()
	if err != nil {
		log.Fatalf("\n[user_input] => Error asking for prefix %v", err)
	}

	prefixIndex, err := strconv.Atoi(selected)
	prefixIndex = prefixIndex - 1
	if prefixIndex > len(prefixes) || prefixIndex < 0 || err != nil {
		log.Println("✖ Invalid option, please try again.")
		robot.askForPrefix()
	}

	return prefixes[prefixIndex]
}

func (robot *UserInputRobot) readline() (string, error) {
	log.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", errors.New("✖ Error reading user's input")
}
