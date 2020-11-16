package robots

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"

	algorithmiaAPI "github.com/algorithmiaio/algorithmia-go"
	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"

	"github.com/sirupsen/logrus"
)

// TextRobot handles NLP and search Wikipedia based on searchTerm and prefix from state
type TextRobot struct{}

type algorithmSearchInput struct {
	Search string `json:"search"`
	Lang   string `json:"lang,omitempty"`
}

type wikipediaSearchResult struct {
	Content string `json:"content"`
}

// Start TextRobot
func (robot *TextRobot) Start(state *State) error {
	logrus.Info("📜 [TEXT] => Starting...")

	wikipediaSearchAlgorithmResult, err := robot.fetchContentFromWikipedia(state.SearchTerm)
	if err != nil {
		return err
	}

	sanitizedContent := robot.sanitizeContent(wikipediaSearchAlgorithmResult.Content)

	sentences, err := robot.splitIntoSentences(sanitizedContent)
	if err != nil {
		return err
	}

	state.SourceContentOriginal = wikipediaSearchAlgorithmResult.Content
	state.SourceContentSanitized = sanitizedContent
	state.Sentences = sentences

	test, _ := json.Marshal(sentences[1:10])
	logrus.Info(string(test))
	return nil
}

func (robot *TextRobot) fetchContentFromWikipedia(searchTerm string) (*wikipediaSearchResult, error) {
	logrus.Info("📜 [TEXT] => Fetching content from Wikipedia")

	wikipediaSearchRawContent, err := robot.fetchWikipediaSearchAndParseResult(searchTerm)
	if err != nil {
		return nil, err
	}
	contentMapped, err := robot.unmarshalWikipediaResponse(wikipediaSearchRawContent)
	if err != nil {
		return nil, err
	}

	logrus.Info("📜 [TEXT] => Fetching done! ")
	return contentMapped, nil
}

func (robot *TextRobot) fetchWikipediaSearchAndParseResult(searchTerm string) (map[string]interface{}, error) {
	algorithmiaAPIKey := os.Getenv("ALGORITHMIA_API_KEY")
	algorithmia := algorithmiaAPI.NewClient(algorithmiaAPIKey, "")

	wikipediaAlgo, err := algorithmia.Algo("web/WikipediaParser/0.1.2")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("📜 [TEXT] => " + ErrAlgorithmInstatiation.Error())
		return nil, ErrAlgorithmInstatiation
	}

	searchInputRaw, err := json.Marshal(&algorithmSearchInput{Search: searchTerm})
	JSONSearchInput := string(searchInputRaw)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("📜 [TEXT] => Error on JSON Marshalling, aborting...")
		return nil, err
	}

	wikipediaSearchRawResponse, err := wikipediaAlgo.Pipe(JSONSearchInput)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("📜 [TEXT] => " + ErrAlgorithm.Error())
		return nil, ErrAlgorithm
	}

	return wikipediaSearchRawResponse.(*algorithmiaAPI.AlgoResponse).Result.(map[string]interface{}), nil
}

func (robot *TextRobot) unmarshalWikipediaResponse(rawResponse map[string]interface{}) (*wikipediaSearchResult, error) {
	JSONResponse, err := json.Marshal(rawResponse)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("📜 [TEXT] => Error on JSON Marshalling, aborting...")
		return nil, err
	}
	unmarshalledResponse := &wikipediaSearchResult{}

	err = json.Unmarshal(JSONResponse, unmarshalledResponse)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("📜 [TEXT] => Error on JSON Unmarshalling, aborting...")
		return nil, err
	}

	return unmarshalledResponse, nil
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

func (robot *TextRobot) splitIntoSentences(sourceContentSanitized string) ([]string, error) {
	sentencesTokenizer, err := robot.instantiateSentencesTokenizer()
	if err != nil {
		return nil, err
	}

	sentences := sentencesTokenizer.Tokenize(sourceContentSanitized)
	sentencesText := []string{}

	for _, s := range sentences {
		sentencesText = append(sentencesText, s.Text)
	}

	return sentencesText, nil
}

func (robot *TextRobot) instantiateSentencesTokenizer() (*sentences.DefaultSentenceTokenizer, error) {
	// Compiling language specific data into a binary file can be accomplished
	// by using `make <lang>` and then loading the `json` data:
	trainingData, err := data.Asset("data/english.json")
	if err != nil {
		return nil, err
	}

	// load the training data
	training, err := sentences.LoadTraining(trainingData)
	if err != nil {
		return nil, err
	}
	// create the default sentence tokenizer
	return sentences.NewSentenceTokenizer(training), nil
}