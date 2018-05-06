package main

import (
	"context"
	"fmt"
	"os"

	btrdb "gopkg.in/btrdb.v4"
)

func printRow(longest int, first string, second string) {
	width := fmt.Sprintf("%v", longest+1)
	fmt.Printf("%-"+width+"v%v\n", first, second)
}

func printCollectionsList(conn *btrdb.BTrDB, collections []string) {
	if len(collections) == 0 {
		fmt.Println("No collections found!")
		return
	}

	longest := 0
	for _, col := range collections {
		if len(col) > longest {
			longest = len(col)
		}
	}
	printRow(longest, "Collection name", "Stream count")
	for _, collection := range collections {
		streams, err := conn.LookupStreams(context.Background(), collection, false, nil, nil)
		streamsCount := ""
		if err != nil {
			streamsCount = fmt.Sprintf("Error getting streams for this collection %v", err)
		} else {
			streamsCount = fmt.Sprintf("%v", len(streams))
		}
		printRow(longest, collection, streamsCount)
	}
}

func printCollectionDetails(conn *btrdb.BTrDB, collection string) {
	streams, err := conn.LookupStreams(context.Background(), collection, false, nil, nil)
	if err != nil {
		fmt.Printf("Error finding streams for %v: %v\n", collection, err)
		return
	}
	fmt.Println("Collection: " + collection + ":")
	fmt.Println("Streams:")
	for _, stream := range streams {
		fmt.Println(" * UUID: " + stream.UUID().String())

		fmt.Println(" * Tags: ")
		tags, err := stream.Tags(context.Background())
		if err != nil {
			fmt.Printf("     - Error getting tags: %v\n", err)
		} else if len(tags) == 0 {
			fmt.Println("     - None")
		} else {
			for k, v := range tags {
				fmt.Printf("     - %v: %v\n", k, v)
			}
		}

		fmt.Println(" * Annontations: ")
		ann, _, err := stream.Annotations(context.Background())
		if err != nil {
			fmt.Printf("     - Error getting annotations: %v\n", err)
		} else if len(ann) == 0 {
			fmt.Println("     - None")
		} else {
			for k, v := range ann {
				fmt.Printf("     - %v: %v\n", k, v)
			}
		}
		fmt.Println()
	}
}

func list(endpoint string, prefix string) {
	conn, err := btrdb.Connect(context.Background(), endpoint)
	if err != nil {
		fmt.Printf("Could not connect to server: %v\n", err)
		os.Exit(1)
	}

	collections, err := conn.ListCollections(context.Background(), prefix)
	if err != nil {
		fmt.Printf("Error listing collectiosn: %v\n", err)
		os.Exit(1)
	}

	if len(collections) == 1 {
		printCollectionDetails(conn, collections[0])
	} else {
		printCollectionsList(conn, collections)
	}
}
