package robots

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

var (
	GOOGLE_API_KEY string
	CX             string
)

// ImageRobot handles image search on google images
type ImageRobot struct{}

// Start ImageRobot
func (robot *ImageRobot) Start(state *State) error {
	logrus.Info("🖼 [IMAGE] => Starting...")

	GOOGLE_API_KEY = os.Getenv("GOOGLE_API_KEY")
	CX = os.Getenv("GOOGLE_CUSTOM_SEARCH_ENGINE_ID")

	robot.fetchImageOfAllSentences(state.Sentences, state.SearchTerm)
	return nil
}

func (robot *ImageRobot) fetchImageOfAllSentences(sentences []*Sentence, searchTerm string) ([]*Sentence, error) {
	customSearchService, _ := robot.createGoogleCustomSearchAPIService()
	search := customSearchService.Cse.List().Cx(CX).SearchType("image").Num(2)

	for index, sentence := range sentences {
		var searchQuery string
		if index == 0 {
			searchQuery = searchTerm
		} else {
			searchQuery = searchTerm + " " + sentence.Keywords[0]
		}

		logrus.Info("🖼 [IMAGE] => Querying images with: " + searchQuery)
		resp, _ := search.Q(searchQuery).Do()
		for _, image := range resp.Items {
			sentence.Images = append(sentence.Images, image.Link)
		}
	}

	return sentences, nil
}

func (robot *ImageRobot) createGoogleCustomSearchAPIService() (*customsearch.Service, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: GOOGLE_API_KEY}}
	svc, err := customsearch.New(client)
	if err != nil {
		return nil, err
	}
	return svc, nil
}
