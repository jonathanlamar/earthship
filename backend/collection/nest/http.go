package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func sendEmptyPost(url string) ([]byte, error) {
	postBody, _ := json.Marshal(map[string]string{})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		return nil, errors.New("Failed to send POST.")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read response body.")
	}

	return body, nil
}

func sendGetRequestWithAccessToken(url, accessToken string) ([]byte, error) {
	// Empty body.  Header has to be set after the request is created.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Failed to initialize request.")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to send GET.")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read response body.")
	}

	return body, nil
}
