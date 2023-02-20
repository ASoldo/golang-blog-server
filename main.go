package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type Post struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Likes    int       `json:"likes"`
	Comments []Comment `json:"comments"`
}

var posts []Post

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/{id}/like", likePost).Methods("POST")
	router.HandleFunc("/posts/{id}/comments", createComment).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.ID = len(posts) + 1
	posts = append(posts, post)
	json.NewEncoder(w).Encode(post)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, post := range posts {
		if post.ID == id {
			json.NewEncoder(w).Encode(post)
			return
		}
	}
	http.NotFound(w, r)
}

func likePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, post := range posts {
		if post.ID == id {
			posts[i].Likes++
			json.NewEncoder(w).Encode(posts[i])
			return
		}
	}
	http.NotFound(w, r)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, post := range posts {
		if post.ID == id {
			comment.ID = len(post.Comments) + 1
			posts[i].Comments = append(posts[i].Comments, comment)
			json.NewEncoder(w).Encode(posts[i])
			return
		}
	}
	http.NotFound(w, r)
}
