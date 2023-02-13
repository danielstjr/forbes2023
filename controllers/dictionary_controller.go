package controllers

import (
	"forbes2023/data"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetDictionary lists all items in the dictionary file exactly as they appear in file
func GetDictionary(c *gin.Context) {
	dictionary, fileReadSuccess := data.BuildDictionaryFromFile()
	if !fileReadSuccess {
		c.JSON(http.StatusOK, gin.H{"error": "Unable to retrieve data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dictionary": dictionary.ToSlice()})
}

type PostDictionaryAdd struct {
	Add []string `json:"add"`
}

type PostAddDictionaryEntryRequest struct {
	Dictionary PostDictionaryAdd `json:"dictionary" binding:"required"`
}

// PostDictionary adds an item to the dictionary file
func PostDictionary(c *gin.Context) {
	dictionary, fileReadSuccess := data.BuildDictionaryFromFile()
	if !fileReadSuccess {
		c.Data(http.StatusInternalServerError, gin.MIMEJSON, nil)
		return
	}

	var input PostAddDictionaryEntryRequest
	if err := c.BindJSON(&input); err != nil || input.Dictionary.Add == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	for _, entry := range input.Dictionary.Add {
		if dictionary.Contains(entry) {
			c.JSON(http.StatusOK, gin.H{"dictionary": input.Dictionary, "error": "duplicate entry"})
			return
		}

		dictionary.Add(entry)
	}

	if !data.SaveDictionaryToFile(dictionary) {
		c.Data(http.StatusInternalServerError, gin.MIMEJSON, nil)
	} else {
		c.JSON(http.StatusAccepted, nil)
	}
}

type DeleteDictionaryRemove struct {
	Remove []string `json:"remove"`
}

type DeleteDictionaryEntryRequest struct {
	Dictionary DeleteDictionaryRemove `json:"dictionary" binding:"required"`
}

// DeleteDictionary removes the word from the dictionary file
func DeleteDictionary(c *gin.Context) {
	dictionary, fileReadSuccess := data.BuildDictionaryFromFile()
	if !fileReadSuccess {
		c.Data(http.StatusInternalServerError, gin.MIMEJSON, nil)
		return
	}

	var input DeleteDictionaryEntryRequest
	if c.BindJSON(&input) != nil || input.Dictionary.Remove == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	for _, entry := range input.Dictionary.Remove {
		if !dictionary.Contains(entry) {
			c.JSON(http.StatusOK, gin.H{"dictionary": input.Dictionary, "error": "entry not found"})
			return
		}

		dictionary.Remove(entry)
	}

	if !data.SaveDictionaryToFile(dictionary) {
		c.Data(http.StatusInternalServerError, gin.MIMEJSON, nil)
	} else {
		c.JSON(http.StatusAccepted, nil)
	}
}
