package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pborman/uuid"
	btrdb "gopkg.in/btrdb.v4"
)

func printRange(ctx context.Context, stream *btrdb.Stream, start int64, end int64, version uint64) (uint64, error) {
	ptChan, verChan, errChan := stream.RawValues(ctx, start, end, version)
	for pt := range ptChan {
		timeString := time.Unix(0, pt.Time).Format(time.UnixDate)
		fmt.Printf("%v: %v\n", timeString, pt.Value)
	}
	ver := <-verChan
	err := <-errChan
	if err != nil {
		return 0, err
	}
	return ver, nil
}

func tail(endpoint string, streamUUID string, follow bool, last time.Duration) {
	ctx := context.Background()
	conn, err := btrdb.Connect(ctx, endpoint)
	if err != nil {
		fmt.Printf("Could not connect to server: %v\n", err)
		os.Exit(1)
	}

	parsed := uuid.Parse(streamUUID)
	stream := conn.StreamFromUUID(parsed)
	exists, err := stream.Exists(ctx)
	if err != nil {
		fmt.Printf("An unknown error occured: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Printf("Stream %v not found.\n", streamUUID)
		os.Exit(1)
	}

	version, err := stream.Version(ctx)
	if err != nil {
		fmt.Printf("Error trying to get stream version: %v\n", err)
		os.Exit(1)
	}

	end := time.Now()
	start := end.Add(-last)
	version, err = printRange(ctx, stream, start.UnixNano(), end.UnixNano(), version)
	if err != nil {
		fmt.Printf("Error printing from BTrDB stream: %v\n", err)
		os.Exit(0)
	}

	if follow {
		for {
			start = end
			end = time.Now()
			version, err = printRange(ctx, stream, start.UnixNano(), end.UnixNano(), version)
			if err != nil {
				fmt.Printf("Error printing from BTrDB stream: %v\n", err)
				os.Exit(0)
			}

			fmt.Println("sleepy")
			time.Sleep(10)
		}
	}
}
