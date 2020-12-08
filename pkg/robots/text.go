package robots

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	algorithmiaAPI "github.com/algorithmiaio/algorithmia-go"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"

	"github.com/IBM/go-sdk-core/core"
	nlu "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
)

// TextRobot handles NLU and search Wikipedia based on searchTerm and prefix from state
type TextRobot struct{}

type algorithmSearchInput struct {
	Search string `json:"search"`
	Lang   string `json:"lang,omitempty"`
}

type wikipediaSearchResult struct {
	Content string `json:"content"`
}

// Start TextRobot
func (robot *TextRobot) Start(state *State) {
	log.Printf("[text] => Starting...")

	wikipediaSearchAlgorithmResult := robot.fetchContentFromWikipedia(state.SearchTerm)
	sanitizedContent := robot.sanitizeContent(wikipediaSearchAlgorithmResult.Content)

	maxSentences := 10
	sentences := robot.splitIntoSentences(sanitizedContent)
	sentencesSliced := sentences[0:maxSentences]
	sentencesWithKeywords := robot.fetchWatsonAndReturnKeywords(sentencesSliced)

	state.SourceContentOriginal = wikipediaSearchAlgorithmResult.Content
	state.SourceContentSanitized = sanitizedContent
	state.Sentences = sentencesWithKeywords

	log.Printf("[text] => Done, goodbye ^^")
}

func (robot *TextRobot) fetchContentFromWikipedia(searchTerm string) *wikipediaSearchResult {
	log.Printf("[text] => Fetching content from Wikipedia...")

	wikipediaSearchRawContent := robot.fetchWikipediaSearchAndParseResult(searchTerm)
	contentMapped := robot.unmarshalWikipediaResponse(wikipediaSearchRawContent)

	log.Printf("[text] => Fetch done!")
	return contentMapped
}

func (robot *TextRobot) fetchWikipediaSearchAndParseResult(searchTerm string) map[string]interface{} {
	wikipediaAlgorithmKey := "web/WikipediaParser/0.1.2"
	algorithmiaAPIKey := os.Getenv("ALGORITHMIA_API_KEY")

	algorithmia := algorithmiaAPI.NewClient(algorithmiaAPIKey, "")
	wikipediaAlgo, err := algorithmia.Algo(wikipediaAlgorithmKey)
	if err != nil {
		log.Fatalf("[text] => Error trying to instantiate algorithmia Wikipedia parser", err)
	}

	searchInputRaw, _ := json.Marshal(&algorithmSearchInput{Search: searchTerm})
	JSONSearchInput := string(searchInputRaw)

	wikipediaSearchRawResponse, err := wikipediaAlgo.Pipe(JSONSearchInput)
	if err != nil {
		log.Fatalf("[text] => Error on wikipedia search algorithm", err)
	}

	return wikipediaSearchRawResponse.(*algorithmiaAPI.AlgoResponse).Result.(map[string]interface{})
}

func (robot *TextRobot) unmarshalWikipediaResponse(rawResponse map[string]interface{}) *wikipediaSearchResult {
	JSONResponse, err := json.Marshal(rawResponse)
	unmarshalledResponse := &wikipediaSearchResult{}
	err = json.Unmarshal(JSONResponse, unmarshalledResponse)

	if err != nil {
		log.Fatalf("[text] => Error on wikipedia response unmarshalling")
	}
	return unmarshalledResponse
}

func (robot *TextRobot) sanitizeContent(originalContent string) string {
	withoutBlankLinesAndMarkdown := robot.removeBlankLinesAndMarkdown(originalContent)
	withoutDatesInParentheses := robot.removeDatesInParentheses(withoutBlankLinesAndMarkdown)
	return withoutDatesInParentheses
}

func (robot *TextRobot) removeBlankLinesAndMarkdown(text string) string {
	allLines := strings.Split(text, "\n")

	filterStringArray := func(stringArray []string, test func(string) bool) (result []string) {
		for _, item := range stringArray {
			if test(item) {
				result = append(result, strings.TrimSpace(item))
			}
		}
		return result
	}
	blankLinesFilter := func(line string) bool {
		trimmedLine := strings.TrimSpace(line)
		if startsWithEqual := strings.HasPrefix(trimmedLine, "="); startsWithEqual == true || len(trimmedLine) == 0 {
			return false
		}
		return true
	}

	return strings.Join(filterStringArray(allLines, blankLinesFilter), " ")
}

func (robot *TextRobot) removeDatesInParentheses(text string) string {
	removeDates := regexp.MustCompile("/\\((?:\\([^()]*\\)|[^()])*\\)/gm")
	removeSpaceDuplicates := regexp.MustCompile("/  /g")
	text = removeDates.ReplaceAllString(text, "")
	text = removeSpaceDuplicates.ReplaceAllString(text, " ")
	return text
}

func (robot *TextRobot) splitIntoSentences(sourceContentSanitized string) []*Sentence {
	sentencesTokenizer := robot.createSentencesTokenizer()

	tokenizedText := sentencesTokenizer.Tokenize(sourceContentSanitized)
	sentencesWithText := []*Sentence{}

	for _, s := range tokenizedText {
		sentencesWithText = append(sentencesWithText, &Sentence{Text: s.Text})
	}

	return sentencesWithText
}

func (robot *TextRobot) createSentencesTokenizer() *sentences.DefaultSentenceTokenizer {
	// Compiling language specific data into a binary file can be accomplished
	// by using `make <lang>` and then loading the `json` data:
	trainingData, _ := data.Asset("data/english.json")

	// load the training data
	training, err := sentences.LoadTraining(trainingData)
	if err != nil {
		log.Fatalf("[text] => Error on loading sentence tokenizer training", err)
	}
	// create the default sentence tokenizer
	return sentences.NewSentenceTokenizer(training)
}

func (robot *TextRobot) fetchWatsonAndReturnKeywords(sentences []*Sentence) []*Sentence {
	log.Printf("[text] => Analyzing content and recognizing keywords...")
	nluSvc := robot.createWatsonNLUService()

	for _, sentence := range sentences {
		analyzeOptions := nluSvc.NewAnalyzeOptions(&nlu.Features{
			Keywords: &nlu.KeywordsOptions{},
		}).SetText(sentence.Text)

		analyzeResult, _, _ := nluSvc.Analyze(analyzeOptions)
		for _, keywordResult := range analyzeResult.Keywords {
			sentence.Keywords = append(sentence.Keywords, *keywordResult.Text)
		}
	}

	log.Printf("[text] => Keywords recognized")
	return sentences
}

func (robot *TextRobot) createWatsonNLUService() *nlu.NaturalLanguageUnderstandingV1 {
	// Instantiate the Watson Natural Language Understanding service
	authenticator := &core.IamAuthenticator{
		ApiKey: os.Getenv("WATSON_API_KEY"),
	}
	service, err := nlu.
		NewNaturalLanguageUnderstandingV1(&nlu.NaturalLanguageUnderstandingV1Options{
			URL:           os.Getenv("WATSON_NLU_URL"),
			Version:       "2017-02-27",
			Authenticator: authenticator,
		})

	if err != nil {
		log.Fatalf("[text] => Error creating NLU service", err)
	}

	return service
}
