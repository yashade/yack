package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func readFile(path string) string {
	if _, err := os.Stat(path); err == nil { // TODO: avoid manual implementation
		f, err2 := ioutil.ReadFile(path)
		check(err2)
		return string(f)
	} else {
		return ""
	}
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}