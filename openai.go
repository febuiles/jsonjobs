package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []Message       `json:"messages"`
	Function string          `json:"function,omitempty"`
	Params   json.RawMessage `json:"params,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

type OpenAIClient struct {
	APIKey   string
	Endpoint string
}

const JSONInstruction = `Take the job entry text and return a job entry in JSON format adhering to this schema:
	{
	  "company_title": "string",
	  "location": "string",
	  "url": "string",
	  "job_title": "string"
	}.`

func NewOpenAIClient(apiKey, endpoint string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:   apiKey,
		Endpoint: endpoint,
	}
}

func (c OpenAIClient) ParseEntry(body string) string {
	requestBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You read job entries and return output in the specified JSON schema.",
			},
			{
				Role:    "user",
				Content: JSONInstruction,
			},
			{
				Role:    "user",
				Content: "Job Entry: \n" + body,
			},
			{
				Role:    "user",
				Content: "Ignore jobs based in the US",
			},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var openAIResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		panic(err)
	}

	var res string
	for _, choice := range openAIResponse.Choices {
		res = choice.Message.Content
	}

	return res
}
