package llama

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/r3labs/sse/v2"
)

type Client struct {
	url    string
	apiKey string
	client *http.Client
}

type completionPayload struct {
	Stream      bool     `json:"stream"`
	Grammar     string   `json:"grammar"`
	CachePrompt bool     `json:"cache_prompt"`
	ApiKey      string   `json:"api_key"`
	Prompt      string   `json:"prompt"`
	SlotId      int      `json:"slot_id"`
	Stop        []string `json:"stop"`
	CompletionOptions
}

type CompletionOptions struct {
	NPredict         int     `json:"n_predict"`
	NProbs           int     `json:"n_probs"`
	MirostatEta      float64 `json:"mirostat_eta"`
	Temperature      float64 `json:"temperature"`
	RepeatLastN      float64 `json:"repeat_last_n"`
	RepeatPenalty    float64 `json:"repeat_penalty"`
	TopK             int     `json:"top_k"`
	TopP             float64 `json:"top_p"`
	MinP             float64 `json:"min_p"`
	TfsZ             float64 `json:"tfs_z"`
	TypicalP         float64 `json:"typical_p"`
	PresencePenalty  float64 `json:"presence_penalty"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	Mirostat         float64 `json:"mirostat"`
	MirostatTau      float64 `json:"mirostat_tau"`
}

var Default = CompletionOptions{
	NPredict:         128,
	NProbs:           0,
	MirostatEta:      0,
	Temperature:      0.2,
	RepeatLastN:      256,
	RepeatPenalty:    1.18,
	TopK:             40,
	TopP:             0.95,
	MinP:             0.05,
	TfsZ:             1,
	TypicalP:         1,
	PresencePenalty:  0,
	FrequencyPenalty: 0,
	Mirostat:         0,
	MirostatTau:      5,
}

type streamResponse struct {
	Content string `json:"content"`
	Stop    bool   `json:"stop"`
}

func NewClient(url, apiKey string) *Client {
	return &Client{
		url:    url,
		apiKey: apiKey,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) Completion(prompt string, stop []string, opts CompletionOptions) chan string {
	payload := completionPayload{
		Stream:            true,
		ApiKey:            c.apiKey,
		Prompt:            prompt,
		Stop:              stop,
		SlotId:            -1,
		CompletionOptions: opts,
	}

	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json; charset=utf8")
	req.Header.Set("Accept", "text/event-stream")

	ch := make(chan string)

	go func() {
		defer close(ch)

		rsp, err := c.client.Do(req)
		if err != nil {
			log.Println("Completion request failed: ", err)
			return
		}

		if rsp.StatusCode != http.StatusOK {
			log.Println("Completion endpoint returned non-200 status code: ", rsp.StatusCode)
			return
		}

		defer rsp.Body.Close()
		reader := sse.NewEventStreamReader(rsp.Body, 16384)

		for {
			msg, err := reader.ReadEvent()
			if err != nil {
				log.Println("Read event error: ", err)
				return
			}
			streamRsp := streamResponse{}
			if err := json.Unmarshal(msg[6:], &streamRsp); err != nil {
				log.Println("Unmarshal event message error: ", err)
				return
			}

			if streamRsp.Stop {
				return
			}

			ch <- streamRsp.Content
		}

	}()

	return ch
}
