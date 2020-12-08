package robots

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

// ImageRobot handles image search on google images
type ImageRobot struct{}

// Start ImageRobot
func (robot *ImageRobot) Start(state *State) {
	log.Println("[image] => Starting...")

	robot.fetchImageOfAllSentences(state.Sentences, state.SearchTerm)

	log.Println("[image] => Done, Tchau :P")
}

func (robot *ImageRobot) fetchImageOfAllSentences(sentences []*Sentence, searchTerm string) []*Sentence {
	cx := os.Getenv("GOOGLE_CUSTOM_SEARCH_ENGINE_ID")
	customSearchService := robot.createGoogleCustomSearchAPIService()
	search := customSearchService.Cse.List().Cx(cx).SearchType("image").Num(2)

	for index, sentence := range sentences {
		var searchQuery string
		if index == 0 {
			searchQuery = searchTerm
		} else {
			searchQuery = searchTerm + " " + sentence.Keywords[0]
		}

		log.Printf("\n[image] => Querying images with: %s", searchQuery)
		resp, err := search.Q(searchQuery).Do()
		if err != nil {
			log.Fatalf("\n[image] => Error on google search %v", err)
		}

		for _, image := range resp.Items {
			sentence.Images = append(sentence.Images, image.Link)
		}
	}

	return sentences
}

func (robot *ImageRobot) createGoogleCustomSearchAPIService() *customsearch.Service {
	googleAPIKey := os.Getenv("GOOGLE_API_KEY")
	client := &http.Client{Transport: &transport.APIKey{Key: googleAPIKey}}
	svc, err := customsearch.New(client)
	if err != nil {
		log.Fatalf("\n[image] => Error trying to create google custom search service %v", err)
	}
	return svc
}
