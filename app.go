package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Post struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func main() {
	db, err := gorm.Open("sqlite3", "./db.sqlite3")
	check(err)
	defer db.Close()

	db.AutoMigrate(&Post{})

	router := httprouter.New()

	router.GET("/", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("hello yack"))
	})

	// get all posts
	router.GET("/posts", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		offset := r.URL.Query().Get("offset")
		limit := r.URL.Query().Get("limit")

		var posts []Post
		if offset != "" && limit != "" {
			db.Limit(limit).Offset(offset).Find(&posts)
		} else {
			db.Find(&posts) // get all posts. ALL OF THEM!
		}

		postsJson, err := json.Marshal(posts)
		check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(postsJson))
	})

	// get single post
	router.GET("/posts/:id", func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		idParam := p.ByName("id")
		var post Post

		db.First(&post, idParam)

		postJson, err := json.Marshal(post)
		check(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(postJson))
	})

	router.POST("/posts", func (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		check(err)

		// not a good solution but i don't wanna unmarshal the key to a struct
		var dataTmp map[string]interface{}
		if isJSON(string(body)) {
			err2 := json.Unmarshal(body, &dataTmp)
			check(err2)

			if dataTmp["key"] == readFile("key") {
				var postData Post
				json.Unmarshal(body, &postData)
				post := Post{Title: postData.Title, Content: postData.Content}
				db.Create(&post)
			} else {
				w.WriteHeader(401)
				w.Write([]byte("Unauthorized"))
			}
		} else {
			w.WriteHeader(400)
			w.Write([]byte("Bad request"))
		}
	})

	http.ListenAndServe(":2001", router)
}
