package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Usage is a struct storing the information for each usage in a Wiktionary entry.
// Each word can have several usages in the same or different languages.
type Usage struct {
	// Part of speech of the word, e.g. Noun, Adjective, Interjection
	PartOfSpeech string `json:"partOfSpeech"`
	// Language for this usage
	Lang string `json::"language"`
	// List of definitions for this usage
	Definitions []Definition `json:"definitions"`
}

// Definition is a struct storing information of a definition of a word usage
type Definition struct {
	// The definition
	Def string `json:"definition"`
	// The examples for this definition
	Examples []string `json:"examples"`
}

func parseHTML(htmlText string) (string) {
	doc, err := html.Parse(strings.NewReader(htmlText))
	if err != nil {
		log.Fatal(err)
	}
	return parseDocTree(doc)
}

func parseDocTree(n *html.Node) (string) {
	if n.Type == html.TextNode {
		return n.Data
	}
	plain := ""
	for c := n.FirstChild; c!= nil; c = c.NextSibling {
		plain += parseDocTree(c)
	}
	return plain
}

func parseDefinitions(json []interface{}) (definitions []Definition) {
	for _, value := range json {
		var definition Definition
		switch typ := value.(type) {
		case map[string]interface{}:
			plainDef := parseHTML(typ["definition"].(string))
			definition.Def = plainDef
			if typ["examples"] != nil {
				switch ex := typ["examples"].(type) {
				case []interface{}:
					for _, s := range ex {
						plainExample := parseHTML(s.(string))
						definition.Examples = append(definition.Examples, plainExample)
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

func parseUsages(json map[string]interface{}, language string) (usages []Usage) {
	for key, value := range json {
		if language != "*" && key != language {
			continue
		}
		var usage Usage
		switch typ := value.(type) {
		case []interface{}:
			for _, u := range typ {
				switch v := u.(type) {
				case map[string]interface{}:
					usage.PartOfSpeech = v["partOfSpeech"].(string)
					usage.Lang = v["language"].(string)
					usage.Definitions = parseDefinitions(v["definitions"].([]interface{}))
					usages = append(usages, usage)
				default:
					fmt.Println("Some other type", v)
				}
			}
		default:
			fmt.Println(key, "is some other type")
		}
	}
	return
}

func main() {
	langPtr := flag.String("l", "en", "2-letter code for the language you want to search,\nuse \"*\" to include all language")
	wordPtr := flag.String("w", "default", "The word you want to look up")
	flag.Parse()
	word := *wordPtr
	fmt.Println("Definition for", word)
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

	if parsedResponse["title"] == "Not found." {
		fmt.Println("That word does not exist.")
		os.Exit(0)
	}

	usages := parseUsages(parsedResponse, *langPtr)
	var currentLang string
	for _, usage := range usages {
		if currentLang != usage.Lang {
			fmt.Println(usage.Lang)
			currentLang = usage.Lang
		}
		fmt.Println("\tPart of Speech:", usage.PartOfSpeech)
		fmt.Println("\tDefinitions:")
		for i, definition := range usage.Definitions {
			fmt.Print("\t", i + 1, ". ")
			fmt.Println(definition.Def)
			if definition.Examples != nil {
				fmt.Println("\t", "Examples")
				for _, example := range definition.Examples {
					fmt.Print("\t", "- ")
					fmt.Println(example)
				}
			}
		}
		fmt.Println()
	}
}
