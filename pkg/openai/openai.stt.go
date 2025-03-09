package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func sttBaseReq(isTranscription, isWordStampReq, isSegmentStampReq bool, req_body OATranscriptionDefaultReq, APIKey string) ([]byte, error) {
	var stt_req OATranscriptionReq
	var req_url string // for checking if the request is for transcription or translation
	if isTranscription {
		req_url = OAUrlSTTTranscription
	} else {
		req_url = OAUrlSTTTranslation
	}

	// check user input validate base on api docs
	if req_body.File == nil {
		return nil, errors.New("file must be provided")
	}

	if req_body.Temperature != 0 && (req_body.Temperature < 0 || req_body.Temperature > 1) {
		return nil, errors.New("temperature must be between 0 and 1")
	}

	// checking file type input and extension. aldo parsing it to proper req struct
	var fileName string
	var fileContent io.Reader

	switch v := req_body.File.(type) {
	case *multipart.FileHeader:
		fileName = v.Filename
		var err error
		fileContent, err = v.Open()
		if err != nil {
			return nil, errors.New("failed to access file content: " + err.Error())
		}
		defer fileContent.(io.Closer).Close()
	case string:
		fileName = filepath.Base(v)
		var err error
		fileContent, err = os.Open(v)
		if err != nil {
			return nil, errors.New("failed to open file: " + err.Error())
		}
		defer fileContent.(io.Closer).Close()
	case io.Reader:
		fileName = req_body.Filename
		if fileName == "" {
			return nil, errors.New("filename must be provided if file is io.Reader")
		}

		fileContent = v
	default:
		return nil, errors.New("file type not supported, supported type is *multipart.FileHeader, string, or io.Reader")
	}

	fileExt := filepath.Ext(fileName)
	validExts := []string{".mp3", ".mp4", ".mpeg", ".mpga", ".m4a", ".webm", ".wav", ".flac", ".ogg"}

	isValid := false
	for _, ext := range validExts {
		if fileExt == ext {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("your file extension is " + fileExt + ", but it must be mp3, mp4, mpeg, mpga, m4a, webm, wav, flac, or ogg")
	}

	// parsing input to proper req struct
	stt_req = OATranscriptionReq{
		// File:  fileContent,
		Model: "whisper-1", // hard coded for now, because on openai docs only support this model
	}

	if req_body.Temperature != 0 {
		stt_req.Temperature = req_body.Temperature
	}

	if req_body.Prompt != "" {
		stt_req.Prompt = req_body.Prompt
	}

	if req_body.Language != "" {
		stt_req.Language = req_body.Language
	}

	// check the request if using word timestamps or segment timestamps or just default
	if isWordStampReq && isSegmentStampReq {
		return nil, errors.New("cannot use both word timestamps and segment timestamps")
	}

	// word timestamps
	if isWordStampReq {
		stt_req.ResponseFormat = "verbose_json"
		stt_req.TimestampGranularities = []string{"word"}
	}

	// segment timestamps
	if isSegmentStampReq {
		stt_req.ResponseFormat = "verbose_json"
		stt_req.TimestampGranularities = []string{"segment"}
	}

	// process form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// add file to form
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return nil, errors.New("failed to create form file")
	}

	if _, err = io.Copy(fw, fileContent); err != nil {
		return nil, errors.New("failed to copy file content")
	}

	// add other field
	// model field (required)
	if fw, err = w.CreateFormField("model"); err != nil {
		return nil, errors.New("failed to create form field for model")
	}
	if _, err := fw.Write([]byte(stt_req.Model)); err != nil {
		return nil, errors.New("failed to write model field")
	}

	// optional field, so do checking here
	if stt_req.Temperature != 0 {
		if fw, err = w.CreateFormField("temperature"); err != nil {
			return nil, errors.New("failed to create form field for temperature")
		}
		if _, err := fw.Write([]byte(strconv.FormatFloat(stt_req.Temperature, 'f', 6, 64))); err != nil {
			return nil, errors.New("failed to write temperature field")
		}
	}

	if stt_req.Prompt != "" {
		if fw, err = w.CreateFormField("prompt"); err != nil {
			return nil, errors.New("failed to create form field for prompt")
		}
		if _, err := fw.Write([]byte(stt_req.Prompt)); err != nil {
			return nil, errors.New("failed to write prompt field")
		}
	}

	if stt_req.Language != "" {
		if fw, err = w.CreateFormField("language"); err != nil {
			return nil, errors.New("failed to create form field for language")
		}
		if _, err := fw.Write([]byte(stt_req.Language)); err != nil {
			return nil, errors.New("failed to write language field")
		}
	}

	// add form for field if using word timestamps or segment timestamps
	if isWordStampReq || isSegmentStampReq {
		// verbose json
		if fw, err = w.CreateFormField("response_format"); err != nil {
			return nil, errors.New("failed to create form field for response_format")
		}
		if _, err := fw.Write([]byte(stt_req.ResponseFormat)); err != nil {
			return nil, errors.New("failed to write response_format field")
		}

		// timestamp granularities
		if fw, err = w.CreateFormField("timestamp_granularities[]"); err != nil {
			return nil, errors.New("failed to create form field for timestamp_granularities")
		}
		if _, err := fw.Write([]byte(stt_req.TimestampGranularities[0])); err != nil {
			return nil, errors.New("failed to write timestamp_granularities field")
		}
	}

	// close writer
	w.Close()

	// http request
	httpReq, err := http.NewRequest("POST", req_url, &b)
	if err != nil {
		return nil, errors.New("failed to create http request")
	}
	httpReq.Header.Set("Content-Type", w.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+APIKey)

	// send the req
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("failed to send request")
	}
	defer resp.Body.Close()

	// req response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body")
	}

	return respBody, nil
}

func (c *openaiAPI) OpenAISpeechToTextWordTimestamps(req_body *OATranscriptionDefaultReq) (*OATranscriptionWordTimestampResp, error) {
	var result OATranscriptionWordTimestampResp
	isWordStamp := true
	isTranscription := true

	respBody, err := sttBaseReq(isTranscription, isWordStamp, false, *req_body, c.apiKey)
	if err != nil {
		return nil, err
	}

	// parse response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New("failed to unmarshal response")
	}

	// check here if error happen on response
	if result.Error.Message != "" {
		return nil, errors.New(result.Error.Message)
	}

	return &result, nil
}

func (c *openaiAPI) OpenAISpeechToTextSegmentTimestamps(req_body *OATranscriptionDefaultReq) (*OATranscriptionSegmentResp, error) {
	var result OATranscriptionSegmentResp
	isSegmentReq := true
	isTranscription := true

	respBody, err := sttBaseReq(isTranscription, false, isSegmentReq, *req_body, c.apiKey)
	if err != nil {
		return nil, err
	}

	// parse response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New("failed to unmarshal response")
	}

	// check here if error happen on response
	if result.Error.Message != "" {
		return nil, errors.New(result.Error.Message)
	}

	return &result, nil
}

func (c *openaiAPI) OpenAISpeechToTextDefault(req_body *OATranscriptionDefaultReq) (*OATranscriptionDefaultResp, error) {
	var result OATranscriptionDefaultResp
	isTranscription := true

	respBody, err := sttBaseReq(isTranscription, false, false, *req_body, c.apiKey)
	if err != nil {
		return nil, err
	}

	// parse response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New("failed to unmarshal response")
	}

	// check here if error happen on response
	if result.Error.Message != "" {
		return nil, errors.New(result.Error.Message)
	}

	return &result, nil
}

func (c *openaiAPI) OpenAISpeechToTextTranslation(req_body *OATranslationDefaultReq) (*OATranscriptionDefaultResp, error) {
	var result OATranscriptionDefaultResp
	isTranscription := false

	req := OATranscriptionDefaultReq{
		File:        req_body.File,
		Filename:    req_body.Filename,
		Temperature: req_body.Temperature,
		Prompt:      req_body.Prompt,
	}

	respBody, err := sttBaseReq(isTranscription, false, false, req, c.apiKey)
	if err != nil {
		return nil, err
	}

	// parse response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New("failed to unmarshal response")
	}

	// check here if error happen on response
	if result.Error.Message != "" {
		return nil, errors.New(result.Error.Message)
	}

	return &result, nil
}
