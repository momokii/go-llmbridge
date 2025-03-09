package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// OACreateResponseFormat creates a response format using a JSON Schema for OpenAI response format data requests.
//
// This function is used to generate a JSON Schema structure that can be passed as a parameter
// to the OpenAISendMessage() function, providing a standard format for the response.
//
// Parameters:
//   - jsonName: A string representing the name of the JSON schema.
//   - jsonSchema: A map of string to interface, representing the schema data, specifically the properties
//     of the schema as defined by the OpenAI structured output documentation.
//
// Returns:
//   - A map[string]interface{} representing the formatted response structure using JSON Schema,
//     including the schema name and its associated properties.
//
// Example usage:
//
//		jsonSchema := map[string]interface{}{
//	 "type": "object",
//		"properties": map[string]interface{}{
//			  "title": map[string]interface{}{
//			    "type": "string",
//			  },
//			  "description": map[string]interface{}{
//			    "type": "string",
//			  },
//		   },
//		}
//
//		formattedResponse := OACreateResponseFormat("MySchema", jsonSchema)
//		fmt.Printf("Formatted response: %v\n", formattedResponse)
//
// JSON Schema Structure:
//   - The structure returned by this function will conform to the schema guidelines provided by OpenAI.
//     More details and examples can be found at the following link:
//     https://platform.openai.com/docs/guides/structured-outputs/examples
//
// Returned Structure Example:
//
//	{
//	  "type": "json_schema",
//	  "json_schema": {
//	    "name": "MySchema",
//	    "schema": {
//	      "type": "object",
//	      "properties": {
//	        "title": {
//	          "type": "string"
//	        },
//	        "description": {
//	          "type": "string"
//	        }
//	      }
//	    }
//	  }
//	}
func OACreateResponseFormat(jsonName string, jsonSchema map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": "json_schema",
		"json_schema": map[string]interface{}{
			"name":   jsonName,
			"schema": jsonSchema,
		},
	}
}

// OACreateOneContentVision constructs a vision content payload for uploading an image (either as a URL or base64-encoded string)
// along with optional text to the OpenAI API.
//
// This function creates a list of `OAContentVisionBaseReq` structures, enabling you to upload both image data (via URL or base64 encoding)
// and optional descriptive text to OpenAI's vision endpoint. Supported media types include JPEG, PNG, JPG, GIF, and WebP.
//
// Parameters:
//   - media_type (string): The MIME type of the image when using base64 encoding. This is required when `using_image_url` is false.
//     Supported types include:
//   - "image/png"
//   - "image/jpeg"
//   - "image/jpg"
//   - "image/gif"
//   - "image/webp"
//   - using_image_url (bool): Specifies whether the image is provided as a URL or a base64-encoded string.
//   - If `true`, the function expects `url_or_base64encoding` to be a valid URL.
//   - If `false`, the function expects `url_or_base64encoding` to be a base64-encoded image string, and `media_type` must be provided.
//   - url_or_base64encoding (string): The image data provided as either a URL (when `using_image_url` is `true`) or a base64-encoded
//     string (when `using_image_url` is `false`). If this value is empty, the function returns an error indicating that both
//     `media_type` and `url_or_base64encoding` must be provided.
//   - text_content (string): An optional text string to accompany the image content. This will be included as a separate text
//     component if provided.
//
// Returns:
//
//	([]OAContentVisionBaseReq, error): A slice of `OAContentVisionBaseReq` structs containing the image (and optional text).
//	If the required parameters are not met or the media type is unsupported, an error is returned.
//
// Example usage:
//
//	// Example URL-based request
//	visionContent, err := OACreateOneContentVision("", true, "https://example.com/sample-image.jpg", "This is an example image.")
//	if err != nil {
//	    log.Fatalf("Error generating vision content: %v", err)
//	}
//
//	// Example base64-encoded request
//	base64Image := "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAYAAABccqhmAAAACXBIWXMAAB7CAAAewgFu0HU+AAAK..." // truncated for example
//	visionContent, err := OACreateOneContentVision("image/png", false, base64Image, "This is an example image.")
//	if err != nil {
//	    log.Fatalf("Error generating vision content: %v", err)
//	}
//
// Function Logic:
//  1. **Input Validation**: Checks if `url_or_base64encoding` is empty. If it is, returns an error requiring both
//     `media_type` and `url_or_base64encoding` to be provided. If using base64 encoding (`using_image_url` is false), `media_type`
//     must also be specified.
//  2. **Supported Media Types**: Validates that `media_type` is one of the supported image types if `using_image_url` is false.
//     Supported types include "image/png", "image/jpeg", "image/jpg", "image/gif", and "image/webp". If an unsupported type
//     is provided, the function returns an error listing the valid types.
//  3. **Data Preparation**: If `using_image_url` is true, `imageData` is set to `url_or_base64encoding`. Otherwise, `imageData`
//     is created as a "data URI" by prepending `media_type` and "base64," to the encoded string. This format complies with the
//     OpenAI API's expectations for base64-encoded images.
//  4. **Image Content**: A `OAContentVisionBaseReq` struct is created for the image content, setting `Type` to `"image_url"`,
//     and the image data is assigned to the `ImageUrl` field.
//  5. **Optional Text Content**: If `text_content` is provided, another `OAContentVisionBaseReq` struct is appended with `Type` set
//     to `"text"` and `Text` containing the provided text. This allows both image and text content to be sent in a single request.
//  6. **Return**: Returns the slice of `OAContentVisionBaseReq` structs, which includes the image content (and text content,
//     if provided), ready for an OpenAI vision request.
//
// Notes:
//   - Ensure that the `url_or_base64encoding` contains a valid URL when `using_image_url` is true or a base64-encoded image string
//     when false.
//   - This function supports only a single image and an optional text. Multiple images or additional text content would
//     require separate calls or modifications to the function.
//   - OpenAI’s API currently supports base64 and URL images as part of its vision feature, making it possible to use both methods
//     with this function.
//
// Considerations:
//   - Base64-encoded images should be appropriately sized or compressed before encoding to avoid excessively large requests.
//   - URLs provided should be publicly accessible or authenticated as needed by the OpenAI API.
//   - this function hope can make you easier for send vision content if just contain one image and optional text content, if you need more than one image, you can create your own structure based on OpenAI Docs with struct OAContentVisionBaseReq & OAContentVisionImageUrl (for content structure) and append it to the slice of OAContentVisionBaseReq
//
// Reference for Vision OpenAI Docs:
// - Official OpenAI API documentation: https://platform.openai.com/docs/guides/vision
func OACreateOneContentVision(media_type string, using_image_url bool, url_or_base64encoding string, text_content string) ([]OAContentVisionBaseReq, error) {
	if url_or_base64encoding == "" {
		return nil, errors.New("media_type and url_or_base64encoding must be provided")
	}

	if media_type == "" && !using_image_url {
		return nil, errors.New("media_type must be provided when using base64 encoding")
	}

	if !using_image_url && media_type != "image/png" && media_type != "image/jpeg" && media_type != "image/jpg" && media_type != "image/gif" && media_type != "image/webp" {
		return nil, errors.New("media_type must be image/png, image/jpeg, or image/jpg")
	}

	var imageData string

	// data url or base64 encoding and the format is based on OpenAI API Docs
	if using_image_url {
		imageData = url_or_base64encoding
	} else {
		imageData = "data:" + media_type + ";base64," + url_or_base64encoding
	}

	contentVision := []OAContentVisionBaseReq{
		{
			Type: "image_url",
			ImageUrl: &OAContentVisionImageUrl{
				Url: imageData,
			},
		},
	}

	if text_content != "" {
		contentVision = append(contentVision, OAContentVisionBaseReq{
			Type: "text",
			Text: &text_content,
		})
	}

	return contentVision, nil
}

func (c *openaiAPI) OpenAISendMessage(content *[]OAMessageReq, with_format_response bool, format_response *map[string]interface{}, with_custom_reqbody bool, req_body_custom *OAReqBodyMessageCompletion) (*OAChatCompletionResp, error) {

	// var reqBody interface{}
	var reqBody interface{}

	if c.apiKey == "" {
		return nil, errors.New("API Key is empty")
	}

	// check if with_format_response is true, format_response must be provided
	if with_format_response && format_response == nil {
		return nil, errors.New("format_response must be provided when with_format_response is true")
	}

	// check if with_custom_reqbody is true, req_body_custom must be provided
	if with_custom_reqbody && req_body_custom.Messages == nil {
		return nil, errors.New("req_body_custom must be provided when with_custom_reqbody is true")
	}

	// check if with_custom_reqbody is false, content must be provided
	if !with_custom_reqbody && content == nil {
		return nil, errors.New("content must be provided")
	}

	// create request body
	if with_custom_reqbody {

		if with_format_response {
			req_body_custom.ResponseFormat = *format_response
		}

		reqBody = req_body_custom

	} else {
		reqData := OAReqBodyMessageCompletion{
			Model:    c.config.openAIModel,
			Messages: content,
		}

		// if using format response add response format to request body
		if with_format_response {
			reqData.ResponseFormat = *format_response
		}

		reqBody = reqData
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("Failed to marshal request body")
	}

	// send req to openai
	req, err := http.NewRequest(http.MethodPost, c.config.openAIBaseUrl, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, errors.New("Failed to create request")
	}

	// header setup
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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

	// decode response
	var result OAChatCompletionResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.New("Failed to decode response: " + err.Error())
	}

	return &result, nil // return response
}

func (c *openaiAPI) OpenAIGetFirstContentDataResp(content *[]OAMessageReq, with_format_response bool, format_response *map[string]interface{}, with_custom_reqbody bool, req_body_custom *OAReqBodyMessageCompletion) (*OAMessage, error) {
	// send request to openai
	resp, err := c.OpenAISendMessage(content, with_format_response, format_response, with_custom_reqbody, req_body_custom)
	if err != nil {
		return nil, err
	}

	// get content first data
	data := resp.Choices[0].Message

	return &data, nil
}
