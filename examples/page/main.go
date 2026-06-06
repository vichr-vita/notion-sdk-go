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

	pageID := os.Getenv("NOTION_PAGE_ID")
	if pageID == "" {
		log.Fatal("NOTION_PAGE_ID must be set")
	}

	client := notion.NewClient(token)

	page, err := client.Pages.Get(context.Background(), pageID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("page: %s\n", page.ID)

	markdown, err := client.Pages.GetMarkdown(context.Background(), pageID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(markdown.Markdown)
}
