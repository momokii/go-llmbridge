package openai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (c *openaiAPI) OpenAITextToSpeech(req_body *OAReqTextToSpeech) (*OATextToSpeechResp, error) {

	// ----------- input checker request
	if req_body.Model == "" || (req_body.Model != "tts-1" && req_body.Model != "tts-1-hd") {
		return nil, errors.New("Model must be gpt-3 or davinci")
	}

	if req_body.Input == "" {
		return nil, errors.New("Input text must be provided")
	}

	if req_body.Voice != "" && (req_body.Voice != "alloy" && req_body.Voice != "echo" && req_body.Voice != "fable" && req_body.Voice != "onyx" && req_body.Voice != "nova" && req_body.Voice != "shimer") {
		return nil, errors.New("Voice must be alloy, echo, fable, onyx, nova, or shimmer")
	}

	if req_body.ResponseFormat != "" && (req_body.ResponseFormat == "mp3" && req_body.ResponseFormat == "opus" && req_body.ResponseFormat == "aac" && req_body.ResponseFormat == "flac" && req_body.ResponseFormat == "wav" && req_body.ResponseFormat == "pcm") {
		return nil, errors.New("ResponseFormat must be mp3, opus, aac, flac, wav, or pcm")
	}

	if req_body.Speed != nil && (*req_body.Speed < 0.25 || *req_body.Speed > 4.0) {
		return nil, errors.New("Speed must be between 0.25 and 4.0")
	}

	apiKey := c.apiKey
	if apiKey == "" {
		return nil, errors.New("API Key is empty")
	}

	// create json ver for req body
	reqBodyJson, err := json.Marshal(req_body)
	if err != nil {
		return nil, errors.New("Failed to marshal request body")
	}

	// create req
	req, err := http.NewRequest(http.MethodPost, OAUrlTextToSpeech, bytes.NewBuffer(reqBodyJson))
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

	// decode file mp3 response to encode base64
	// because from the docs will be return file extension for audio, so for the response will be base64 encoded version of the audio we received
	var b64audio, fileExt string
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read response body3: " + err.Error())
	}

	b64audio = base64.StdEncoding.EncodeToString(fileBytes)

	if req_body.ResponseFormat == "" {
		fileExt = ".mp3"
	} else {
		fileExt = "." + req_body.ResponseFormat
	}

	result := OATextToSpeechResp{
		B64JSON:     b64audio,
		FormatAudio: fileExt,
	}

	return &result, nil
}
