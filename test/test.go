package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/controller"
	"urlshortner/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShortTheUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := types.ShortUrlBody{
		LongUrl: "https://example.com",
	}
	jsonBody, _ := json.Marshal(body)

	c.Request, _ = http.NewRequest("POST", "/url/short", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.ShortUrl(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	assert.False(t, response["error"].(bool))
	assert.NotEmpty(t, response["data"])
	assert.NotEmpty(t, response["short_url"])
}
