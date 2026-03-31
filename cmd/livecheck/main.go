package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	seaavey "github.com/seaavey/seaavey-sdk-go"
)

func main() {
	apiKey := os.Getenv("SEAAVEY_API_KEY")
	targetURL := os.Getenv("SEAAVEY_TARGET_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := seaavey.NewClient(apiKey)

	resp, err := client.Downloader.TikTok(ctx, targetURL)
	if err != nil {
		log.Fatalf("live check failed: %v", err)
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("marshal response: %v", err)
	}

	fmt.Println(string(out))
}
