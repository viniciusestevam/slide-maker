package robots

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/slides/v1"
)

// SlidesRobot creates and publishes the presentation on google slides
type SlidesRobot struct{}

// Start SlidesRobot
func (robot *SlidesRobot) Start(state *State) {
	fmt.Println("[slides] => Starting...")
	robot.createSlide(state.SearchTerm)
	fmt.Println("[slides] => Done, Aloha :D")
}

func (robot *SlidesRobot) createSlide(title string /* sentences []*Sentence*/) *slides.Presentation {
	fmt.Println("[slides] => Creating slide...")
	slidesService := robot.createGoogleSlidesAPIService()
	presentation, err := slidesService.Presentations.Create(&slides.Presentation{Title: title}).Do()
	if err != nil {
		log.Fatalf("\n[slides] => Could not create presentation on google slides %v", err)
	}
	return presentation
}

func (robot *SlidesRobot) createGoogleSlidesAPIService() *slides.Service {
	fmt.Println("[slides] => Asking for user permission...")
	credentialsFile, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("\n[slides] => Could not find credentials file %v", err)
	}
	config, err := google.ConfigFromJSON(credentialsFile, "https://www.googleapis.com/auth/presentations")
	if err != nil {
		log.Fatalf("\n[slides] => Could not create client config from credentials file %v", err)
	}
	client := robot.getHTTPClient(config)

	slidesService, err := slides.New(client)
	if err != nil {
		log.Fatalf("\n[slides] => Unable to create google slides service %v", err)
	}

	return slidesService
}

func (robot *SlidesRobot) getHTTPClient(config *oauth2.Config) *http.Client {
	token := robot.requestUserAuthenticationToken(config)
	return config.Client(context.Background(), token)
}

func (robot *SlidesRobot) requestUserAuthenticationToken(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("[slides] => Go to the following link in your browser then type the authorization code:")
	fmt.Printf("%v\n\n>", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		fmt.Printf("\n[slides] => Could not read user authorization code %v", err)
		fmt.Println("[slides] => Please try again")
		robot.requestUserAuthenticationToken(config)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("\n[slides] => Unable to retrieve token from web: %v", err)
	}
	return token
}
