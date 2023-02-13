package controllers

import (
	"forbes2023/data"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// Regex to remove non alphanumeric characters and - and '
var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9'\- ]+`)

// Format post story request bodies are in
type PostStoryRequest struct {
	Story string `json:"story" binding:"required"`
}

// Individual unmatched word for a story
type UnmatchedWord struct {
	Word       string `json:"word"`
	CloseMatch string `json:"closeMatch"`
}

// Concurrent unmatched words list with a lock to mutate words parameter
type concurrentUnmatchedWords struct {
	sync.RWMutex
	words []UnmatchedWord
}

// Append locks access to the unmatchedWords instance and updates its items
func (unmatchedWords *concurrentUnmatchedWords) Append(word UnmatchedWord) {
	unmatchedWords.Lock()
	defer unmatchedWords.Unlock()

	unmatchedWords.words = append(unmatchedWords.words, word)
}

// PostStory accepts the post request from the OpenAPI3 spec and notes what words aren't in the dictionary
func PostStory(c *gin.Context) {
	var input PostStoryRequest
	if c.BindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	words := make(map[string]bool)
	storyWords := strings.Fields(
		strings.ToLower(nonAlphanumericRegex.ReplaceAllString(strings.Clone(input.Story), "")),
	)

	for _, word := range storyWords {
		if _, exists := words[word]; !exists {
			words[word] = true
		}
	}

	var unmatchedWords concurrentUnmatchedWords
	dictionary, fileReadSuccess := data.BuildDictionaryFromFile()
	if !fileReadSuccess {
		c.Data(http.StatusInternalServerError, gin.MIMEJSON, nil)
		return
	}

	var closestWordsLookup sync.WaitGroup
	for word, _ := range words {
		if !dictionary.Contains(word) {
			closestWordsLookup.Add(1)
			go func() {
				defer closestWordsLookup.Done()
				unmatchedWords.Append(UnmatchedWord{
					Word:       word,
					CloseMatch: dictionary.NearestMatch(word),
				})
			}()
		}
	}

	closestWordsLookup.Wait()

	c.JSON(http.StatusOK, gin.H{"story": input.Story, "unmatchedWords": unmatchedWords.words})
}
