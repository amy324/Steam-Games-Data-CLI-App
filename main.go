package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	//"strings"

	"github.com/PuerkitoBio/goquery"
)

// Game struct represents the data structure for each game
type Game struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Price       string `json:"price"`
	ReleaseDate string `json:"release_date"`
}

func main() {
	// Get game keyword input from user
	var gameKeyword string
	fmt.Print("Enter the game keyword you want to search: ")
	fmt.Scanln(&gameKeyword)

	// Create a folder to hold the extracted files
	err := os.Mkdir("resultfiles", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Scrape data from the first page
	url := fmt.Sprintf("https://store.steampowered.com/search/?term=%s", gameKeyword)
	games, err := scrapePage(url)
	if err != nil {
		fmt.Println("Error scraping page:", err)
		return
	}

	// Output data to JSON file
	jsonFile, err := os.Create("resultfiles/games.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer jsonFile.Close()

	jsonEncoder := json.NewEncoder(jsonFile)
	jsonEncoder.SetIndent("", "  ")
	err = jsonEncoder.Encode(games)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON file created successfully.")

	// Output data to CSV file
	csvFile, err := os.Create("resultfiles/games.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	err = csvWriter.Write([]string{"Title", "Link", "Price", "Release Date"})
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write game data to CSV
	for _, game := range games {
		err := csvWriter.Write([]string{game.Title, game.Link, game.Price, game.ReleaseDate})
		if err != nil {
			fmt.Println("Error writing CSV record:", err)
			return
		}
	}

	fmt.Println("CSV file created successfully.")
}
func scrapePage(url string) ([]Game, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var games []Game
	doc.Find(".search_result_row").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title").Text()
		link, _ := s.Find(".search_name a").Attr("href")
		price := s.Find(".search_price").Text()
		releaseDate := s.Find(".search_released").Text()

		game := Game{
			Title:       strings.TrimSpace(title),
			Link:        link,
			Price:       strings.TrimSpace(price),
			ReleaseDate: strings.TrimSpace(releaseDate),
		}
		games = append(games, game)
	})

	return games, nil
}
