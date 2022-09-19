package main

import (
	"fmt"
	"github.com/bitfield/script"
	"strings"
)

const schedulesConfigPath = "/config/schedules"
const osmConfigPath = "/config/osm"

const schedulesDataPath = "/input/schedule"
const osmDataFolder = "/input"

func main() {
	scheduleUrls, err := script.File(schedulesConfigPath).Slice()
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
	script.ListFiles(tmpDir).ExecForEach(fmt.Sprintf("unzip {{.}} -d %v", schedulesDataPath)).Stdout()

	mapUrls, err := script.File(osmConfigPath).Slice()
	if err != nil {
		panic(fmt.Errorf("error reading osm URLs: %w", err))
	}

	fmt.Println("Downloading OpenStreetMap data...")
	for _, url := range mapUrls {
		script.Exec(fmt.Sprintf("wget -P %v %v", osmDataFolder, url)).Wait()
	}
}
