package function

import (
	"encoding/json"
	"log"
)

type inference struct {
	Score float64 `json:"score"`
	Name  string  `json:"name"`
}

type payload struct {
	URL        string      `json:"url"`
	Inferences []inference `json:"inferences"`
}

// OK string to be returned
const OK = "OK"

// Handle a serverless request
func Handle(req []byte) string {

	d := payload{}
	if err := json.Unmarshal(req, &d); err != nil {
		log.Fatalf("Unable to unmarshal object for: %s", string(req))
	}

	sendMessage(d.URL, d.Inferences)

	return OK
}
