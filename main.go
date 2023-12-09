package main

import (
	"fmt"
	"net/http"
	"port_scanner/utils/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	router.GET("/port", func(c *gin.Context) {
		domain := c.Query("host")

		if domain == "" {
			c.JSON(400, gin.H{
				"message": "Domain is required",
			})
			return
		}

		openPortsJson, err := handlers.ScanPortsMain(domain)
		if err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"message": "Error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": openPortsJson,
		})
	})

	router.POST("/directories", func(c *gin.Context) {

		var formData struct {
			URL string `form:"url" binding:"required"`
		}

		// ShouldBindWith binds the form data from the request body into the struct
		if err := c.ShouldBindWith(&formData, binding.Form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"error at formFile": err.Error(),
			})
			return
		}

		uploadPath := "./uploads/" + file.Filename

		err = c.SaveUploadedFile(file, uploadPath)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"error at saving file": err.Error(),
			})
			return
		}

		fmt.Println(uploadPath)
		availableDirectories := handlers.HandleURL(formData.URL, uploadPath)

		c.JSON(http.StatusOK, gin.H{
			"availableDirectories": availableDirectories,
		})
	})

	router.Run()
}
