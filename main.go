package main

import (
	"internet-availability-gr/collector"
	"internet-availability-gr/models"
	"internet-availability-gr/utils"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	crawlBaseURL := os.Getenv("CRAWL_BASE_URL")
	crawlEnabled, _ := strconv.ParseBool(os.Getenv("CRAWL_ENABLED"))

	querylBaseURL := os.Getenv("QUERY_BASE_URL")

	stage := new(int)
	stateId := new(int)
	municipalityId := new(int)

	c := collector.NewCollector(stage, stateId, municipalityId)

	var states []models.State
	var municipalities []models.Municipality
	var prefectures []models.Prefecture
	var packages []models.PackageInfo

	if crawlEnabled {
		log.Println("Crawling started.")

		log.Println("Fetching states...")
		*stage = 1
		collector.FetchStates(c, crawlBaseURL, &states)
		utils.WriteJSON("data/states.json", states)
		log.Println("Completed states fetching...")

		log.Println("Fetching municipalities...")
		*stage = 2
		collector.FetchMunicipalities(c, crawlBaseURL, states, &municipalities, stateId)
		utils.WriteJSON("data/municipalities.json", municipalities)
		log.Println("Completed municipalities fetching...")

		log.Println("Fetching prefectures...")
		*stage = 3
		collector.FetchPrefectures(c, crawlBaseURL, municipalities, &prefectures, stateId, municipalityId)
		utils.WriteJSON("data/prefectures.json", prefectures)
		log.Println("Completed prefectures fetching...")

		log.Println("Crawling complete.")
	}

	data, err := os.ReadFile("input.yaml")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var inputList models.Input
	err = yaml.Unmarshal(data, &inputList)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	for _, tel := range inputList.Telephones {
		var pck []models.PackageInfo
		collector.FetchPackagesViaTelephone(c, querylBaseURL, tel, &pck)

		packages = append(packages, pck...)
	}

	utils.WriteJSON("data/packages.json", packages)
}
