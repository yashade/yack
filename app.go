package main

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Post struct {
	Id int `json: "id"`
	Title string `json: "title"`
	Content string `json: "content"`
}

func main() {
	db, err := gorm.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&Post{})

	router := httprouter.New()

	router.GET("/", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("hello gratuscms"))
	})

	// get all posts
	router.GET("/posts", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var posts []Post
		db.Find(&posts)

		postsJson, err := json.Marshal(posts)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(postsJson))
	})

	// get single post
	router.GET("/posts/:id", func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		idParam := p.ByName("id")
		var post Post

		db.First(&post, idParam)

		postJson, err := json.Marshal(post)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(postJson))
	})

	router.POST("/posts", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Fatal(err)
			return
	}

		var postData Post
		json.Unmarshal(body, &postData)

		post := Post{Title: postData.Title, Content: postData.Content}
		db.Create(&post)
	})

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))


	http.ListenAndServe(":2001", router)
}