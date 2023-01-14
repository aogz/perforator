package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ResponseChoice is a choice returned in OpenAI response
type ResponseChoice struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

// Response from OpenAI
type Response struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Model   string           `json:"model"`
	Choices []ResponseChoice `json:"choices"`
}

// DaVinciRequest explains commit messages using Chat GPT API
func DaVinciRequest(text string) (string, error) {
	payload := map[string]interface{}{
		"model":             "text-davinci-003",
		"prompt":            text,
		"temperature":       0,
		"max_tokens":        60,
		"top_p":             1,
		"frequency_penalty": 0.5,
		"presence_penalty":  0,
	}
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/completions",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_SECRET_TOKEN")))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid status code: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Text, nil
}
