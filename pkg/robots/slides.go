package robots

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kjk/betterguid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/slides/v1"
)

// SlidesRobot creates and publishes the presentation on google slides
type SlidesRobot struct{}

// Start SlidesRobot
func (robot *SlidesRobot) Start(state *State) {
	fmt.Println("[slides] => Starting...")
	robot.createPresentation(state.SearchTerm, state.Sentences)
	fmt.Println("[slides] => Done, Aloha :D")
}

func (robot *SlidesRobot) createPresentation(title string, sentences []*Sentence) *slides.Presentation {
	fmt.Println("[slides] => Creating presentation...")
	slidesService := robot.createGoogleSlidesAPIService()

	presentation, err := slidesService.Presentations.Create(&slides.Presentation{Title: title}).Fields("presentationId").Do()
	if err != nil {
		log.Fatalf("\n[slides] => Could not create presentation on google slides %v", err)
	}

	createSlidesRequests := robot.createRequestsForSlidePages(sentences, presentation.PresentationId)
	requestBody := &slides.BatchUpdatePresentationRequest{
		Requests: createSlidesRequests,
	}
	_, err = slidesService.Presentations.BatchUpdate(presentation.PresentationId, requestBody).Do()
	if err != nil {
		log.Fatalf("\n[slides] => Error on create slides requests %v", err)
	}

	fmt.Println("[slides] => Presentation created")
	return presentation
}

func (robot *SlidesRobot) createRequestsForSlidePages(sentences []*Sentence, presentationID string) []*slides.Request {
	fmt.Println("[slides] => Creating slide pages...")
	requests := []*slides.Request{}
	for _, sentence := range sentences {
		slideID := "S" + betterguid.New()

		createSlideBody := robot.createSlideBodyRequest(slideID, 0)
		insertText := robot.insertTextRequests(sentence.Text, slideID, 16, 350, 350, 100)
		insertImage := robot.insertImageRequest(sentence.Images[0], slideID, 100000.0, 100000.0)

		requests = append(requests, createSlideBody)
		requests = append(requests, insertText...)
		requests = append(requests, insertImage)
	}
	return requests
}

func (robot *SlidesRobot) createSlideBodyRequest(slideID string, index int64) *slides.Request {
	return &slides.Request{
		CreateSlide: &slides.CreateSlideRequest{
			ObjectId:       slideID,
			InsertionIndex: index,
			SlideLayoutReference: &slides.LayoutReference{
				PredefinedLayout: "BLANK",
			},
		},
	}
}

func (robot *SlidesRobot) insertTextRequests(text string, slideID string, fontSize float64, boxSize float64, translateX float64, translateY float64) []*slides.Request {
	textBoxID := "T" + betterguid.New()
	boxSizeDimension := slides.Dimension{
		Magnitude: boxSize,
		Unit:      "PT",
	}
	fontSizeDimension := slides.Dimension{
		Magnitude: fontSize,
		Unit:      "PT",
	}
	createTextBoxRequest := &slides.Request{
		CreateShape: &slides.CreateShapeRequest{
			ObjectId:  textBoxID,
			ShapeType: "TEXT_BOX",
			ElementProperties: &slides.PageElementProperties{
				PageObjectId: slideID,
				Size: &slides.Size{
					Height: &boxSizeDimension,
					Width:  &boxSizeDimension,
				},
				Transform: &slides.AffineTransform{
					ScaleX:     1.0,
					ScaleY:     1.0,
					TranslateX: translateX,
					TranslateY: translateY,
					Unit:       "PT",
				},
			},
		},
	}

	insertTextOnTextboxRequest := &slides.Request{
		InsertText: &slides.InsertTextRequest{
			ObjectId:       textBoxID,
			InsertionIndex: 0,
			Text:           text,
		},
	}
	updateTextStyleRequest := &slides.Request{
		UpdateTextStyle: &slides.UpdateTextStyleRequest{
			ObjectId: textBoxID,
			Style: &slides.TextStyle{
				FontSize:   &fontSizeDimension,
				FontFamily: "Montserrat",
			},
			Fields: "fontSize, fontFamily",
		},
	}

	return []*slides.Request{createTextBoxRequest, insertTextOnTextboxRequest, updateTextStyleRequest}
}

func (robot *SlidesRobot) insertImageRequest(imageURL string, slideID string, translateX float64, translateY float64) *slides.Request {
	imageID := "I" + betterguid.New()
	dimension := slides.Dimension{
		Magnitude: 4000000,
		Unit:      "EMU",
	}
	return &slides.Request{
		CreateImage: &slides.CreateImageRequest{
			ObjectId: imageID,
			Url:      imageURL,
			ElementProperties: &slides.PageElementProperties{
				PageObjectId: slideID,
				Size: &slides.Size{
					Height: &dimension,
					Width:  &dimension,
				},
				Transform: &slides.AffineTransform{
					ScaleX:     1.0,
					ScaleY:     1.0,
					TranslateX: translateX,
					TranslateY: translateY,
					Unit:       "EMU",
				},
			},
		},
	}
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
