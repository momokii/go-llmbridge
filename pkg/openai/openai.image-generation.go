package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (c *openaiAPI) OpenAICreateImageDallE(req_body *OAReqImageGeneratorDallE) (*OAImageGeneratorDallEResp, error) {

	// ----------- input checker request
	if req_body.Model == "" || (req_body.Model != "dall-e-2" && req_body.Model != "dall-e-3") {
		return nil, errors.New("Model must be dall-e-2 or dall-e-3")
	}

	if req_body.N != nil && (*req_body.N < 1 || *req_body.N > 10) {
		return nil, errors.New("N must be between 1 and 10")
	}

	if req_body.Model != "dall-e-3" && req_body.Quality != nil {
		return nil, errors.New("Quality is only supported for dall-e-3 model")
	}

	if req_body.Quality != nil && (*req_body.Quality != "standard" && *req_body.Quality != "hd") {
		return nil, errors.New("Quality must be standard or hd")
	}

	if req_body.Model != "dall-e-3" && req_body.Style != nil {
		return nil, errors.New("Style is only supported for dall-e-3 model")
	}

	if req_body.Style != nil && (*req_body.Style != "vivid" && *req_body.Style != "natural") {
		return nil, errors.New("Style must be vivid or natural")
	}

	if req_body.ResponseFormat != nil && (*req_body.ResponseFormat != "url" && *req_body.ResponseFormat != "b64_json") {
		return nil, errors.New("ResponseFormat must be url or b64_json")
	}

	apiKey := c.apiKey
	if apiKey == "" {
		return nil, errors.New("API Key is empty")
	}

	reqBodyJson, err := json.Marshal(req_body)
	if err != nil {
		return nil, errors.New("Failed to marshal request body")
	}

	// create and send request
	req, err := http.NewRequest(http.MethodPost, OAUrlImageGenerationsDallE, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return nil, errors.New("Failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := c.config.httpClient

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to send request: " + err.Error())
	}
	defer func() {
		if resp.StatusCode != http.StatusOK {
			io.ReadAll(resp.Body)
		}
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to send request: " + resp.Status)
	}

	var respDataDallE OAImageGeneratorDallEResp
	if err := json.NewDecoder(resp.Body).Decode(&respDataDallE); err != nil {
		return nil, errors.New("Failed to decode response: " + err.Error())
	}

	return &respDataDallE, nil
}
