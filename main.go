package main

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB
var err error

type Todo struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	IsDone  bool   `json:"isDone"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	USERNAME := os.Getenv("USERNAME")
	PASSWORD := os.Getenv("PASSWORD")
	DATABASE := os.Getenv("DATABASE")
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")

	db, err = gorm.Open("mysql", USERNAME+":"+PASSWORD+"@tcp("+HOST+":"+PORT+")/"+DATABASE)
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&Todo{})

	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.POST("/todos", createTodo)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run(":8080")
}

func getTodos(c *gin.Context) {
	var todos []Todo

	err := db.Find(&todos).Error
	if err != nil {
		c.AbortWithError(404, errors.New("No data"))
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"toDos": todos,
	})
}

func getTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")

	err := db.Where("id = ?", id).First(&todo).Error
	if err != nil {
		c.AbortWithError(500, errors.New("ID not found"))
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"toDo": todo,
	})
}

func createTodo(c *gin.Context) {
	var todo Todo

	c.BindJSON(&todo)

	db.Create(&todo)
	c.JSON(200, gin.H{
		"toDo": todo,
	})
}

func updateTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")

	err := db.Where("id = ?", id).First(&todo).Error
	if err != nil {
		c.AbortWithError(404, errors.New("ID not found"))
		panic(err.Error())
	}

	c.BindJSON(&todo)

	db.Save(&todo)
	c.JSON(200, gin.H{
		"toDo": todo,
	})
}

func deleteTodo(c *gin.Context) {
	var todo Todo
	id := c.Params.ByName("id")

	err := db.Where("id = ?", id).Delete(&todo).Error
	if err != nil {
		c.AbortWithError(404, errors.New("ID not found"))
		panic(err.Error())
	}

	c.JSON(200, gin.H{
		"message": "ID = " + id + " is deleted",
	})
}
