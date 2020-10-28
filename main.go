package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

type Entry struct {
    PartOfSpeech string `json:"partOfSpeech"`
    Language string `json::"language"`
    Definitions []Definition `json:"definitions"`
}

type Definition struct {
    Def string `json:"definition"`
    examples []string `json:"examples"`
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

    for key, value := range entries {
        switch typ := value.(type) {
        case string:
            fmt.Println(key, "is string", typ)
        case int:
            fmt.Println(key, "is int", typ)
        case []interface{}:
            fmt.Println(key, "is an array:")
            for i, u := range typ {
                fmt.Println(i, u)
            }
        default:
            fmt.Println(key, "is some other type")
        }
    }
}
