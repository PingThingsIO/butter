package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jawher/mow.cli"
)

func listCommand(cmd *cli.Cmd) {
	cmd.Spec = "[ENDPOINT] [PREFIX]"
	endpoint := cmd.StringArg("ENDPOINT", "localhost:4410", "The BTrDB endpoint to list")
	prefix := cmd.StringArg("PREFIX", "", "A prefix to filter collections.")
	cmd.Action = func() {
		list(*endpoint, *prefix)
	}
}

func removeCommand(cmd *cli.Cmd) {
	cmd.Spec = "[ENDPOINT] UUID [-y]"
	endpoint := cmd.StringArg("ENDPOINT", "localhost:4410", "The BTrDB endpoint to list")
	uuid := cmd.StringArg("UUID", "", "UUID of the stream to delete")
	confirmed := cmd.BoolOpt("y yes", false, "Skip confirmation prompt")
	cmd.Action = func() {
		remove(*endpoint, *uuid, *confirmed)
	}
}

type streamConfigs []StreamConfig

func (s *streamConfigs) Set(v string) error {
	parts := strings.Split(v, ",")
	if len(parts) < 2 {
		return fmt.Errorf("Stream config requires src_collection and dest_collection, %v is invalid", v)
	}
	stream := StreamConfig{}
	stream.Tags = make(map[string]string)
	stream.SrcCollection = parts[0]
	stream.DstCollection = parts[1]
	tags := parts[2:len(parts)]
	for _, tag := range tags {
		tagParts := strings.Split(tag, "=")
		if len(tagParts) < 2 {
			return fmt.Errorf("Tag must follow format tagname=tagvalue, %v is invalid", tag)
		}
		tagName := tagParts[0]
		tagValue := tagParts[1]
		stream.Tags[tagName] = tagValue
	}
	*s = append(*s, stream)
	return nil
}

func (s *streamConfigs) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *streamConfigs) Clear() {
	*s = []StreamConfig{}
}

func copyCommand(cmd *cli.Cmd) {
	cmd.Spec = "FROMSERVER TOSERVER [-sea] STREAMCONFIG..."
	fromServer := cmd.StringArg("FROMSERVER", "", "BTrDB endpoint to copy from")
	toServer := cmd.StringArg("TOSERVER", "", "BTrDB endpoint to copy to")
	start := cmd.StringOpt("s start", "", "Start time of the range to copy (in format 2006-01-02T15:04:05+07:00)")
	end := cmd.StringOpt("e end", "", "End time of the range to copy (in format 2006-01-02T15:04:05+07:00)")
	abortIfExists := cmd.BoolOpt("a abort", false, "Abort the copy if the collection already exists")
	streamCfg := &streamConfigs{}
	cmd.VarArg("STREAMCONFIG", streamCfg, "Config for the streams to copy (follows the format src_collection,dest_collection,tagname=tagvalue)")
	cmd.Action = func() {
		cfg := Config{*fromServer, *toServer, *start, *end, *abortIfExists, *streamCfg}
		cp(cfg)
	}
}

type duration time.Duration

func (d *duration) Set(v string) error {
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return err
	}
	*d = duration(parsed)
	return nil
}

func (d *duration) String() string {
	duration := time.Duration(*d)
	return duration.String()
}

func tailCommand(cmd *cli.Cmd) {
	cmd.Spec = "[-fl] [ENDPOINT] UUID"
	endpoint := cmd.StringArg("ENDPOINT", "localhost:4410", "The BTrDB endpoint to print")
	uuid := cmd.StringArg("UUID", "", "UUID of the stream to print from")
	follow := cmd.BoolOpt("f follow", false, "Output values as they are added to the stream")
	last := duration(time.Second)
	cmd.VarOpt("l last", &last, "Duration decribing how far back to print (i.e. last 5m)")
	cmd.Action = func() {
		tail(*endpoint, *uuid, *follow, time.Duration(int64(last)))
	}
}

func main() {
	app := cli.App("butter", "Useful BTrDB CLI tools for development")

	app.Command("ls", "List collections for a BTrDB endpoint. If only one collection is returned, its streams will be listed.", listCommand)
	app.Command("rm", "Remove a stream from BTrDB", removeCommand)
	app.Command("cp", "Copy a collection from one BTrDB server to another", copyCommand)
	app.Command("tail", "Prints the latest values inserted into BTrDB", tailCommand)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
