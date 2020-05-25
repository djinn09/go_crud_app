package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var router *chi.Mux
var db *sql.DB

const (
	dbName = "YOUR_DB_NAME"
	dbPass = "********"
	dbHost = "localhost"
	dbPort = "3306"
)

func routers() *chi.Mux {
	// router.Get("/posts", AllPosts)
	// router.Get("/posts/{id}", DetailPost)
	router.Post("/posts/create", CreatePost)
	router.Put("/posts/{id}", UpdatePost)
	router.Delete("/posts/{id}", DeletePost)

	return router
}

func init() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)
	dbSource := fmt.Sprintf("dev:%s@tcp(%s:%s)/%s?charset=utf8", dbPass, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", dbSource)
	catch(err)
}

// Post : for database method
type Post struct {
	ID      int `json: "id"`
	Title   int `json: "title"`
	Content int `json: "content"`
}

// CreatePost create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	json.NewDecoder(r.Body).Decode(&post)

	query, err := db.Prepare("Insert posts Set title=?,content=?")
	catch(err)

	_, er := query.Exec(post.Title, post.Content)
	catch(er)

	defer query.Close()
	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "succesfully created"})
}

// UpdatePost is update
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&post)

	query, err := db.Prepare("Update posts set title=?,content=? where id=?")
	catch(err)
	_, er := query.Exec(post.Title, post.Content, id)
	catch(er)
	defer query.Close()
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "update sucessfully"})
}

// DeletePost is delete
func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query, err := db.Prepare("delete from posts where id=?")
	catch(err)
	_, er := query.Exec(id)

	catch(er)
	query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "successfully deleted"})

}

func main() {
	var PORT = ":8005"
	fmt.Println("Server runnning on " + PORT)
	routers()
	http.ListenAndServe(PORT, Logger())

}
