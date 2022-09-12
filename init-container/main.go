package main

import (
	"fmt"
	"github.com/bitfield/script"
	"strings"
)

func main() {
	scheduleUrls, err := script.File("schedules.txt").Slice()
	if err != nil {
		panic(fmt.Errorf("error reading schedule URLs: %w", err))
	}

	fmt.Println("Downloading schedules...")
	tmpDir, err := script.Exec("mktemp -d").String()
	if err != nil {
		panic(fmt.Errorf("error creating tmp directory: %w", err))
	}
	tmpDir = strings.Replace(tmpDir, "\n", "", -1)
	for _, url := range scheduleUrls {

		script.Exec(fmt.Sprintf("wget -P %v %v", tmpDir, url)).Wait()
	}
	script.Echo(fmt.Sprintf("Downloaded files to %v", tmpDir)).Stdout()

	_, err = script.ListFiles(tmpDir).Stdout()
	if err != nil {
		fmt.Printf("Error listing files: %v", err)
	}
	script.ListFiles(tmpDir).ExecForEach("unzip {{.}} -d schedules").Stdout()

	mapUrls, err := script.File("osm.txt").Slice()
	if err != nil {
		panic(fmt.Errorf("error reading osm URLs: %w", err))
	}

	fmt.Println("Downloading OpenStreetMap data...")
	for _, url := range mapUrls {
		script.Exec(fmt.Sprintf("wget %v", url)).Wait()
	}
}
