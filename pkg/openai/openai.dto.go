package openai

// OPEN AI DOCS api Reference
// https://platform.openai.com/docs/api-reference/chat/create

// OpenAIReqBodyMessageCompletion create document for openai dto
// if omitempty is used, the field will be omitted from the JSON representation of the object if the field has an empty value
// the omit empty value is optional in openai docs

// ----------------- CHAT COMPLETIONS ----------------------
type OAReqBodyMessageCompletion struct {
	Messages         interface{}            `json:"messages"` // required
	Model            string                 `json:"model"`    // required
	Store            bool                   `json:"store,omitempty"`
	Metadata         interface{}            `json:"metadata,omitempty"`
	FrequencyPenalty float64                `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`
	Logprobe         bool                   `json:"logprobe,omitempty"`
	Modalities       []string               `json:"modalities,omitempty"`
	ResponseFormat   map[string]interface{} `json:"response_format,omitempty"`
}

type OAMessageReq struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type OAContentVisionImageUrl struct {
	Url string `json:"url"`
}

type OAContentVisionBaseReq struct {
	Type     string                   `json:"type"`
	Text     *string                  `json:"text,omitempty"`
	ImageUrl *OAContentVisionImageUrl `json:"image_url,omitempty"`
}

// response COMPLETION OpenAI structure
type OAChatCompletionResp struct {
	ID                string     `json:"id"`
	Object            string     `json:"object"`
	Created           int64      `json:"created"`
	Model             string     `json:"model"`
	SystemFingerprint string     `json:"system_fingerprint"`
	Choices           []OAChoice `json:"choices"`
	Usage             OAUsage    `json:"usage"`
}

type OAChoice struct {
	Index        int       `json:"index"`
	Message      OAMessage `json:"message"`
	Logprobs     *string   `json:"logprobs"` // Could be null, so pointer
	FinishReason string    `json:"finish_reason"`
}

type OAMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// support for audio output gpt-4o-audio-preview
	Refusal string              `json:"refusal,omitempty"`
	Audio   OAAudioDataResponse `json:"audio,omitempty"`
}

type OAAudioDataResponse struct {
	Id         string `json:"id"`
	ExpiresAt  int64  `json:"expires_at"`
	Data       string `json:"data"`
	Transcript string `json:"transcript"`
}

type OAUsage struct {
	PromptTokens           int          `json:"prompt_tokens"`
	CompletionTokens       int          `json:"completion_tokens"`
	TotalTokens            int          `json:"total_tokens"`
	CompletionTokensDetail TokensDetail `json:"completion_tokens_details"`
}

type TokensDetail struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

// ----------------- DALL E IMAGE GENERATIONS ------ Reference for Image Generation Request Body
// 	   - OpenAI Docs: https://platform.openai.com/docs/api-reference/images/create
type OAReqImageGeneratorDallE struct {
	Prompt         string  `json:"prompt"`                    // required
	Model          string  `json:"model"`                     // required dall-e-2 or dall-e-3
	N              *int    `json:"n,omitempty"`               // total image to generate, max 10 default 1
	Quality        *string `json:"quality,omitempty"`         // "standard" (default), "hd" // just support for dall-e 3
	ResponseFormat *string `json:"response_format,omitempty"` // url (default) or b64_json
	Size           *string `json:"size,omitempty"`            // default "1024x1024",  Must be one of 256x256, 512x512, or 1024x1024 for dall-e-2. Must be one of 1024x1024, 1792x1024, or 1024x1792 for dall-e-3 models.
	Style          *string `json:"style,omitempty"`           // vivid (default) or natural, only support for dall-e-3
	User           *string `json:"user,omitempty"`            //A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
}

// response image create DALL e
type OAImageGeneratorDallEResp struct {
	Created int64                       `json:"created"`
	Data    []OAImageGeneratorDallEData `json:"data"`
}

type OAImageGeneratorDallEData struct {
	Url     string `json:"url"`      // if using response format url this data will contain the url image
	B64JSON string `json:"b64_json"` // if using response format b64_json this data will contain the base64 image
}

// ----------------- TTS TEXT TO SPEECH ------ Reference for TTS Request Body
// 	   - OpenAI Docs: https://platform.openai.com/docs/api-reference/audio/createSpeech
type OAReqTextToSpeech struct {
	Model          string   `json:"model"`           // required (tts-1 or tts-1-hd)
	Input          string   `json:"input"`           // required (max 4096)
	Voice          string   `json:"voice"`           // required (alloy, echo, fable, onyx, nova, and shimmer)
	ResponseFormat string   `json:"response_format"` // required (mp3, opus, aac, flac, wav, and pcm)
	Speed          *float64 `json:"speed,omitempty"` // optional (0.25 to 4.0. 1.0 is the default.)
}

type OATextToSpeechResp struct {
	FormatAudio string `json:"format_audio"` // will be like ".mp3"
	B64JSON     string `json:"b64_json"`
}

// ----------------- STT SPEECH TO TEXT ------ Reference for stt Request Body
// 	   - OpenAI Docs: https://platform.openai.com/docs/api-reference/audio/createTranscription

// transcription  -> transcribe audio to the input language
// will using inside the function
type OATranscriptionReq struct {
	File                   interface{} `json:"file" form:"file"`   // required
	Model                  string      `json:"model" form:"model"` // required
	Language               string      // optional, The language of the input audio. Supplying the input language in ISO-639-1 (e.g. en) format will improve accuracy and latency.
	Prompt                 string      `json:"prompt" form:"prompt"`                   // optional, An optional text to guide the model's style or continue a previous audio segment
	ResponseFormat         string      `json:"response_format" form:"response_format"` // default to json, The format of the response. Either json or text, srt, verbose_json, or vtt.
	Temperature            float64     `json:"temperature" form:"temperature"`         // The sampling temperature, between 0 and 1. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. If set to 0, the model will use log probability to automatically increase the temperature until certain thresholds are hit.
	TimestampGranularities []string    `json:"timestamp_granularities"`                // optional, The timestamp granularities to populate for this transcription. response_format must be set verbose_json to use timestamp granularities. Either or both of these options are supported: word, or segment. Note: There is no additional latency for segment timestamps, but generating word timestamps incurs additional latency.
}

// req parameter using different struct
type OATranscriptionDefaultReq struct {
	File        interface{} `json:"file" form:"file"`               // required
	Filename    string      `json:"filename" form:"filename"`       // optional, must be provided if File using io.reader as input
	Prompt      string      `json:"prompt" form:"prompt"`           // optional, An optional text to guide the model's style or continue a previous audio segment
	Language    string      `json:"language" form:"language"`       // optional, The language of the input audio. Supplying the input language in ISO-639-1 (e.g. en) format will improve accuracy and latency.
	Temperature float64     `json:"temperature" form:"temperature"` // The sampling temperature, between 0 and 1. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. If set to 0, the model will use log probability to automatically increase the temperature until certain thresholds are hit.
}

type OASTTError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    int    `json:"code"`
}

type OATranscriptionDefaultResp struct {
	Text  string     `json:"text"`
	Error OASTTError `json:"error"`
}

type wordTimestampResp struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type OATranscriptionWordTimestampResp struct {
	Task     string              `json:"task"`
	Language string              `json:"language"`
	Duration float64             `json:"duration"`
	Text     string              `json:"text"`
	Words    []wordTimestampResp `json:"words"`
	Error    OASTTError          `json:"error"`
}

type segmentResp struct {
	Id               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

type OATranscriptionSegmentResp struct {
	Task     string        `json:"task"`
	Language string        `json:"language"`
	Duration float64       `json:"duration"`
	Text     string        `json:"text"`
	Segments []segmentResp `json:"segments"`
	Error    OASTTError    `json:"error"`
}

type OATranslationDefaultReq struct {
	File        interface{} `json:"file" form:"file"`               // required
	Filename    string      `json:"filename" form:"filename"`       // optional, must be provided if File using io.reader as input
	Prompt      string      `json:"prompt" form:"prompt"`           // optional, An optional text to guide the model's style or continue a previous audio segment
	Temperature float64     `json:"temperature" form:"temperature"` // The sampling temperature, between 0 and 1. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic. If set to 0, the model will use log probability to automatically increase the temperature until certain thresholds are hit.
}
