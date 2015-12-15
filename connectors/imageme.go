package connectors

import (
	"encoding/json"
	"fmt"
	"github.com/projectjane/jane/models"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type ImageMe struct {
}

func (x ImageMe) Listen(commandMsgs chan<- models.Message, connector models.Connector) {
	defer Recovery(connector)
	return
}

func (x ImageMe) Command(message models.Message, publishMsgs chan<- models.Message, connector models.Connector) {
	if strings.Index(message.In.Text, "image me") == 0 {
		msg := strings.TrimSpace(strings.Replace(message.In.Text, "image me", "", 1))
		message.Out.Text = callImageMe(msg, connector.Key, connector.Pass, false)
		publishMsgs <- message
	}
	if strings.Index(message.In.Text, "animate me") == 0 {
		msg := strings.TrimSpace(strings.Replace(message.In.Text, "animate me", "", 1))
		message.Out.Text = callImageMe(msg, connector.Key, connector.Pass, true)
		publishMsgs <- message
	}
}

func (x ImageMe) Publish(connector models.Connector, message models.Message, target string) {
	return
}

func (x ImageMe) Help(connector models.Connector) (help string) {
	help += "image me <image keywords> - pulls back an image url\n"
	help += "animate me <image keywords> - pulls back an animated gif url\n"
	return help
}

type SearchResult struct {
	Items []Items `json:"items"`
}

type Items struct {
	Link string `json:"link"`
}

var client = &http.Client{}

var baseUrl = "https://www.googleapis.com/customsearch/v1?key="
var errorMessage = "Error retrieving image"
var animated bool

func callImageMe(msg string, apiKey string, cx string, animated bool) string {
	start := rand.Intn(3)
	if start < 1 {
		start = 1
	}

	cx = "&cx=" + cx
	returnFields := fmt.Sprintf("&fields=items(link)&start=%v", start)
	query := "&q=" + url.QueryEscape(msg)
	fields := "&searchType=image"
	if animated {
		fields += "&fileType=gif&hq=animated&tbs=itp:animated"
	}

	url := baseUrl + apiKey + cx + returnFields + query + fields

	resp, err := client.Get(url)
	if err != nil {
		log.Print(err)
		return errorMessage
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		FindDeprecatedImage(msg, animated)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return errorMessage
	}

	var result SearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Print(err)
		return errorMessage
	}

	if len(result.Items) > 0 {
		randomLink := result.Items[rand.Intn(len(result.Items))]
		return randomLink.Link
	} else {
		return FindDeprecatedImage(msg, animated)
	}
}

type DeprecatedResult struct {
	ResponseData ResponseData `json:"responseData"`
}

type ResponseData struct {
	Results []Result `json:"results"`
}

type Result struct {
	Url string `json:"url"`
}

func FindDeprecatedImage(query string, animated bool) string {
	baseUrl := "https://ajax.googleapis.com/ajax/services/search/images?v=1.0&rsz=8"
	if animated {
		baseUrl += "&as_filetype=gif"
	}

	baseUrl += "&q="
	searchUrl := baseUrl + url.QueryEscape(query)

	if animated {
		searchUrl += url.QueryEscape(" animated")
	}

	resp, err := client.Get(searchUrl)
	if err != nil {
		log.Print(err)
		return errorMessage
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return errorMessage
	}

	var result DeprecatedResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Print(err)
		return errorMessage
	}

	index := rand.Intn(len(result.ResponseData.Results))

	if len(result.ResponseData.Results) > 0 {
		return result.ResponseData.Results[index].Url
	}

	return "No results found"
}