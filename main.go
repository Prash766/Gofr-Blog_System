package main

import (
	
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
	"log"
	"net/http"
	"strconv"
	"time"
)

//post represents a blog post
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var posts []Post // Slice to hold blog posts
var postIDCounter = 1

func main() {
	app := gofr.New()

	// Endpoints for blog posts
	app.POST("/posts", createPost)
	app.GET("/posts", getAllPosts)
	app.GET("/posts/:id", getPostByID)
	app.PUT("/posts/:id", updatePost)
	app.DELETE("/posts/:id", deletePost)

	log.Println("Server running on port 8000...")
	app.Start()
}

func createPost(ctx *gofr.Context) (interface{}, error) {
	var post Post
	if err := ctx.Bind(&post); err != nil {
		return nil, err
	}

	post.ID = postIDCounter
	post.CreatedAt = time.Now()

	posts = append(posts, post)
	postIDCounter++

	return post, nil
}

func getAllPosts(ctx *gofr.Context) (interface{}, error) {
	return posts, nil
}

func getPostByID(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Reason:     "Invalid post ID",
		}
	}

	for _, post := range posts {
		if post.ID == id {
			return post, nil
		}
	}

	return nil, &errors.Response{
		StatusCode: http.StatusNotFound,
		Reason:     "Post not found",
	}
}

func updatePost(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Reason:     "Invalid post ID",
		}
	}

	var updatedPost Post
	if err := ctx.Bind(&updatedPost); err != nil {
		return nil, err
	}

	for i, post := range posts {
		if post.ID == id {
			updatedPost.ID = id
			updatedPost.CreatedAt = post.CreatedAt
			posts[i] = updatedPost
			return updatedPost, nil
		}
	}

	return nil, &errors.Response{
		StatusCode: http.StatusNotFound,
		Reason:     "Post not found",
	}
}

func deletePost(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Reason:     "Invalid post ID",
		}
	}

	for i, post := range posts {
		if post.ID == id {
			posts = append(posts[:i], posts[i+1:]...)
			return "Post deleted successfully", nil
		}
	}

	return nil, &errors.Response{
		StatusCode: http.StatusNotFound,
		Reason:     "Post not found",
	}
}
