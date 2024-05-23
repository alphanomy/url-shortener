package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ShortlyLink struct {
	gorm.Model
	OriginalUrl  string `gorm:"unique"`
	ShortenedUrl string `gorm:"unique"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		url.QueryEscape(dbUsername),
		url.QueryEscape(dbPassword),
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&ShortlyLink{})

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/shorten", func(c *gin.Context) {
		var data struct {
			Url string `json:"url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var link ShortlyLink
		result := db.Where("original_url = ?", data.Url).First(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				shortenedUrl := generateStrings(6)
				link = ShortlyLink{OriginalUrl: data.Url, ShortenedUrl: shortenedUrl}
				result = db.Create(&link)
				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"shortened_url": link.ShortenedUrl})
	})

	r.GET("/:shortenedUrl", func(c *gin.Context) {
		shortenedUrl := c.Param("shortenedUrl")
		var link ShortlyLink
		result := db.Where("shortened_url = ?", shortenedUrl).First(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			return
		}

		c.Redirect(http.StatusMovedPermanently, link.OriginalUrl)
	})

	r.Run(":8000")
}

func generateStrings(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var shortenedUrl string
	for i := 0; i < length; i++ {
		shortenedUrl += string(chars[rand.Intn(len(chars))])
	}

	return shortenedUrl
}
