package robots

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

var (
	GOOGLE_API_KEY = os.Getenv("GOOGLE_API_KEY")
)

// ImageRobot handles image search on google images
type ImageRobot struct{}

// Start ImageRobot
func (robot *ImageRobot) Start(state *State) error {
	logrus.Info("🖼 [IMAGE] => Starting...")
	return nil
}

func (robot *ImageRobot) createGoogleCustomSearchAPIService() (*customsearch.Service, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: GOOGLE_API_KEY}}

	svc, err := customsearch.New(client)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (robot *ImageRobot) fetchImageOfAllSentences(sentences []string) {

}
