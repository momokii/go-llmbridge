# Go LLM Wrapper

The Go LLM Wrapper provides a streamlined, pure Go interface for interacting with leading LLM APIs, including Claude (Anthropic) and OpenAI. This lightweight wrapper supports both text and vision capabilities, simplifying the integration of advanced language model functionality into Go applications.

## Changelog
### New Update Features
- 🆕 Added OpenAI Speech-To-Text (STT) support

## Features

- 🚀 Pure Go implementation with zero external dependencies
- 💬 Multiple Provider Support
- 👁️ Comprehensive support for both text-based and vision capabilities
- 🎙️ Text-to-Speech generation support
- 🎙️ Speech-to-Text support 🆕
- 🖼️ Image generation with DALL-E
- ⚡ Minimal and efficient
- 🛠️ Flexible configuration options for each provider

## Supported Providers (Documentation)

- [Claude (Anthropic)](https://docs.anthropic.com/)
- [OpenAI](https://platform.openai.com/docs)

## Requirements

- Go 1.21.0

## Installation

```bash
go get github.com/momokii/go-llmbridge
```

## Quick Start

### Claude & OpenAI Client Initialization

```go
package main

import (
    "github.com/momokii/go-llmbridge/claude"
    "github.com/momokii/go-llmbridge/openai"
)

func main() {
    // if you want to use custom http client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

    // init claude with some custom options (still the client also have default config)
    claudeClient, err := claude.New(
		os.Getenv("CLAUDE_API_KEY"), // rqeuired
        // additional config
		claude.WithHTTPClient(httpClient),
		claude.WithBaseUrl(os.Getenv("CLAUDE_BASE_URL")),
		claude.WithModel(os.Getenv("CLAUDE_MODEL")),
		claude.WithAnthropicVersion(os.Getenv("CLAUDE_ANTHROPIC_VERSION")),
	)
	if err != nil {
		panic(err)
	}

    // init openai with some custom options (still the client also have default config)
	gptClient, err = openai.New(
		os.Getenv("OA_APIKEY"), // required
        // additional config
		os.Getenv("OA_ORGANIZATIONID"),
		os.Getenv("OA_PROJECTID"),
		openai.WithHTTPClient(httpClient),
		openai.WithModel("gpt-4o-mini"),
		openai.WithBaseUrl("https://api.openai.com/v1/chat/completions"),
	)
	if err != nil {
		panic(err)
	}
}
```

## Usage Examples

### Text Completion

Documentation for text completion for supported LLM Provider
- [Claude Text](https://docs.anthropic.com/en/api/messages)
- [OpenAI Text](https://platform.openai.com/docs/api-reference/chat)

#### Claude
##### Simple Completion Chat
```go
// create body message request
claudeMessageBodyText := []claude.ClaudeMessageReq{
    {
        Role:    "user",
        Content: "Hello, how are you?",
    },
}

// send request to Claude to Get The Content response from model
if content, err := claudeClient.ClaudeGetFirstContentDataResp(&claudeMessageBodyText, 256, false, nil); err != nil {
    fmt.Println("error claude req: " + err.Error())
} else {
    fmt.Println("claude response content: ")
    fmt.Println(content)
}

```
##### Chaining Completion Chat
```go
// for chaining message, basically you just add "all" conversation so the model can known the context
claudeMessageBodyChain := []claude.ClaudeMessageReq{
    {
        Role:    "user",
        Content: "Hello, how are you?",
    },
    {
        Role:    "assistant",
        Content: "I'm fine, thank you. How can I help you today?",
    },
    {
        Role:    "user",
        Content: "Give me a joke",
    },
    {
        Role: "assistant",
        Content: `Sure, here's a joke for you:

        Why don't scientists trust atoms?

        Because they make up everything!

        I hope that gives you a little chuckle. Do you have any favorite types of jokes or would you like to hear another one?`,
    },
    {
        Role:    "user",
        Content: "yes.",
    },
}

// Claude has some "optional" parameters that you can find in the Claude Docs. 
// If you need to use these, you can create a custom request body like below.
claudeRequestBody := claude.ClaudeReqBody{
    Model:       os.Getenv("CLAUDE_MODEL"),
    MaxTokens:   2560,
    Messages:    claudeMessageBodyChain,
    Temperature: 1.0,
    System:      "You are AI assistant with great knowledge at comedy like a stand-up comedian. You can give a joke, comedy, or funny story to user. The joke can be in any form, such as a pun, a one-liner, or a short story.",
}

// Try using a custom request body with chained messages
// Retrieve the full response content from Claude’s response
// When using a "custom" request, a different approach is required
if content, err := claudeClient.ClaudeSendMessage(nil, 256, true, &claudeRequestBody); err != nil {
    fmt.Println("error claude req: " + err.Error())
} else {
    fmt.Println("\nclaude response custom body: ")
    fmt.Println(content)
}

```


#### OpenAI GPT
##### Simple Completion Chat
```go
// message req structure for openai model
gptMessageBodyText := []openai.OAMessageReq{
    {
        Role:    "user",
        Content: "Hello, how are you?",
    },
}

// send request to openai model and get the full content response from model
if content, err := gptClient.OpenAISendMessage(&gptMessageBodyText, false, nil, false, nil); err != nil {
    fmt.Println("error gpt req: " + err.Error())
} else {
    fmt.Println("gpt response content: ")
    fmt.Println(content)
}
```
##### Chaining Completion Chat & Using Format Response
```go
// for chaining message, basically you just add "all" conversation so the model can known the context
gptMessageChain := []openai.OAMessageReq{
    {
        Role:    "user",
        Content: "Hello, how are you?",
    },
    {
        Role:    "assistant",
        Content: "I'm fine, thank you. How can I help you today?",
    },
    {
        Role:    "user",
        Content: "Give me some joke from a stand-up comedian.",
    },
}

// OpenAI supports a parameter that allows you to create a structured response format, making it easier to work with the response structure (this feature is not yet supported in Claude). 
// For reference, see OpenAI Docs: https://platform.openai.com/docs/guides/structured-outputs

// I have provided a function here to help you easily define a response format. 
// Example usage is shown below:
gptResponseFormatCustom := openai.OACreateResponseFormat(
    "testing_response_format",
    map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "data": map[string]interface{}{
                "type": "array",
                "items": map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "name": map[string]interface{}{
                            "type": "string",
                        },
                        "joke": map[string]interface{}{
                            "type": "string",
                        },
                    },
                },
            },
        },
    },
)

// In the GPT version, you can also create a custom request body, using "optional" parameters based on the OpenAI Docs
gptMessageCustomWithResponseFormat := openai.OAReqBodyMessageCompletion{
    Messages:       gptMessageChain,
    Model:          "gpt-4o-mini",
    ResponseFormat: gptResponseFormatCustom,
}

// send request to openai model and get the full content response from mpdel
if content, err := gptClient.OpenAISendMessage(nil, false, nil, true, &gptMessageCustom); err != nil {
    fmt.Println("error gpt req: " + err.Error())
} else {
    fmt.Println("\ngpt response content only but custom req body: ")
    fmt.Println(content)
}

// OR if you are not using a custom request body, you can apply the response format directly as shown below
if content, err := gptClient.OpenAISendMessage(&gptMessageChain, true, gptResponseFormatCustom, false, nil); err != nil {
    fmt.Println("error gpt req: " + err.Error())
} else {
    fmt.Println("\ngpt response content only but custom req body: ")
    fmt.Println(content)
}
```

### Vision Analysis

Currently, support for vision analysis varies between Claude and OpenAI GPT, as detailed in their respective documentation:

- **Claude** supports image data in base64 encoding only, as specified in the [Claude Vision Docs](https://docs.anthropic.com/en/docs/build-with-claude/vision).
- **OpenAI GPT** supports image data in both URL and base64 formats, as outlined in the [OpenAI Vision Docs](https://platform.openai.com/docs/guides/vision).

Both models, Claude and OpenAI GPT, currently support only image-based data.

```go
// image data

// encode file data
file, err := os.Open("image.png")
if err != nil {
    fmt.Println("error open file: " + err.Error())
}
defer file.Close()

fileContent, err := io.ReadAll(file)
if err != nil {
    fmt.Println("error read file: " + err.Error())
}

encodeFile := base64.StdEncoding.EncodeToString(fileContent) // encode file data
imageUrl := "your-url.jpg"

```

#### Claude
##### Single Image
```go
// Vision supports handling multiple image data, but if you only need to upload a single image, 
// there is a simplified function to structure content data for a single image, as shown below:
contentVision, err := claude.ClaudeCreateOneContentImageVisionBase64("image/png", encodeFile, "make joke with this image")
if err != nil {
    fmt.Println("error create content vision: " + err.Error())
}

// create message structure
claudeMessageBodyFile := []claude.ClaudeMessageReq{
    {
        Role:    "user",
        Content: contentVision,
    },
}

// Send request
// The function below is a simplified version of ClaudeSendMessage(), providing just the content response from OpenAI
// rather than the full response data. Use this function if you only need the model's answer (content).
if content, err := claudeClient.ClaudeGetFirstContentDataResp(&claudeMessageBodyFile, 2048, false, nil); err != nil {
    fmt.Println("error claude req: " + err.Error())
} else {
    fmt.Println("\nclaude response vision: ")
    fmt.Println(content)
}

```
##### Multiple Image
```go
// Currently, there is no built-in function for handling multiple image contents in this repo, 
// but you can create the content data yourself using the available structs.
// To structure image vision content, combine ClaudeVisionContentBase{} and ClaudeVisionSource{}.
// Note: Vision on Claude currently only supports data payload in base64 encoding.
textContent := "what is the difference between these two images? or the two image is the same?"

// Create a message structure with multiple image content in the request data
claudeMessageBodyMultipleFile := []claude.ClaudeMessageReq{
    {
        Role: "user",
        Content: []claude.ClaudeVisionContentBase{
            {
                Type: "text",
                Text: &textContent,
            },
            {
                Type: "image",
                Source: &claude.ClaudeVisionSource{
                    Type:      "base64",
                    MediaType: "image/png",
                    Data:      encodeFile,
                },
            },
            {
                Type: "image",
                Source: &claude.ClaudeVisionSource{
                    Type:      "base64",
                    MediaType: "image/png",
                    Data:      encodeFile,
                },
            },
        },
    },
}

// send request
if content, err := claudeClient.ClaudeGetFirstContentDataResp(&claudeMessageBodyMultipleFile, 2048, false, nil); err != nil {
    fmt.Println("error claude req: " + err.Error())
} else {
    fmt.Println("\nclaude response vision: ")
    fmt.Println(content)
}
```

#### OpenAI GPT
##### Single Image
```go
// Vision supports multiple image data. If you only need to upload one image, 
// there is a function to simplify content structure creation for a single image, as shown below.
gptMessageVisionContentUrl, err := openai.OACreateOneContentVision("image/png", false, encodeFile, "make joke with this image")
if err != nil {
    panic(err)
}

// create message structure
gptMessageVisionContent := []openai.OAMessageReq{
    {
        Role:    "user",
        Content: gptMessageVisionContentUrl,
    },
}

// Send request
// The function below is a simplified version of OpenAISendMessage() that returns only the model’s response content, not the full data from OpenAI.
// Use this if you only need the response content.
if content, err := gptClient.OpenAIGetFirstContentDataResp(&gptMessageVisionContent, false, nil, false, nil); err != nil {
    fmt.Println("error gpt req: " + err.Error())
} else {
    fmt.Println("gpt response content: ")
    fmt.Println(content)
}
```
##### Multiple Image Send
```go
// Currently, there is no built-in function for handling multiple image content, but you can create the data using the available structs in this repo.
// For image vision content, you can combine OAContentVisionBaseReq{} and OAContentVisionImageUrl{}.
// Note: Vision on GPT supports both base64 and URL formats, so you can mix them if needed.
gptMessageVisionMultipleImageText := "Tell me the difference between these two images and also make joke with combined image"

// create message structure and also within on create message req data
gptMessageVisionMultipleImage := []openai.OAMessageReq{
    {
        Role: "user",
        Content: []openai.OAContentVisionBaseReq{
            {
                Type: "text",
                Text: &gptMessageVisionMultipleImageText,
            },
            {
                Type: "image_url",
                ImageUrl: &openai.OAContentVisionImageUrl{
                    Url: "data:image/png;base64," + encodeFile,
                },
            },
            {
                Type: "image_url",
                ImageUrl: &openai.OAContentVisionImageUrl{
                    Url: imageUrl,
                },
            },
        },
    },
}

// send request
if content, err := gptClient.OpenAIGetFirstContentDataResp(&gptMessageVisionMultipleImage, false, nil, false, nil); err != nil {
    fmt.Println("error gpt req: " + err.Error())
} else {
    fmt.Println("gpt response content: ")
    fmt.Println(content)
}

```

### Image Generation (OpenAI DALL-E)
```go
size := "1792x1024" // choosing size for image

// there are 2 response format for DALL-E image generator response format b64_json and url
response_b64 := "b64_json"
// response_url := "url"

// DALL-E image generator message body
dalleMessage := openai.OAReqImageGeneratorDallE{
    Model:          "dall-e-3",
    Prompt:         "A painting of a flower vase in the style of Picasso",
    Size:           &size,
    ResponseFormat: &response_b64,
}

if dalleRes, err := gptClient.OpenAICreateImageDallE(&dalleMessage); err != nil {
    fmt.Println("error dalle req: " + err.Error())
} else {
    fmt.Println("\ndalle response: ")
    fmt.Println(dalleRes) // response data
}

```

### Text To Speech (OpenAI TTS)
```go
// output data format for TTS is just base64 encode audio data
text := "Hello, my name is Momokii. I am a virtual assistant. I can help you with anything you need. How can I help you today?"
ttsReBody := openai.OAReqTextToSpeech{
    Model:          "tts-1",
    Voice:          "alloy",
    ResponseFormat: "mp3",
    Input:          text,
}

if ttsData, err := gptClient.OpenAITextToSpeech(&ttsReBody); err != nil {
    fmt.Println("error tts req: " + err.Error())
} else {
    fmt.Println("\ntts response: ")
    fmt.Println(ttsData)
}

```

### Speech To Text (OpenAI Whisper)

The library supports OpenAI's Whisper model for speech-to-text conversion with multiple output formats and capabilities. These functions accept audio in various formats (mp3, mp4, mpeg, mpga, m4a, wav, webm, flac, ogg) and offer different levels of detail in the transcription output.

Documentation for speech-to-text for OpenAI:
- [OpenAI Whisper API (Transcription)](https://platform.openai.com/docs/api-reference/audio/createTranscription)
- [OpenAI Whisper API (Translation)](https://platform.openai.com/docs/api-reference/audio/createTranslation)

#### Transcription Example Usage
Convert audio to text

```go
// NEW: Text To Speech API - Converts text to audio

// Data load example

// filepath example
req := &OATranslationDefaultReq{
    File:     "/path/to/audio.mp3",
}

// file open/ io.reader as input
file, err := os.Open("/path/to/audio.mp3")
if err != nil {
    log.Fatalf("Failed to open file: %v", err)
}
defer file.Close()

req := &OATranscriptionDefaultReq{
    File:     file,
    Filename: "audio.mp3", // Required when using io.Reader
}

```
For transcription requests, there will be three different functions, and the main difference lies in the response structure that the user will receive.

##### Default Request
```go
// NEW: Text To Speech API - Converts text to audio

// Transcription - Default Request
resp, err := gptClient.OpenAISpeechToTextDefault(&req)
if err != nil {
    log.Println("error stt req: " + err.Error())
} else {
    log.Println("\nstt response: ")
    log.Println(resp)
}

// example transcript response for default transcription
// structure of json/ struct will be similar like openai docs structure
// {
//   "text": "Imagine the wildest idea that you've ever had, and you're curious about how it might scale to something that's a 100, a 1,000 times bigger. This is a place where you can get to do that."
// }


```

##### Word Timestamps Request
```go
// NEW: Text To Speech API - Converts text to audio

// Transcription - Word timestamps Request
resp, err := gptClient.OpenAISpeechToTextWordTimestamps(&req)
if err != nil {
    log.Println("error stt req: " + err.Error())
} else {
    log.Println("\nstt WORD TIMESTAMPS response: ")
    log.Println(resp)
}

// example transcript response for word timestamps transcription
// structure of json/ struct will be similar like openai docs structure
// {
//   "task": "transcribe",
//   "language": "english",
//   "duration": 8.470000267028809,
//   "text": "The beach was a popular spot on a hot summer day. People were swimming in the ocean, building sandcastles, and playing beach volleyball.",
//   "words": [
//     {
//       "word": "The",
//       "start": 0.0,
//       "end": 0.23999999463558197
//     },
//     ...
//     {
//       "word": "volleyball",
//       "start": 7.400000095367432,
//       "end": 7.900000095367432
//     }
//   ]
// }


```

##### Segments Timestamps Request
```go
// NEW: Text To Speech API - Converts text to audio

// Transcription - Segments Timestamps Request
resp, err := gptClient.OpenAISpeechToTextSegmentTimestamps(&req)
if err != nil {
    log.Println("error stt req: " + err.Error())
} else {
    log.Println("\nstt SEGMENTS TIMESTAMPS response: ")
    log.Println(resp)
}

// example transcript response for segments timestamps transcription
// structure of json/ struct will be similar like openai docs structure
// {
//   "task": "transcribe",
//   "language": "english",
//   "duration": 8.470000267028809,
//   "text": "The beach was a popular spot on a hot summer day. People were swimming in the ocean, building sandcastles, and playing beach volleyball.",
//   "segments": [
//     {
//       "id": 0,
//       "seek": 0,
//       "start": 0.0,
//       "end": 3.319999933242798,
//       "text": " The beach was a popular spot on a hot summer day.",
//       "tokens": [
//         50364, 440, 7534, 390, 257, 3743, 4008, 322, 257, 2368, 4266, 786, 13, 50530
//       ],
//       "temperature": 0.0,
//       "avg_logprob": -0.2860786020755768,
//       "compression_ratio": 1.2363636493682861,
//       "no_speech_prob": 0.00985979475080967
//     },
//     ...
//   ]
// }


```


#### Translation Example Usage
Convert audio to text and translate to to english
```go
// NEW: Speech to Text API (Translation) - Convert audio to text and translate to to english

// input data can use with format like form, local filepath, or io.reader format

// filepath example
req := &OATranslationDefaultReq{
    File:     "/path/to/french_audio.mp3",
}

// file open/ io.reader as input
file, err := os.Open("/path/to/french_audio.mp3")
if err != nil {
    log.Fatalf("Failed to open file: %v", err)
}
defer file.Close()

req := &OATranscriptionDefaultReq{
    File:     file,
    Filename: "french_audio.mp3", // Required when using io.Reader
    Language: "en",
}

if ttsData, err := gptClient.OpenAITextToSpeech(&ttsReBody); err != nil {
    fmt.Println("error tts req: " + err.Error())
} else {
    fmt.Println("\ntts response: ")
    fmt.Println(ttsData)
}

resp, err := openAIClient.OpenAISpeechToTextTranslation(req)
if err != nil {
    log.Fatalf("Audio translation failed: %v", err)
}
fmt.Println("English Translation:", resp.Text)

// example real transcript from german
// Heute ist Donnerstag. Am Nachmittag gehe ich in den Supermarkt, um einzukaufen. Ich brauche Milch, fünf Eier, Brot und Obst. Im Supermarkt finde ich alles, was ich brauche. Zuerst gehe ich zur Obstabteilung und nehme frische Orangen. Dann hole ich Milch, Eier und Brot. An der Kasse bezahle ich oft mit meiner Kreditkarte, manchmal bar. Danach gehe ich nach Hause. Einkaufen macht Spaß.

// will be translate to
// Today is Thursday. In the afternoon I go to the supermarket to buy groceries. I need milk, five eggs, bread and fruit. In the supermarket I find everything I need. First I go to the fruit section and take fresh oranges. Then I get milk, eggs and bread. At the checkout I often pay with my credit card, sometimes cash. Then I go home. Shopping is fun!

```

## Error Handling

The current wrapper does not yet provide specific error types for all scenarios. We plan to expand error coverage in future updates. However, for now, we provide error handling for nearly all available functions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Roadmap

- [ ] Add support for more functions
- [ ] Comprehensive error type system

