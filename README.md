# Steam Scraper CLI APP

Steam Scraper is a command-line tool written in Go for scraping game data from the Steam store. This project utilizes web scraping, JSON and CSV file handling, concurrent HTTP requests, and command-line interface (CLI) implementation using Cobra.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Usage Walkthrough](#usage-walkthrough)
- [Code Explanation](#code-explanation)
- [Compatibility](#compatibility)
- [Legal](#legal)
- [Contributing](#contributing)
- [License](#license)

## Overview

### Steam:

[Steam](https://store.steampowered.com/) is a digital distribution platform developed by Valve Corporation, offering a vast array of video games, software, and other digital content to millions of users worldwide. It serves as a centralized hub for purchasing, downloading, and managing digital entertainment across multiple platforms.

### Purpose of the Application:

The purpose of this application is to provide users with a streamlined way to access and explore information about games available on the Steam platform. By leveraging web scraping techniques, the application retrieves valuable data such as game titles, prices, release dates, reviews, tags, developers, and publishers directly from the Steam store website.

## Features

- **Game Data Retrieval**: The application scrapes data from the Steam store website, including game titles, prices, release dates, reviews, and tags.
- **User Interaction**: Users can interact with the application through a command-line interface (CLI), providing keywords to search for specific games and selecting actions such as opening JSON or CSV files.
- **Detailed Information**: Users can access detailed information about individual games, including developer, publisher, description, and system requirements.
- **Customizable Output:** The application allows users to export retrieved data in both JSON and CSV formats for further analysis or integration with other tools.
- **Additional Details**: Prompt the user to input a game link for scraping additional details such as developer, publisher, description, and system requirements.
- **Compliance**: The application ensures compliance with the Steam Data Terms of Use, respecting user privacy and data handling practices outlined by Valve Corporation.

## Installation

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/amy234/steam-scraper.git
    ```

2. **Navigate to the Project Directory**:

    ```bash
    cd steam-scraper
    ```

3. **Add .env File**:

   Create an `.env` file and set your Date of Birth in as a [Unix Timestamp](https://www.unixtimestamp.com/), for example, if you were born on the 1st January 2000, the Unix Timestamp for your Date of Birth is 946684800. This environmental variable will enable the scraper to work for age restricted content.

    ```env
    BIRTHTIME=<DOB in Unix Timestamp format>
    ```

4. **Build the Executable**:

    ```bash
    go build
    ```

5. **Run the Executable**:

    ```bash
    ./steam-scraper
    ```


## Usage

1. **Search and Scrape**:
    - Enter the game keyword you want to search when prompted.
    - The tool will scrape relevant game data from the Steam store.
    - Results will be saved in the `resultfiles` directory as `games.json` and `games.csv`.

2. **Additional Details**:
    - After scraping, you can input a game link from the JSON/CSV file for additional details.
    - The tool will scrape and display developer, publisher, description, and system requirements.

3. **Open Results**:
    - Choose to open the JSON or CSV file with your preferred text editor or spreadsheet application.
    - Alternatively, type 'next' to move on to the next operation.

## Usage Walkthrough

Below is a step-by-step guide on how to use the program after compiling it with `go build` and opening the executable file via `cmd.exe`.

```cmd
Welcome to Steam Scraper!
----------------------------
Enter the game keyword you want to search (or 'quit' to exit): doom
JSON and CSV files created successfully.
Do you want to open the results? (Type 'J' to open JSON, 'C' to open CSV, or 'next' to move on): J
File opened successfully.
Type 'next' to move on:
```

Both JSON and CSV files were created and saved as `resultfiles\games.csv` and `resultfiles\games.json`.

In the example above, I chose to open the JSON file by typing `J`, which contains the results for games retrieved using the keyword `doom`. Here's a snippet of the first result, and you can explore both the full CSV and JSON files in the `resultfiles` folder of this repository:

```json
{
  "title": "DOOM",
  "link": "https://store.steampowered.com/app/379720/DOOM/?snr=1_7_7_151_150_1",
  "price": "£15.99",
  "release_date": "12 May, 2016",
  "reviews": "Overwhelmingly Positive",
  "tags": [
    "FPS",
    "Gore",
    "Action",
    "Shooter",
    "Demons",
    "Great Soundtrack",
    "First-Person"
  ]
}
```

After typing `next` to proceed, the program prompted for a game link for additional details. Using the link from the previous example, the program fetched additional details including developer, publisher, and description:

```cmd
Type 'next' to move on: next
Moving on to the next operation.
Enter the link of the game from your JSON/CSV file for additional details (or 'quit' to exit): https://store.steampowered.com/app/379720/DOOM/?snr=1_7_7_151_150_1
Scraping additional details from: https://store.steampowered.com/app/379720/DOOM/?snr=1_7_7_151_150_1

Developer:
id Software

Publisher:
Bethesda Softworks

Description:
Now includes all three premium DLC packs (Unto the Evil, Hell Followed, and Bloodfall), maps, modes, and weapons, as well as all feature updates including Arcade Mode, Photo Mode, and the latest Update 6.66, which brings further multiplayer improvements as well as revamps multiplayer progression.

Do you want to see system requirements? (yes/no): yes

System Requirements:
Minimum:
OS *: Windows 7/8.1/10 (64-bit versions)
Processor: Intel Core i5-2400/AMD FX-8320 or better
Memory: 8 GB RAM
Graphics: NVIDIA GTX 670 2GB/AMD Radeon HD 7870 2GB or better
Storage: 55 GB available space
Additional Notes: Requires Steam activation and broadband internet connection for Multiplayer and SnapMap

Recommended:
OS *: Windows 7/8.1/10 (64-bit versions)
Processor: Intel Core i7-3770/AMD FX

-8350 or better
Memory: 8 GB RAM
Graphics: NVIDIA GTX 970 4GB/AMD Radeon R9 290 4GB or better
Storage: 55 GB available space
Additional Notes: Requires Steam activation and broadband internet connection for Multiplayer and SnapMap
* Starting January 1st, 2024, the Steam Client will only support Windows 10 and later versions.

Enter the game keyword you want to search (or 'quit' to exit):
```



## Code Explanation

The codebase is organized into several files and functions:


1. **Main Function (`main`):**
   - Displays a welcome message.
   - Defines the root command for the CLI tool using the Cobra library.
   - Executes the root command.

2. **Scraping Function (`scrape`):**
   - Loads tags data from a JSON file that maps tag IDs to tag names.
   - Creates a folder to store extracted files.
   - Uses a loop to repeatedly prompt the user for a game keyword until they choose to quit.
   - Scrapes data from the Steam store based on the user's keyword input.
   - Outputs the scraped data to JSON and CSV files.
   - Prompts the user to open the results and continues to the next operation.
   - Utilizes the `birthtime` environmental variable to bypass age verification on the Steam website. This variable should be set to the user's birthdate and is sent as a cookie in the HTTP request headers to access age-restricted content, eg: the game used in the walkthrough above, *Doom*, requires age verification to search for due to violent content in the game. 

3. **Additional Details Prompt Function (`promptAdditionalDetails`):**
   - Prompts the user to input a link for additional details of a specific game.
   - Searches for the selected game in the JSON file.
   - Calls the `scrapeAdditionalDetails` function to fetch and display additional details.

4. **Results Opening Function (`openResults`):**
   - Prompts the user to open the JSON or CSV files containing the scraped data.
   - Opens the selected file using the appropriate system command.
   - Waits for the user to type 'next' to continue to the next operation.

5. **Page Scraping Function (`scrapePage`):**
   - Makes an HTTP request to the Steam store search page.
   - Parses the HTML response using **Goquery** to extract game information such as title, link, price, release date, reviews, and tags.
   - Returns a slice of `Game` structs containing the scraped data.

6. **Additional Details Scraping Function (`scrapeAdditionalDetails`):**
   - Makes an HTTP request to the provided game link.
   - Uses Goquery to parse the HTML response and extract additional details such as developer, publisher, description, and system requirements.
   - Displays the extracted details and prompts the user to view system requirements if available.

## Techniques Used

- **Web Scraping**: Utilizes the `goquery` library to extract structured data from HTML documents retrieved via HTTP requests.
- **Concurrency**: Makes concurrent HTTP requests to scrape multiple game pages simultaneously, improving performance.
- **JSON and CSV Handling**: Saves the scraped data into JSON and CSV files for easy storage and further analysis.
- **Command-Line Interface (CLI)**: Implements a CLI using the Cobra library, providing a user-friendly interface for interacting with the tool.

### HTML Scraping Process:

1. **HTTP Request:**
   - The program starts by making an HTTP GET request to the Steam store search page using the `http.Get` method provided by Go's standard library.

2. **Response Handling:**
   - Upon receiving the response, the program checks if the status code is OK (200) to ensure a successful request.
   
3. **HTML Parsing:**
   - Goquery library is used to parse the HTML content of the response body. Goquery allows easy traversal and manipulation of HTML documents using CSS selectors.
   
4. **Scraping Game Information:**
   - The program extracts game information such as title, link, price, release date, reviews, and tag IDs from specific HTML elements on the Steam store search page.
   - Each game's information is encapsulated in a `Game` struct.
   
5. **Matching Tags:**
   - For each game scraped, the program retrieves the associated tag IDs from the HTML attributes.
   - It then matches these tag IDs with their corresponding tag names stored in a JSON file (`tags.json`) located in the `data` directory.
   
6. **Tag Inclusion:**
   - Tags play a crucial role in categorizing games. The program ensures that the tags associated with each game are included in the results.
   - By referencing the `tags.json` file, the program maps tag IDs to their respective tag names.
   - The tag names are then added to the `Game` struct as a slice of strings.

### How Tag Matching Works:

- **Tags JSON File (`tags.json`):**
  - This file contains a mapping of tag IDs to tag names. Each entry in the JSON file represents a tag ID with its corresponding name.
  - Information relating to tag IDs from https://store.steampowered.com/tagdata/populartags/english
  - For example:
   ```json
    {
    
    "4400": "Abstract",
    "19": "Action",
    "4231": "Action RPG",
    "4106": "Action-Adventure",
    "21": "Adventure",
    "22602": "Agriculture",
    "1673": "Aliens",
    }
    ```
  - In the example game used in the walkthrough above (*DOOM 2016*), the HTML on https://store.steampowered.com/search/?term=doom for this partciular game includes :
    
    ```html
    data-ds-tagids="[1663,4345,19,1774,9541,1756,3839]"
    ```
    These tag IDs correspond to   "FPS", "Gore", "Action", "Shooter", "Demons",  "Great Soundtrack", and "First-Person" in the tags.json file.
- **Matching Process:**
  - When scraping game data, the program extracts tag IDs associated with each game from the HTML.
  - It then searches for these tag IDs in the `tags.json` file to find their corresponding names.
  - If a match is found, the tag name is added to the `Game` struct. Otherwise, the tag is skipped.
  
- **Inclusion in Results:**
  - Once all relevant information, including tags, is gathered for a game, it is included in the final results.
  - The tags provide additional context and categorization for each game, enhancing the usefulness of the scraped data.

## Compatibility

- **Operating System**: This tool is primarily developed and tested on Windows environments.For best results, consider running the .exe file you compile from cmd.exe. While the core functionality should work on other platforms, the file opening mechanism may be Windows-specific. Compatibility with macOS and Linux is not guaranteed.
- **Dependencies**: Requires Go version 1.16 or higher and the `goquery` library. All dependencies are managed via Go Modules.

## Legal 

### Compliance with Steam Data Terms of Use:

This program ensures compliance with the [Steam Data Terms of Use](https://steamcommunity.com/dev/apiterms) set forth by Valve Corporation, despite not directly utilizing the Steam Web API. The following measures are implemented to adhere to these terms:

1. **Privacy Policy:**
   - A privacy policy is provided, detailing the handling of nonpublic end user data, including Steam Data obtained through web scraping. The program treats Steam Data consistently with this policy as it does not handle any nonpublic end user data - it only handles public data relating to games available for purchase on the Steam website.

2. **User Consent:**
   - This program does not collect or store any user information. 

3. **Data Handling:**
   - Steam Data is stored in countries identified in the privacy policy. The program treats Steam Data with disclaimers substantially equivalent to those specified in Sections 7 and 8 of the Steam Data Terms of Use. It provides Steam Data "as is", with it purpose only being the retrieval of public games information.

4. **Endorsement and Affiliation:**
   - The program does not present Steam Data in a manner that suggests endorsement or affiliation with Valve or Steam.

5. **Compliance with Steam Subscriber Agreement:**
   - The program ensures compliance with the Steam Subscriber Agreement, avoiding actions that may violate the agreement or degrade the operation of Steam or any games distributed via Steam.

6. **Fair Play and Competitive Advantage:**
   - The program does not utilize Steam Data in any way that may provide users with an unfair competitive advantage in multiplayer versions of Steam games.

7. **Unsolicited Marketing Communications:**
   - Steam Data obtained through web scraping is not used for unsolicited marketing communications. It does not retrieve or store any user-related data; It only retrieves public information relating to games.

8. **Breach Reporting and Enforcement:**
   - In the event of any breach of the Steam Data Terms of Use or Steam Subscriber Agreement by users of the program, prompt reporting and corrective action are taken to remedy such breaches.

By implementing these measures, the program respects the rights and privacy of Steam users and ensures compliance with the Steam Data Terms of Use, despite not directly interfacing with the Steam Web API.


### Disclaimer:
This application was created by a video game enthusiast and is not affiliated with Steam or Valve Corporation. It is designed to provide access to publicly available information about video games from the Steam platform. However, please note that this application does not facilitate the purchase or play of games directly. For purchasing and playing games, please visit the official [Steam website](https://store.steampowered.com/) or other authorized games retailers. The data provided by this application is for informational purposes only and should not be considered as an endorsement or promotion of any specific game or product.


## Contributing

Contributions are welcome! If you find any bugs or want to add new features, please submit an issue or open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

