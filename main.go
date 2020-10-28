package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

type Usage struct {
    PartOfSpeech string `json:"partOfSpeech"`
    Language string `json::"language"`
    Definitions []Definition `json:"definitions"`
}

type Definition struct {
    Def string `json:"definition"`
    examples []string `json:"examples"`
}

func PrintJSON(json map[string]interface{}) {
    for key, value := range json {
        switch typ := value.(type) {
        case []interface{}:
            fmt.Println("Language:", key)
            for _, u := range typ {
                fmt.Println("{")
                switch v := u.(type) {
                case map[string]interface{}:
                    fmt.Println(v["partOfSpeech"])
                    fmt.Println(v["language"])
                    fmt.Println(v["definitions"])
                default:
                    fmt.Println("Some other type", v)
                }
                fmt.Println("},")
            }
            fmt.Println("]\n")
        default:
            fmt.Println(key, "is some other type", typ)
        }
    }
}

func main() {
    response, err := http.Get("https://en.wiktionary.org/api/rest_v1/page/definition/général")

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

    entries := result.(map[string]interface{})

    PrintJSON(entries)
}
