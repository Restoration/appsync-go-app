package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type Message struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

var db *pg.DB

func main() {
	// PostgreSQLの設定
	db = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "user",
		Password: "password",
		Database: "chatdb",
	})

	r := gin.Default()

	r.GET("/messages", fetchMessages)
	r.POST("/send", sendMessage)

	r.Run(":8080")
}

func fetchMessages(c *gin.Context) {
	var messages []Message
	err := db.Model(&messages).Select()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, messages)
}

func sendMessage(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, insertErr := db.Model(&msg).Returning("id").Insert()
	if insertErr != nil {
		c.JSON(500, gin.H{"error": insertErr.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Message sent successfully", "data": msg})
}
