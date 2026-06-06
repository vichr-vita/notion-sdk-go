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

	blockID := os.Getenv("NOTION_BLOCK_ID")
	if blockID == "" {
		log.Fatal("NOTION_BLOCK_ID must be set")
	}

	client := notion.NewClient(token)

	err := client.Blocks.ForEachChild(context.Background(), blockID, &notion.Pagination{
		PageSize: 100,
	}, func(block notion.Block) error {
		fmt.Printf("%s %s\n", block.Type, block.ID)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
