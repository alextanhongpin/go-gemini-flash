// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//var resp genai.GenerateContentResponse
	var resp Response
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		panic(err)
	}
	fmt.Println(resp.Candidates[0].Content.Parts[0])
}

type Response struct {
	Candidates []struct {
		Content struct {
			Parts []string
		}
	}
}

var raw = `
{
  "Candidates": [
    {
      "Index": 0,
      "Content": {
        "Parts": [
          "\"Gemini\" can refer to a few different things, so to give you the most helpful answer, I need a little more context. Could you tell me what you're interested in learning about? \n\nFor example, are you asking about:\n\n* **The Gemini constellation?** This is a constellation in the Northern Hemisphere, known for its twin stars Castor and Pollux.\n* **The Gemini spacecraft?** This was a series of American crewed spacecraft used in the 1960s.\n* **The Gemini language model?** This is a large language model developed by Google AI.\n* **The Gemini Project?** This was a code name for a project during World War II.\n\nOnce I know what you're interested in, I can give you a more specific and detailed answer. \n"
        ],
        "Role": "model"
      },
      "FinishReason": 1,
      "SafetyRatings": [
        {
          "Category": 9,
          "Probability": 1,
          "Blocked": false
        },
        {
          "Category": 8,
          "Probability": 1,
          "Blocked": false
        },
        {
          "Category": 7,
          "Probability": 1,
          "Blocked": false
        },
        {
          "Category": 10,
          "Probability": 1,
          "Blocked": false
        }
      ],
      "CitationMetadata": null,
      "TokenCount": 0
    }
  ],
  "PromptFeedback": null,
  "UsageMetadata": {
    "PromptTokenCount": 67,
    "CandidatesTokenCount": 1,
    "TotalTokenCount": 68
  }
}
`
