package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/theghostmac/todo-api-with-gin/internal/handlers"
	"github.com/theghostmac/todo-api-with-gin/internal/repository"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection setup
	connStr := "user=youruser password=yourpassword dbname=yourdbname sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	router := gin.Default()

	// Instantiate dependencies
	todoRepository := repository.NewTodoRepository(db)
	todoHandler := handlers.NewTodoHandler(todoRepository)

	// Register routes
	router.GET("/todos", todoHandler.GetAllTodos)
	router.GET("/todos/:id", todoHandler.GetTodoByID)
	router.POST("/todos", todoHandler.CreateTodo)
	router.PUT("/todos/:id", todoHandler.UpdateTodo)
	router.DELETE("/todos/:id", todoHandler.DeleteTodo)

	// Start the server
	router.Run()
}
