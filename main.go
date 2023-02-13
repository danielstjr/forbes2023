package main

import (
	"forbes2023/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("dictionary", controllers.GetDictionary)
	router.POST("dictionary", controllers.PostDictionary)
	router.DELETE("dictionary", controllers.DeleteDictionary)

	router.POST("story", controllers.PostStory)

	router.Run("localhost:8080")
}
