package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Details string `json:"details"`
	Done    bool   `json:"done"`
}

var todolists = []Todo{
	{ID: 1, Title: "ทำการบ้าน", Details: "การบ้านวิชา Database", Done: false},
	{ID: 2, Title: "อ่านหนังสือ", Details: "อ่านบทที่ 3 วิชา Network", Done: false},
}

func getTodos(c *gin.Context) {
	doneQuery := strings.TrimSpace(c.Query("done"))
	titleQuery := strings.TrimSpace(c.Query("title"))

	filtered := []Todo{}

	for _, todo := range todolists {
		match := true

		if doneQuery != "" {
			if doneQuery == "true" && !todo.Done {
				match = false
			}
			if doneQuery == "false" && todo.Done {
				match = false
			}
		}

		if titleQuery != "" && todo.Title != titleQuery {
			match = false
		}

		if match {
			filtered = append(filtered, todo)
		}
	}

	if doneQuery != "" || titleQuery != "" {
		c.IndentedJSON(http.StatusOK, filtered)
		return
	}

	c.IndentedJSON(http.StatusOK, todolists)
}

func addTodo(c *gin.Context) {
	var newTodo Todo

	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doneQuery := strings.TrimSpace(c.Query("done"))

	if doneQuery != "" {
		if doneQuery == "true" && !newTodo.Done {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "todo ต้อง Done=true"})
			return
		}
		if doneQuery == "false" && newTodo.Done {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "todo ต้อง Done=false"})
			return
		}
	}

	newTodo.ID = len(todolists) + 1
	todolists = append(todolists, newTodo)

	c.IndentedJSON(http.StatusCreated, newTodo)
}

func main() {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/todolists", getTodos)
		api.POST("/todolists", addTodo)
	}

	r.Run(":8080")
}
