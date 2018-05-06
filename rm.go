package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pborman/uuid"
	btrdb "gopkg.in/btrdb.v4"
)

// TODO: support deleting by collection and tags

func doObliterate(ctx context.Context, stream *btrdb.Stream) {
	err := stream.Obliterate(ctx)
	if err != nil {
		fmt.Printf("Error deleting stream: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func remove(endpoint string, streamUuid string, confirmed bool) {
	ctx := context.Background()
	conn, err := btrdb.Connect(ctx, endpoint)
	if err != nil {
		fmt.Printf("Could not connect to server: %v\n", err)
		os.Exit(1)
	}

	parsed := uuid.Parse(streamUuid)
	stream := conn.StreamFromUUID(parsed)
	exists, err := stream.Exists(ctx)
	if err != nil {
		fmt.Printf("An unknown error occured: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("Stream %v not found.\n", streamUuid)
		os.Exit(1)
	}

	collection, err := stream.Collection(ctx)
	if err != nil {
		fmt.Printf("Problem getting collection from stream: %v", err)
		os.Exit(1)
	}

	if confirmed {
		doObliterate(ctx, stream)
	}

	var response string
	for {
		fmt.Printf("Are you sure you want to delete stream %v from collection %v? [y/n] ", streamUuid, collection)
		_, err := fmt.Scanln(&response)
		if err != nil {
			fmt.Print("There was a problem getting a response. Are you sure you want to delete? [y/n]")
			continue
		}
		response = strings.ToLower(response)
		if strings.Contains(response, "y") {
			doObliterate(ctx, stream)
		} else if strings.Contains(response, "n") {
			fmt.Println("Quitting")
			os.Exit(0)
		} else {
			fmt.Println("Please answer with y or n.")
		}
	}
}
