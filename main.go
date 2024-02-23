package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Game struct represents the data structure for each game
type Game struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Price       string `json:"price"`
	ReleaseDate string `json:"release_date"`
	Reviews     string `json:"reviews"`
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
	err = csvWriter.Write([]string{"Title", "Link", "Price", "Release Date", "Reviews"})
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write game data to CSV
	for _, game := range games {
		err := csvWriter.Write([]string{game.Title, game.Link, game.Price, game.ReleaseDate, game.Reviews})
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
	doc.Find("#search_resultsRows > a").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title").Text()
		link, _ := s.Attr("href") // Extract link directly from the anchor tag
		price := s.Find(".col.search_price_discount_combined .discount_final_price").Text()
		releaseDate := s.Find(".search_released").Text()

		// Extract review summary
		reviewsSummary := ""
		if tooltipHTML, exists := s.Find(".search_review_summary").Attr("data-tooltip-html"); exists {
			// Split the HTML string by the <br> tag
			parts := strings.Split(tooltipHTML, "<br>")
			if len(parts) > 0 {
				// Extract only the first part
				reviewsSummary = strings.TrimSpace(parts[0])
			}
		}

		game := Game{
			Title:       strings.TrimSpace(title),
			Link:        link,
			Price:       strings.TrimSpace(price),
			ReleaseDate: strings.TrimSpace(releaseDate),
			Reviews:     reviewsSummary,
		}
		games = append(games, game)
	})

	return games, nil
}






