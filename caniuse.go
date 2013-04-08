package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var remoteFile = "http://raw.github.com/Fyrd/caniuse/master/data.json"
var localFile = "data.json"

type StatsType map[string]map[string]string
type UsgGblType map[string]string

type DataObject struct {
	Eras     map[string]string      "eras"
	Agents   map[string]BrowserType "agents"
	Statuses map[string]string      "statuses"
	Cats     map[string][]string    "cats"
	Data     map[string]DataType    "data"
}

type BrowserType struct {
	Browser        string     "browser"
	Abbr           string     "abbr"
	Prefix         string     "prefix"
	Type           string     "type"
	UsageGlobal    UsgGblType "usage_global"
	Versions       []*string  "versions"
	CurrentVersion string     "current_version"
}

type DataType struct {
	Title       string "title"
	Description string "description"
	Spec        string "spec"
	Status      string "status"
	Links       []struct {
		Url   string "url"
		Title string "title"
	} "links"
	Categories []*string "categories"
	Stats      StatsType "stats"
	Notes      string    "notes"
	UsagePercY string    "usage_perc_y"
	UsagePercA string    "usage_perc_a"
	UcPrefix   bool      "ucprefix"
	Parent     string    "parent"
	Keywords   string    "keywords"
}

//Load data from file
func loadData(fileLocation string) []byte {
	
	var content []byte
	var errFile error
	
	if fileLocation=="local" {
		content, errFile = ioutil.ReadFile(localFile)
		
		if errFile != nil {
			log.Print("File error", errFile)
		}

	} else {
		resp, errHttp := http.Get(remoteFile)
		defer resp.Body.Close()

		if errHttp != nil {
			log.Fatal("Error downloading file: %v\n", errHttp)
		}

		content, errFile = ioutil.ReadAll(resp.Body)

		if errFile != nil {
			log.Fatal("File error: %v\n", errFile)
		}
	}

	return content
}

//Parse data
func parseData(file []byte) *DataObject {
	var data DataObject
	
	errParse := json.Unmarshal(file, &data)

	if errParse != nil {
		log.Print("Parser error:", errParse)
	}
	
	return &data
}

func main() {
	//Remote file or Local file
	var file = loadData("local")
	var data = parseData(file)

	feature:=os.Args[1]

	switch feature {
		default:
			//Details for feature
			fmt.Println("----------------");
			fmt.Printf("Details for: %v \n", feature);
			fmt.Println("----------------");
			fmt.Printf("- Title: %v  \n", data.Data[feature].Title);
			fmt.Printf("- Description: %v  \n", data.Data[feature].Description);
			fmt.Printf("- Spec: %v  \n", data.Data[feature].Spec);
			fmt.Printf("- Status: %v  \n", data.Statuses[data.Data[feature].Status]);
			fmt.Printf("- Categories: %v  \n", data.Cats[*data.Data[feature].Categories[0]]); //TODO: for range
			
			fmt.Println("- Stats:");
			
			//Stats for every agent
			for agent, stats := range data.Data[feature].Stats {
				fmt.Printf("-- %v | %v \n", data.Agents[agent].Browser, data.Agents[agent].Type);
				
				//Stats for every agent-version
				for version, status := range stats {
					fmt.Printf("--- Version: %v | Status: %v  \n", version, status);
				}
			}
			
			fmt.Printf("- Notes: %v  \n", data.Data[feature].Notes);

		case "--list":
			for key, value := range data.Data	{
				fmt.Printf("- %v | %v \n", key, value.Title)
			}
	}
}
