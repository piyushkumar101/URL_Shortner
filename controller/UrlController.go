package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"urlshortner/constant"
	"urlshortner/database"
	"urlshortner/helper"
	"urlshortner/types"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ShortUrl(c *gin.Context) {
	var shortUrlBody types.ShortUrlBody
	err := c.BindJSON(&shortUrlBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": constant.BindError})
		return
	}
	cachedURL, err := RedisClient.Get(shortUrlBody.LongUrl).Result()
	if err == nil {
		// URL found in cache we are returning cached URL
		c.JSON(http.StatusOK, gin.H{"error": false, "data": cachedURL})
		return
	}
	code := helper.GenRandomString(6)

	record, _ := database.Mgr.GetUrlFromCode(code, constant.UrlCollection)

	if record.UrlCode != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "this code is already in use"})
		return
	}

	var url types.UrlDb

	url.CreatedAt = time.Now().Unix()
	url.ExpiredAt = time.Now().Unix()
	url.UrlCode = code
	url.LongUrl = shortUrlBody.LongUrl
	url.ShortUrl = constant.BaseUrl + code
	err = RedisClient.Set(shortUrlBody.LongUrl, url.ShortUrl, 24*time.Hour).Err()
	if err != nil {
		log.Println("Not found in redis cache")
	}

	resp, err := database.Mgr.Insert(url, constant.UrlCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "data": resp, "short_url": url.ShortUrl})
}

func RedirectURL(c *gin.Context) {
	code := c.Param("code")
	cachedURL, err := RedisClient.Get(code).Result()
	if err == nil {
		// Redirect to cached URL
		c.Redirect(http.StatusPermanentRedirect, cachedURL)
		return
	}
	record, _ := database.Mgr.GetUrlFromCode(code, constant.UrlCollection)

	if record.UrlCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "there is no url found"})
		return
	}
	err = RedisClient.Set(code, record, 24*time.Hour).Err()
	if err != nil {
		log.Println("url not found in redis cache")
	}
	fmt.Println(record.LongUrl)

	c.Redirect(http.StatusPermanentRedirect, record.LongUrl)
}
