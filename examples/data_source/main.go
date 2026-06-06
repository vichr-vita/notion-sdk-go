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

	dataSourceID := os.Getenv("NOTION_DATA_SOURCE_ID")
	if dataSourceID == "" {
		log.Fatal("NOTION_DATA_SOURCE_ID must be set")
	}

	client := notion.NewClient(token)

	resp, err := client.DataSources.Query(context.Background(), dataSourceID, notion.QueryDataSourceRequest{
		PageSize: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("%s %s\n", result.Object, result.ID)
	}
}
