package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Usage struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Lang string `json::"language"`
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	Def string `json:"definition"`
	Examples []string `json:"examples"`
}

func ParseHTML(htmlText string) (string) {
	doc, err := html.Parse(strings.NewReader(htmlText))
	if err != nil {
		log.Fatal(err)
	}
	return parseHTML(doc)
}

func parseHTML(n *html.Node) (string) {
	if n.Type == html.TextNode {
		return n.Data
	} else {
		plain := ""
		for c := n.FirstChild; c!= nil; c = c.NextSibling {
			plain += parseHTML(c)
		}
		return plain
	}
}

func ParseDefinitions(json []interface{}) (definitions []Definition) {
	for _, value := range json {
		var definition Definition
		switch typ := value.(type) {
		case map[string]interface{}:
			plain_def := ParseHTML(typ["definition"].(string))
			definition.Def = plain_def
			if typ["examples"] != nil {
				switch ex := typ["examples"].(type) {
				case []interface{}:
					for _, s := range ex {
						plain_example := ParseHTML(s.(string))
						definition.Examples = append(definition.Examples, plain_example)
					}
				default:
					fmt.Println("Error: some other typ")
				}
			}
		default:
			fmt.Println("Some other type", value)
		}
		definitions = append(definitions, definition)
	}
	return
}

func ParseUsages(json map[string]interface{}) (usages []Usage) {
	for key, value := range json {
		var usage Usage
		switch typ := value.(type) {
		case []interface{}:
			for _, u := range typ {
				switch v := u.(type) {
				case map[string]interface{}:
					usage.PartOfSpeech = v["partOfSpeech"].(string)
					usage.Lang = v["language"].(string)
					usage.Definitions = ParseDefinitions(v["definitions"].([]interface{}))
				default:
					fmt.Println("Some other type", v)
				}
			}
		default:
			fmt.Println(key, "is some other type")
		}
		usages = append(usages, usage)
	}
	return
}

func main() {
	if len(os.Args) == 1 {
	    panic("You must be looking for some word")
	}
	word := os.Args[1]
	response, err := http.Get("https://en.wiktionary.org/api/rest_v1/page/definition/" + word)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}


	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var result interface{}

	err = json.Unmarshal([]byte(responseData), &result)

	parsedResponse := result.(map[string]interface{})

	usages := ParseUsages(parsedResponse)
	for _, usage := range usages {
		fmt.Println(usage.Lang)
		fmt.Println("Part of Speech:", usage.PartOfSpeech)
		fmt.Println("Definitions:")
		for i, definition := range usage.Definitions {
			fmt.Print(i + 1, ". ")
			fmt.Println(definition.Def)
			if definition.Examples != nil {
				fmt.Println("Examples")
				for _, example := range definition.Examples {
					fmt.Print("- ")
					fmt.Println(example)
				}
			}
		}
		fmt.Println()
	}
}
