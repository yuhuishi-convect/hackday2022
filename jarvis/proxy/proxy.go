package main

import (
	"context"
	"jarvis/customerrors"
	"jarvis/db"
	"jarvis/ent"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

var mu sync.Mutex

func handle(f func(c *gin.Context) error) gin.HandlerFunc {
	return func(context *gin.Context) {
		if err := f(context); err != nil {
			if ae, ok := err.(customerrors.AppError); ok {
				context.JSON(ae.StatusCode, gin.H{
					"message": ae.ErrorText,
				})
			} else {
				log.Println(err.Error())
				context.JSON(500, gin.H{
					"message": "Internal server error",
				})
			}
		}

	}
}

func main() {
	router := gin.Default()
	router.GET("/uploadSchema", uploadSchema)

	router.GET("/api/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome!",
		})
	})
	router.GET("/api/:typeName", handle(db.FetchAll))
	router.POST("/api/:typeName", handle(db.Create))
	router.GET("/api/:typeName/:id", handle(db.Fetch))
	// router.PATCH("/api/:typeName/:id", handle(controller.UpdateOneReadingItem))
	router.DELETE("/api/:typeName/:id", handle(db.Remove))

	router.Run("localhost:8080")
}

func uploadSchema(c *gin.Context) {
	mu.Lock()
	client, err := ent.Open("sqlite3", "file:ent.db?mode=rwc&cache=shared&_fk=1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed connecting to database")
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed updating schema in database")
		log.Fatalf("failed creating schema resources: %v", err)
	}
	mu.Unlock()
}
