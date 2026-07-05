package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const upstreamURL = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"

var thinkingInject = json.RawMessage(`{"thinking":{"type":"disabled"}}`)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{
		Timeout: 300 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 128,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      http.HandlerFunc(handler(client)),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 310 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("gateway listening on :%s, upstream %s", port, upstreamURL)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handler(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST is supported", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		modified, err := injectThinking(body)
		if err != nil {
			http.Error(w, "failed to process body", http.StatusInternalServerError)
			return
		}

		upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, bytes.NewReader(modified))
		if err != nil {
			http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
			return
		}

		upstreamReq.Header.Set("Content-Type", "application/json")
		if auth := r.Header.Get("Authorization"); auth != "" {
			upstreamReq.Header.Set("Authorization", auth)
		}

		resp, err := client.Do(upstreamReq)
		if err != nil {
			log.Printf("upstream error: %v", err)
			http.Error(w, "upstream request failed", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, vs := range resp.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func injectThinking(body []byte) ([]byte, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	m["thinking"] = json.RawMessage(`{"type":"disabled"}`)
	m["stream"] = json.RawMessage(`false`)
	return json.Marshal(m)
}
