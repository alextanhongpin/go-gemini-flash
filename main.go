package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GOOGLE_GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with most use cases
	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text("You are a cat. Your name is Neko.")},
	}
	b, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /button", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		fmt.Fprintf(w, `
  <div hx-ext="sse" sse-connect="/event?prompt=%s" sse-swap="message" hx-swap='beforeend'>
      Contents of this box will be updated in real time
  </div>
`, r.FormValue("prompt"))
	})
	mux.HandleFunc("GET /event", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ctx := r.Context()

		prompt := r.URL.Query().Get("prompt")
		iter := model.GenerateContentStream(ctx, genai.Text(prompt))
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			for _, c := range resp.Candidates {
				for _, p := range c.Content.Parts {
					fmt.Fprintf(w, "data: %s\n\n", p)
					w.(http.Flusher).Flush()
				}
			}
		}

		resp := iter.MergedResponse()
		if resp.PromptFeedback.BlockReason != genai.BlockReasonUnspecified {
			// Don't save
		}
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))

		// Simulate closing the connection
		closeNotify := w.(http.CloseNotifier).CloseNotify()
		<-closeNotify
	})
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("listening to port *:3000")
	panic(http.ListenAndServe(":3000", mux))
}
