package collector

import (
	"fmt"
	"internet-availability-gr/models"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func NewCollector(stage, stateId, municipalityId *int) *colly.Collector {
	c := colly.NewCollector()
	var mu sync.Mutex

	c.OnHTML("li > a[id]", func(e *colly.HTMLElement) {
		id, err := strconv.Atoi(e.Attr("id"))
		if err != nil || id == 0 {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		switch *stage {
		case 1:
			states = append(states, models.State{Id: id, Name: e.Text})
		case 2:
			municipalities = append(municipalities, models.Municipality{
				Id:      id,
				Name:    e.Text,
				StateId: *stateId,
			})
		case 3:
			prefectures = append(prefectures, models.Prefecture{
				Id:             id,
				Name:           e.Text,
				StateId:        *stateId,
				MunicipalityId: *municipalityId,
			})
		}
	})

	c.OnHTML(".available-programm-container.res", func(e *colly.HTMLElement) {
		name := e.DOM.Find(".main-desc").First().Text()
		speed := e.DOM.Find(".secondary-desc").Text()

		packages = append(packages, models.PackageInfo{
			Category: "residential",
			Name:     strings.TrimSpace(name),
			Speed:    speed,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %s\n", err)
		if r.StatusCode == 0 {
			log.Println("Retrying:", r.Request.URL)
			time.Sleep(10 * time.Second)
			r.Request.Retry()
		}
	})

	return c
}

var states []models.State
var municipalities []models.Municipality
var prefectures []models.Prefecture
var packages []models.PackageInfo

func FetchStates(c *colly.Collector, url string, output *[]models.State) {
	err := c.Visit(url)
	if err != nil {
		log.Fatal("Error fetching states: ", err)
	}
	*output = states
}

func FetchMunicipalities(c *colly.Collector, baseURL string, states []models.State, output *[]models.Municipality, stateId *int) {
	var wg sync.WaitGroup
	for _, state := range states {
		wg.Add(1)
		go func(s models.State) {
			defer wg.Done()
			*stateId = s.Id
			err := c.Visit(baseURL + "?stateId=" + strconv.Itoa(s.Id))
			if err != nil {
				log.Printf("Error fetching municipalities for state %d: %v\n", s.Id, err)
			}
		}(state)
	}
	wg.Wait()
	*output = municipalities
}

func FetchPrefectures(c *colly.Collector, baseURL string, municipalities []models.Municipality, output *[]models.Prefecture, stateId, municipalityId *int) {
	var wg sync.WaitGroup
	for i, mun := range municipalities {
		if i > 0 && i%50 == 0 {
			log.Println("Iteration #" + strconv.Itoa(i) + ": Pausing for 5 seconds to avoid overloading the server...")
			time.Sleep(5 * time.Second)
		}
		wg.Add(1)
		go func(m models.Municipality) {
			defer wg.Done()
			*stateId = m.StateId
			*municipalityId = m.Id
			err := c.Visit(baseURL + "?stateId=" + strconv.Itoa(m.StateId) + "&municipalityId=" + strconv.Itoa(m.Id))
			if err != nil {
				log.Printf("Error fetching prefectures for municipality %d: %v\n", m.Id, err)
			}
		}(mun)
	}
	wg.Wait()
	*output = prefectures
}

func FetchPackagesViaStreet(c *colly.Collector, baseURL string, streetName string, stateName string, municipalityName string, prefectureName string, number int, output *[]models.PackageInfo) {
	u := fmt.Sprintf("%s?Accept-Language=en-US,en;q=0.9,el-GR;q=0.8,el;q=0.7&mTelno=&mAddress=%s&mState=%s&mPrefecture=%s&mNumber=%d&mArea=%s&searchcriteria=address&ct=res",
		baseURL,
		url.QueryEscape(streetName),
		url.QueryEscape(stateName),
		url.QueryEscape(municipalityName),
		number,
		url.QueryEscape(prefectureName),
	)
	println(u)
	err := c.Visit(u)
	if err != nil {
		log.Printf("Error fetching packages for streetName %s and number %d: %v\n", streetName, number, err)
	}

	*output = packages
}

func FetchPackagesViaTelephone(c *colly.Collector, baseURL string, telephone string, output *[]models.PackageInfo) {
	u := fmt.Sprintf("%s?Accept-Language=en-US,en;q=0.9,el-GR;q=0.8,el;q=0.7&mTelno=%s&searchcriteria=tel&ct=res",
		baseURL,
		url.QueryEscape(telephone),
	)

	err := c.Visit(u)
	if err != nil {
		log.Printf("Error fetching packages for telephone %s: %v\n", telephone, err)
	}

	for index := range packages {
		packages[index].Telephone = telephone
	}

	*output = packages
}
