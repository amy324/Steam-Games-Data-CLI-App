package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	//"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
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
	color.Cyan("Welcome to Steam Scraper!")
	color.Cyan("----------------------------")
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

    var jsonFilePath, csvFilePath string

    reader := bufio.NewReader(os.Stdin)

    var fileOpened bool // Flag to indicate whether the file has been opened
    var games []Game    // Define games slice here to access it outside the loop

    for {
        if !fileOpened {
            // Get game keyword input from user
            var gameKeyword string
            fmt.Print("Enter the game keyword you want to search (or 'quit' to exit): ")
            gameKeyword, err = reader.ReadString('\n')
            if err != nil {
                fmt.Println("Error reading input:", err)
                return
            }
            gameKeyword = strings.TrimSpace(gameKeyword)

            if gameKeyword == "quit" {
                break
            }

            // Scrape data from the first page
            url := fmt.Sprintf("https://store.steampowered.com/search/?term=%s", gameKeyword)
            games, err = scrapePage(url, tags) // Assign to the games variable defined outside the loop
            if err != nil {
                fmt.Println("Error scraping page:", err)
                continue
            }

            // Output data to JSON file
            jsonFilePath = "resultfiles/games.json"
            jsonFile, err := os.Create(jsonFilePath)
            if err != nil {
                color.Red("Error creating JSON file:", err)
                return
            }
            defer jsonFile.Close()

            jsonEncoder := json.NewEncoder(jsonFile)
            jsonEncoder.SetIndent("", "  ")
            jsonEncoder.SetEscapeHTML(false) // Prevent escaping HTML characters
            err = jsonEncoder.Encode(games)
            if err != nil {
                color.Red("Error encoding JSON:", err)
                return
            }

            // Output data to CSV file
            csvFilePath = "resultfiles/games.csv"
            csvFile, err := os.Create(csvFilePath)
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
                color.Red("Error writing CSV header:", err)
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

            color.Green("JSON and CSV files created successfully.")

            // Prompt user to open the results
            fileOpened = openResults(jsonFilePath, csvFilePath)
        } else {
            // Prompt user to enter the link for additional details
            if err := promptAdditionalDetails(games, reader, birthtime); err != nil {
                color.Red("Error processing additional details:", err)
                continue
            }

            fileOpened = false // Reset fileOpened flag for the next iteration
        }
    }
}

func promptAdditionalDetails(games []Game, reader *bufio.Reader, birthtime string) error {
    fmt.Print("Enter the link of the game from your JSON/CSV file for additional details (or 'quit' to exit): ")
    input, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("error reading input: %v", err)
    }

    // Remove newline character from the end of the input
    input = strings.TrimSpace(input)

    if input == "quit" {
        return nil // Exit without error if user chooses to quit
    }

    // Find the game in the JSON file
    var selectedGame *Game
    for _, g := range games {
        if g.Link == input {
            selectedGame = &g
            break
        }
    }

    if selectedGame == nil {
        return fmt.Errorf("game not found")
    }

    // Scrape additional details from the selected game's link
    if err := scrapeAdditionalDetails(selectedGame, birthtime); err != nil {
        return fmt.Errorf("error scraping additional details: %v", err)
    }

    return nil
}


func openResults(jsonFilePath, csvFilePath string) bool {
	var cmd *exec.Cmd

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want to open the results? (Type 'J' to open JSON, 'C' to open CSV, or 'next' to move on): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	input = strings.ToLower(strings.TrimSpace(input))
	switch input {
	case "j":
		cmd = exec.Command("notepad", jsonFilePath) // Open JSON file with Notepad
	case "c":
		cmd = exec.Command("cmd", "/c", "start", csvFilePath) // Open CSV file with default program
	case "next":
		fmt.Println("Moving on to the next operation.")
		return true // Return true when user selects 'next'
	default:
		fmt.Println("Unsupported input.")
		return false
	}

	// Run the command to open the file
	if cmd != nil {
		if err := cmd.Start(); err != nil {
			fmt.Println("Error opening file:", err)
			return false
		}
		fmt.Println("File opened successfully.")

		// Wait for the user to explicitly type 'next' to continue
		for {
			fmt.Print("Type 'next' to move on: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return false
			}
			input = strings.ToLower(strings.TrimSpace(input))
			if input == "next" {
				fmt.Println("Moving on to the next operation.")
				return true
			}
		}
	}

	return false
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

	// Print titles with emphasis
	fmt.Println("\nDeveloper:")
	fmt.Println(developer)
	fmt.Println("\nPublisher:")
	fmt.Println(publisher)
	fmt.Println("\nDescription:")
	fmt.Println(description)

	// Ask the user if they want to see system requirements
	var showSysReqs string
	fmt.Print("\nDo you want to see system requirements? (yes/no): ")
	fmt.Scanln(&showSysReqs)

	if strings.ToLower(showSysReqs) == "yes" {
		// Find system requirements section
		sysReqSection := doc.Find(".game_page_autocollapse.sys_req")
		if sysReqSection.Length() > 0 {
			// Extract system requirements text
			systemRequirements := strings.TrimSpace(sysReqSection.Find(".sysreq_contents").Text())
			// Clean up system requirements text
			systemRequirements = strings.ReplaceAll(systemRequirements, "\n", " ") // Remove newlines
			systemRequirements = strings.ReplaceAll(systemRequirements, "  ", "")  // Remove excessive spaces
			systemRequirements = strings.TrimSpace(systemRequirements)             // Trim leading and trailing spaces
			// Replace multiple spaces with a single space
			systemRequirements = regexp.MustCompile(`\s+`).ReplaceAllString(systemRequirements, " ")
			// Print system requirements with distinct formatting
			fmt.Println("\nSystem Requirements:")
			fmt.Println(systemRequirements)
			fmt.Println() // Extra newline
		} else {
			fmt.Println("\nSystem requirements not found.")
		}
	}

	return nil
}
