package main

import (
	"context"
	"fmt"
	"log"
	"os"

	notion "github.com/vichr-vita/notion-sdk-go"
)

func main() {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		log.Fatal("NOTION_TOKEN must be set")
	}

	client := notion.NewClient(token)

	resp, err := client.Search.Query(context.Background(), notion.SearchRequest{
		Query:    "Tasks",
		PageSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("%s %s\n", result.Object, result.ID)
	}
}
