package openai

import (
	"errors"
	"net/http"
	"time"
)

const (
	OAUrlBase                  = "https://api.openai.com/v1"
	OAUrlTextCompletions       = OAUrlBase + "/chat/completions"
	OAUrlImageGenerationsDallE = OAUrlBase + "/images/generations"
	OAUrlTextToSpeech          = OAUrlBase + "/audio/speech"
	OAUrlSTTTranscription      = OAUrlBase + "/audio/transcriptions"
	OAUrlSTTTranslation        = OAUrlBase + "/audio/translations"
	OABaseModel                = "gpt-4o-mini"
)

type OpenAI interface {

	// OpenAISendMessage sends a message to OpenAI's API and handles the request and response format.
	//
	// This function creates and sends a request to the OpenAI API, allowing for custom request bodies and response formats.
	// It either uses a provided custom request body or constructs a request body based on the provided message content.
	// If response formatting is required, the `OACreateResponseFormat()` function can be used to generate the response format schema.
	//
	// Parameters:
	//   - content: A pointer to a slice of OAMessageReq, which represents the request message content to be sent to OpenAI.
	//     This is used if `with_custom_reqbody` is set to false.
	//   - with_format_response: A boolean indicating whether a response format should be applied. If true, `format_response` must be provided.
	//   - format_response: A map containing the JSON schema for formatting the response (can be created using OACreateResponseFormat).
	//   - with_custom_reqbody: A boolean indicating whether a custom request body (`req_body_custom`) should be used.
	//   - req_body_custom: A pointer to an OAReqBodyMessageCompletion struct. This is used if `with_custom_reqbody` is true.
	//
	// Returns:
	//   - A pointer to an OAChatCompletionResp struct containing the API response.
	//   - An error if the request fails, or if invalid parameters are provided.
	//
	// Example usage:
	//
	//	content := []OAMessageReq{
	//	  {Role: "user", Content: "What is the weather like today?"},
	//	}
	//
	//	formatResponse := OACreateResponseFormat("WeatherResponse", map[string]interface{}{
	//	  "temperature": map[string]interface{}{"type": "string"},
	//	  "condition": map[string]interface{}{"type": "string"},
	//	})
	//
	//	response, err := openaiAPIInstance.OpenAISendMessage(&content, true, formatResponse, false, nil)
	//	if err != nil {
	//	    log.Fatalf("Failed to send message: %v", err)
	//	}
	//	fmt.Printf("API response: %+v\n", response)
	//
	// Notes:
	//   - The function checks for invalid states, such as missing content or custom request bodies when required.
	//   - The request is sent as a POST request with a JSON payload, and the response is decoded into the OAChatCompletionResp struct.
	//
	// References:
	// - Official OpenAI API documentation: https://platform.openai.com/docs/api-reference/chat/create
	OpenAISendMessage(content *[]OAMessageReq, with_format_response bool, format_response *map[string]interface{}, with_custom_reqbody bool, req_body_custom *OAReqBodyMessageCompletion) (*OAChatCompletionResp, error)

	// OpenAIGetFirstContentDataResp retrieves the first content data from an OpenAI API response.
	//
	// This function sends a message request to the OpenAI API using the given content,
	// and then extracts the first response content that basically the message response message from the API's response that can use for simplicity reason if you just need to use it like "one shot" request, so you can only the return content straight away and not the full response structure of OpenAI Response.
	//
	// Parameters:
	//   - content: A pointer to a slice of OAMessageReq, which represents the request message content to be sent to OpenAI.
	//   - with_format_response: A boolean indicating whether the response should be formatted.
	//   - format_response: A map that contains additional formatting options for the response. if you need to use the format_response that supported by OpenAI API. Official Docs and structure about structured response OpenAPI schema in: https://platform.openai.com/docs/guides/structured-outputs/examples
	//
	// Returns:
	//   - A pointer to an OAMessage struct that contains the first content data from the response.
	//   - An error if the request to OpenAI fails.
	//
	// Example usage:
	//
	//	content := []OAMessageReq{...}
	//	formatOptions := map[string]interface{}{
	//	  "option1": "value1",
	//	  // add formatting options here
	//	}
	//	firstContent, err := openaiAPIInstance.OpenAIGetFirstContentDataResp(&content, true, formatOptions)
	//	if err != nil {
	//	    log.Fatalf("Failed to get first content data: %v", err)
	//	}
	//	fmt.Println("First response content:", firstContent)
	//
	// References:
	// - Official OpenAI API documentation: https://platform.openai.com/docs/api-reference/chat/create
	OpenAIGetFirstContentDataResp(content *[]OAMessageReq, with_format_response bool, format_response *map[string]interface{}, with_custom_reqbody bool, req_body_custom *OAReqBodyMessageCompletion) (*OAMessage, error)

	// OpenAICreateImageDallE generates images based on a text prompt using either the DALL-E 2 or DALL-E 3 model.
	//
	// This method constructs an HTTP request to OpenAI's image generation API, validates input requirements for each model,
	// and includes optional parameters for fine-tuning image quality, response format, style, size, and number of generated images.
	//
	// Parameters:
	//
	//   - req_body (*OAReqImageGeneratorDallE): A pointer to a struct containing image generation request parameters.
	//
	//     Fields in OAReqImageGeneratorDallE struct:
	//
	//   - Prompt (string): A descriptive prompt that the DALL-E model will use to generate images. This field is required.
	//
	//   - Model (string): Specifies the DALL-E model version, either "dall-e-2" or "dall-e-3". This is required.
	//
	//   - N (*int): Optional. The number of images to generate, which can range between 1 and 10. Defaults to 1 if omitted.
	//
	//   - Quality (*string): Optional. Available only for the DALL-E 3 model, it can be set to "standard" (default) or "hd" for high definition.
	//
	//   - ResponseFormat (*string): Optional. Specifies the format of the generated image response, either "url" (default) or "b64_json" for a base64-encoded image.
	//
	//   - Size (*string): Optional. The image size, dependent on the model:
	//
	//   - DALL-E 2 supports "256x256", "512x512", and "1024x1024".
	//
	//   - DALL-E 3 supports "1024x1024", "1792x1024", and "1024x1792".
	//
	//   - Style (*string): Optional. Available only for the DALL-E 3 model. Choices are "vivid" (default) or "natural" for a more realistic style.
	//
	//   - User (*string): Optional. A unique identifier for the end user to monitor and detect abuse, helping OpenAI with usage tracking.
	//
	// Returns:
	//   - (*OAImageGeneratorDallEResp, error): On success, returns a pointer to an `OAImageGeneratorDallEResp` struct containing the
	//     generated image details. Returns an error if any validation fails or if the request is unsuccessful.
	//
	// Example usage:
	//
	//	reqBody := &OAReqImageGeneratorDallE{
	//	    Prompt: "A futuristic cityscape at sunset",
	//	    Model: "dall-e-3",
	//	    N: ptr(3),
	//	    Quality: ptr("hd"),
	//	    Style: ptr("vivid"),
	//	    Size: ptr("1024x1024"),
	//	    ResponseFormat: ptr("url"),
	//	}
	//
	//	imageResp, err := apiClient.OpenAICreateImageDallE(reqBody)
	//	if err != nil {
	//	    log.Fatalf("Image generation failed: %v", err)
	//	}
	//
	// Function Logic:
	//  1. **Model Validation**: Ensures `Model` is either "dall-e-2" or "dall-e-3". If not, returns an error.
	//  2. **N Validation**: If `N` is provided, checks if it falls between 1 and 10 (inclusive). If out of range, returns an error.
	//  3. **Quality Validation**:
	//     - If `Model` is "dall-e-2", `Quality` should be nil. Returns an error if a quality value is provided for this model.
	//     - If `Model` is "dall-e-3", `Quality` can be "standard" or "hd". If any other value is provided, returns an error.
	//  4. **Style Validation**:
	//     - Ensures `Style` is only available for "dall-e-3". For "dall-e-2", returns an error if `Style` is provided.
	//     - Valid values for `Style` are "vivid" or "natural"; any other value results in an error.
	//  5. **ResponseFormat Validation**:
	//     - If `ResponseFormat` is specified, it must be either "url" or "b64_json". Returns an error for other values.
	//  6. **API Key Check**: Confirms that the API key is set; returns an error if it's empty.
	//  7. **JSON Marshalling**: Serializes `req_body` to JSON format for the request body.
	//  8. **Request Creation and Headers**: Sets up an HTTP POST request with necessary headers (`Content-Type` and `Authorization`).
	//  9. **Response Handling**:
	//     - If the HTTP response status is not 200 OK, reads and closes the response body, then returns an error.
	//     - On successful response, decodes JSON data into `OAImageGeneratorDallEResp` struct and returns it.
	//
	// Considerations:
	//   - The function relies on an HTTP client specified in the API client’s configuration (c.config.httpClient).
	//   - In the event of a non-200 HTTP status, the function reads the body to allow graceful closing of the connection.
	//   - OpenAI API may apply rate limiting, so ensure retry or error handling mechanisms are in place for high-frequency requests.
	//
	// References:
	//   - OpenAI DALL E Image Generation API: https://platform.openai.com/docs/api-reference/images/create
	OpenAICreateImageDallE(req_body *OAReqImageGeneratorDallE) (*OAImageGeneratorDallEResp, error)

	// OpenAITextToSpeech converts a text input into a speech audio file using OpenAI's TTS models.
	// This function validates the input parameters, prepares the request, sends it to the OpenAI API,
	// and returns the audio response encoded in base64 format.
	//
	// Parameters:
	//   - req_body (*OAReqTextToSpeech): A pointer to the OAReqTextToSpeech struct containing the TTS parameters.
	//
	// Returns:
	//   - (*OATextToSpeechResp, error): On success, returns a pointer to an OATextToSpeechResp struct containing:
	//   - FormatAudio: The file extension for the audio format (e.g., ".mp3").
	//   - B64JSON: A base64-encoded string representing the audio file content.
	//     On failure, returns an error.
	//
	// Errors:
	//   - Returns an error if required fields are missing or invalid, including:
	//   - Invalid Model (must be "tts-1" or "tts-1-hd").
	//   - Missing Input text.
	//   - Invalid Voice option (allowed values: "alloy", "echo", "fable", "onyx", "nova", "shimmer").
	//   - Invalid ResponseFormat (allowed values: "mp3", "opus", "aac", "flac", "wav", "pcm").
	//   - Speed out of range (0.25 to 4.0).
	//   - Also returns an error if the API key is missing, or if any part of the HTTP request/response fails.
	//
	// Example Usage:
	//
	//	reqBody := OAReqTextToSpeech{
	//	    Model:          "tts-1",
	//	    Input:          "Hello, world!",
	//	    Voice:          "alloy",
	//	    ResponseFormat: "mp3",
	//	}
	//
	//	resp, err := openAI.OpenAITextToSpeech(&reqBody)
	//	if err != nil {
	//	    log.Fatalf("Text-to-Speech conversion failed: %v", err)
	//	}
	//
	//	fmt.Println("Audio Format:", resp.FormatAudio)
	//	fmt.Println("Base64 Encoded Audio:", resp.B64JSON)
	//
	// References:
	//   - TTS OpenAI: https://platform.openai.com/docs/api-reference/audio/createSpeech
	OpenAITextToSpeech(req_body *OAReqTextToSpeech) (*OATextToSpeechResp, error)

	OpenAISpeechToTextDefault(req_body *OATranscriptionDefaultReq) (*OATranscriptionDefaultResp, error)

	OpenAISpeechToTextWordTimestamps(req_body *OATranscriptionDefaultReq) (*OATranscriptionWordTimestampResp, error)

	OpenAISpeechToTextSegmentTimestamps(req_body *OATranscriptionDefaultReq) (*OATranscriptionSegmentResp, error)

	OpenAISpeechToTextTranslation(req_body *OATranslationDefaultReq) (*OATranscriptionDefaultResp, error)
}

// Config holds the configuration for OpenAI API client
type Config struct {
	httpClient    *http.Client
	openAIBaseUrl string
	openAIModel   string
}

// default configuration for OpenAI API client
func DefaultConfig() *Config {
	return &Config{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		// user base url for chat completions endpoint with using gpt-4o-mini model
		openAIBaseUrl: OAUrlTextCompletions,
		openAIModel:   OABaseModel,
	}
}

// client implementation for OpenAI API interfaces
type openaiAPI struct {
	apiKey             string
	openaiOrganization string
	openaiProject      string
	config             *Config
}

// client options for configuring the OpenAI API client
type ClientOption func(*Config)

// New creates and returns a new instance of the OpenAI API client.
//
// This function initializes the OpenAI client with the provided API key, organization, and project ID.
// Additionally, it applies any optional configuration options passed via the variadic `opts` parameter.
//
// Parameters:
//   - apiKey: The API key used to authenticate with OpenAI's API. This is required and cannot be empty.
//   - openaiOrganization: (Optional) The organization ID associated with your OpenAI account. This can be omitted if not applicable.
//   - openaiProject: (Optional) The project ID associated with your OpenAI account. This can be omitted if not applicable.
//   - opts: A variadic parameter that accepts one or more `ClientOption` functions to customize the client's behavior.
//
// Returns:
//   - An OpenAI client instance, or an error if the API key is missing or another issue occurs during initialization.
//
// Default Configuration Values:
//   - httpClient: The HTTP client used for making requests. By default, the client has a
//     timeout of 60 seconds (`http.Client{ Timeout: 60 * time.Second }`).
//   - openAIBaseUrl: The default base URL for the OpenAI is set to the `/chat/completions` endpoint
//     The default value is `"https://api.openai.com/v1/chat/completions"`.
//   - openAIModel: The default model for message processing is `"gpt-4o-mini"`, which specifies
//     the Claude model version that will be used to generate responses.
//
// Example usage:
//
//	// Initialize the OpenAI client with an API key
//	openaiClient, err := New("your-api-key", "your-org-id", "your-project-id")
//	if err != nil {
//	    log.Fatalf("Failed to create OpenAI client: %v", err)
//	}
//
//	// Optionally, provide custom options such as timeouts, base URL, or HTTP clients
//	customClient, err := New("your-api-key", "your-org-id", "your-project-id", WithCustomTimeout(30 * time.Second))
//	if err != nil {
//	    log.Fatalf("Failed to create OpenAI client with custom options: %v", err)
//	}
//
// Notes:
//   - The `apiKey` is required and must be provided, otherwise an error will be returned.
//   - The `openaiOrganization` and `openaiProject` parameters are optional and can be left empty if not needed.
//   - `ClientOption` is a functional option pattern that allows customization of the client, such as setting custom HTTP clients or changing API base URLs.
//
// References:
//   - Official OpenAI API authentication: https://platform.openai.com/docs/api-reference/authentication
func New(apiKey string, openaiOrganization string, openaiProject string, opts ...ClientOption) (OpenAI, error) {
	// from openai docs on
	// https://platform.openai.com/docs/api-reference/authentication
	// organization and project id is optional
	if apiKey == "" {
		return nil, errors.New("API Key is empty")
	}

	// create new OpenAI instance from private struct
	config := DefaultConfig()

	// apply custom options
	for _, opt := range opts {
		opt(config)
	}

	return &openaiAPI{
		apiKey:             apiKey,
		openaiOrganization: openaiOrganization,
		openaiProject:      openaiProject,
		config:             config,
	}, nil
}

// use if need custom http client setup, use it on New function initiate
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Config) {
		c.httpClient = httpClient
	}
}

// custom base url setup if need using different endpoint maybe like dalle or whisper or other, use it on New function initiate
func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *Config) {
		c.openAIBaseUrl = baseUrl
	}
}

// custom model setup if need using different model maybe like gpt-4o or gpt-4o-turbo or other, use it on New function initiate
func WithModel(model string) ClientOption {
	return func(c *Config) {
		c.openAIModel = model
	}
}
