package main
import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// Game struct represents the data structure for each game
type Game struct {
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Price       string   `json:"price"`
	ReleaseDate string   `json:"release_date"`
	Reviews     string   `json:"reviews"`
	Tags        []string `json:"tags"`
}

// TagsMap represents the structure of the tags JSON file
type TagsMap map[string]string

func main() {
	rootCmd := &cobra.Command{
		Use:   "steam-scraper",
		Short: "Scrape game data from Steam",
		Run: func(cmd *cobra.Command, args []string) {
			scrape()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func scrape() {
	// Load tags data from JSON file
	tagsFile, err := os.Open("data/tags.json")
	if err != nil {
		fmt.Println("Error opening tags file:", err)
		return
	}
	defer tagsFile.Close()

	var tags TagsMap
	err = json.NewDecoder(tagsFile).Decode(&tags)
	if err != nil {
		fmt.Println("Error decoding tags JSON:", err)
		return
	}

	// Create a folder to hold the extracted files
	err = os.Mkdir("resultfiles", 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Get birthtime from environment
	birthtime := os.Getenv("BIRTHTIME")

	for {
		// Get game keyword input from user
		var gameKeyword string
		fmt.Print("Enter the game keyword you want to search (or 'quit' to exit): ")
		fmt.Scanln(&gameKeyword)

		if gameKeyword == "quit" {
			break
		}

		// Scrape data from the first page
		url := fmt.Sprintf("https://store.steampowered.com/search/?term=%s", gameKeyword)
		games, err := scrapePage(url, tags)
		if err != nil {
			fmt.Println("Error scraping page:", err)
			continue
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
		jsonEncoder.SetEscapeHTML(false) // Prevent escaping HTML characters
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
		err = csvWriter.Write([]string{"Title", "Link", "Price", "Release Date", "Reviews", "Tags"})
		if err != nil {
			fmt.Println("Error writing CSV header:", err)
			return
		}

		// Write game data to CSV
		for _, game := range games {
			// Combine tags into a single string
			tagsStr := strings.Join(game.Tags, ", ")
			err := csvWriter.Write([]string{game.Title, game.Link, game.Price, game.ReleaseDate, game.Reviews, tagsStr})
			if err != nil {
				fmt.Println("Error writing CSV record:", err)
				return
			}
		}

		fmt.Println("CSV file created successfully.")

		// Prompt user to select a game for additional details
		var input string
		fmt.Print("Enter the title or link of the game for additional details: ")
		fmt.Scanln(&input)

		// Find the game in the JSON file
		var selectedGame *Game
		for _, g := range games {
			if g.Title == input || g.Link == input {
				selectedGame = &g
				break
			}
		}

		if selectedGame == nil {
			fmt.Println("Game not found.")
			continue
		}

		// Scrape additional details from the selected game's link
		if err := scrapeAdditionalDetails(selectedGame, birthtime); err != nil {
			fmt.Println("Error scraping additional details:", err)
			continue
		}
	}
}

func scrapePage(url string, tags TagsMap) ([]Game, error) {
	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check if response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML document: %v", err)
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

		// Extract tag IDs from HTML
		tagIDsStr, _ := s.Attr("data-ds-tagids")
		tagIDs := strings.Split(strings.Trim(tagIDsStr, "[]"), ",")

		// Debug print for tag IDs extracted from HTML
		// fmt.Println("Tag IDs from HTML for", title, ":", tagIDs)

		// Print corresponding tag names from tags.json
		// fmt.Println("Tag names for", title, ":")
		var tagNames []string
		for _, tagID := range tagIDs {
			tagName, ok := tags[tagID]
			if ok {
				tagNames = append(tagNames, tagName)
				// fmt.Println(tagName)
			} else {
				// fmt.Println("Tag name not found for ID:", tagID)
			}
		}

		game := Game{
			Title:       strings.TrimSpace(title),
			Link:        link,
			Price:       strings.TrimSpace(price),
			ReleaseDate: strings.TrimSpace(releaseDate),
			Reviews:     reviewsSummary,
			Tags:        tagNames, // Assign tag names directly
		}

		games = append(games, game)
	})

	return games, nil
}

func scrapeAdditionalDetails(game *Game, birthtime string) error {
	fmt.Println("Scraping additional details from:", game.Link)

	// Create HTTP client with a cookie jar
	client := &http.Client{}

	// Create a request object
	req, err := http.NewRequest("GET", game.Link, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers to bypass age verification
	req.Header.Set("Cookie", fmt.Sprintf("birthtime=%s", birthtime)) // Use the birthtime value from environment

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check if response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read HTML content
	htmlContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTML content: %v", err)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return fmt.Errorf("failed to parse HTML document: %v", err)
	}

	// Find developer
	developer := doc.Find("#developers_list > a").Text()

	// Find publisher
	publisher := doc.Find("#game_highlights > div.rightcol > div > div.glance_ctn_responsive_left > div:nth-child(4) > div.summary.column > a").Text()

	// Find description
	description := strings.TrimSpace(doc.Find(".game_description_snippet").Text())
	description = strings.ReplaceAll(description, "\n", " ") // Remove newlines

	fmt.Println("Developer:", developer)
	fmt.Println("Publisher:", publisher)
	fmt.Println("Description:", description)

	// Ask the user if they want to see system requirements
	var showSysReqs string
	fmt.Print("Do you want to see system requirements? (yes/no): ")
	fmt.Scanln(&showSysReqs)

	if strings.ToLower(showSysReqs) == "yes" {
		// Find system requirements section
		sysReqSection := doc.Find(".game_page_autocollapse.sys_req")
		if sysReqSection.Length() > 0 {
			// Extract system requirements text
			systemRequirements := strings.TrimSpace(sysReqSection.Find(".sysreq_contents").Text())
			// Clean up system requirements text
			systemRequirements = strings.ReplaceAll(systemRequirements, "\n", " ") // Remove newlines
			systemRequirements = strings.ReplaceAll(systemRequirements, "  ", "")    // Remove excessive spaces
			systemRequirements = strings.TrimSpace(systemRequirements)            // Trim leading and trailing spaces
			// Replace multiple spaces with a single space
			systemRequirements = regexp.MustCompile(`\s+`).ReplaceAllString(systemRequirements, " ")
			fmt.Println("System Requirements:", systemRequirements)
		} else {
			fmt.Println("System requirements not found.")
		}
	}

	return nil
}
