package main

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func readFile(path string) string {
	f, err := ioutil.ReadFile(path)
	check(err)

	return string(f)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}