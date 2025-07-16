package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"

	"github.com/alvincrespo/glypto-go/pkg/metadata"
	"github.com/alvincrespo/glypto-go/pkg/scraper"
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape [URL]",
	Short: "Scrape metadata from a webpage",
	Long: `Scrape metadata from a webpage including Open Graph tags, Twitter Cards, 
standard meta tags, and other HTML elements.

You can provide a URL as an argument or you will be prompted to enter one.

Examples:
  glypto scrape https://example.com
  glypto scrape`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScrape,
}

func getURLFromInput(args []string) (string, error) {
	var url string

	if len(args) > 0 {
		url = args[0]
	} else {
		reader := bufio.NewReader(os.Stdin)
		color.Blue("Enter the URL to scrape metadata from: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}
		url = strings.TrimSpace(input)
	}

	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	return url, nil
}

func fetchWebpage(url string) (*http.Response, error) {
	color.Yellow("Fetching metadata from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP error! status: %d", resp.StatusCode)
	}

	return resp, nil
}

func parseHTML(resp *http.Response) (*html.Node, error) {
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	return doc, nil
}

func scrapeMetadata(doc *html.Node) (*metadata.Metadata, error) {
	scraperInstance, err := scraper.CreateScraper()
	if err != nil {
		return nil, fmt.Errorf("failed to create scraper: %w", err)
	}

	metadata, err := scraperInstance.Scrape(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape metadata: %w", err)
	}

	return metadata, nil
}

func displayResults(metadata *metadata.Metadata) {
	color.Green("\n✓ Metadata scraped successfully:\n")

	printField("Title", metadata.Title())
	printField("Description", metadata.Description())
	printField("Image", metadata.Image())
	printField("URL", metadata.URL())
	printField("Site Name", metadata.SiteName())

	favicon := metadata.Favicon()
	printField("Favicon", &favicon)

	if len(metadata.Feeds) > 0 {
		color.New(color.Bold).Println("\nFeeds:")
		for i, feed := range metadata.Feeds {
			title := "Untitled"
			if feed.Title != nil {
				title = *feed.Title
			}
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, title, feed.Type, feed.Href)
		}
	}

	printProviderData("Open Graph Tags", metadata.OpenGraph())
	printProviderData("Twitter Card Tags", metadata.TwitterCard())
}

func runScrape(cmd *cobra.Command, args []string) error {
	url, err := getURLFromInput(args)
	if err != nil {
		return err
	}

	resp, err := fetchWebpage(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := parseHTML(resp)
	if err != nil {
		return err
	}

	metadata, err := scrapeMetadata(doc)
	if err != nil {
		return err
	}

	displayResults(metadata)
	return nil
}

func printField(name string, value *string) {
	bold := color.New(color.Bold)
	if value != nil {
		bold.Printf("%s: ", name)
		fmt.Println(*value)
	} else {
		bold.Printf("%s: ", name)
		fmt.Println("Not found")
	}
}

func printProviderData(title string, data map[string][]string) {
	if len(data) > 0 {
		color.New(color.Bold).Printf("\n%s:\n", title)
		for key, values := range data {
			fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
		}
	}
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
