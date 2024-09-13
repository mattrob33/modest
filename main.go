package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
    "os"

	"github.com/gin-gonic/gin"
)

type SlackResponse struct {
    ResponseType string `json:"response_type"`
    Text         string `json:"text"`
}

func main() {
	r := gin.Default()

	// Route to render the form where user can input a prompt
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

    // Route to handle the form submission and make the API request
    r.GET("/generate", func(c *gin.Context) {
        // Get the prompt from the form
        prompt := c.Query("prompt")

        // Call the function to generate the article using the API
        article, err := generateArticle(prompt)
        if err != nil {
            log.Printf("Error generating article: %v", err)
            c.String(http.StatusInternalServerError, "Error generating article")
            return
        }

        // Serve the raw HTML content
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(article))
    })

	// Route to handle the form submission and make the API request
	r.POST("/generate-api", func(c *gin.Context) {
		// Get the prompt from the form
		title := c.PostForm("text")

		// Call the function to generate the article using the API
		preview, err := generatePreview(title)
		if err != nil {
			log.Printf("Error generating article: %v", err)
			c.String(http.StatusInternalServerError, "Error generating article")
			return
		}

        c.String(http.StatusOK, "*" + title + ":* " + preview);
	})

	// Route to handle the form submission and make the API request
	r.POST("/slack", func(c *gin.Context) {
		// Get the prompt from the form
		text := c.PostForm("text")
		command := c.PostForm("command")
		userId := c.PostForm("user_id")

		c.String(200, text + " ~ " + command + " ~ " + userId)
	})

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Start the Gin server
	r.Run(getPort())
}

// Function to generate the article by calling the API
func generateArticle(prompt string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=AIzaSyAB22gdjZYFhFdRO3qnSODsXwA-Sz0Qpgw"
	
	// Prepare the request body
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "Generate a satirical news article with the following title: `" + prompt + "`. Output in HTML instead of markdown. Style it like a real news website."},
				},
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Error marshalling JSON: %v", err)
	}

	// Create and send the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	// Parse the response JSON
	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling response: %v", err)
	}

	// Return the generated text
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("No content found")
}

// Function to generate the article by calling the API
func generatePreview(prompt string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=AIzaSyAB22gdjZYFhFdRO3qnSODsXwA-Sz0Qpgw"
	
	// Prepare the request body
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "Return the preview text for a satirical news article with the following title: `" + prompt + "`. Just a couple sentences.'"},
				},
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Error marshalling JSON: %v", err)
	}

	// Create and send the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	// Parse the response JSON
	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling response: %v", err)
	}

	// Return the generated text
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return response.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("No content found")
}


// Get the Port from the environment so we can run on Heroku
func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
