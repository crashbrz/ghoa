package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// User represents the structure of the GitHub user data returned by the API
type User struct {
	Login   string `json:"login"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Bio     string `json:"bio"`
	Company string `json:"company"`
	Type    string `json:"type"`
	TwoFA   string `json:"two_factor_authentication"`
}

// Repository represents the structure of GitHub repository data
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

// Colors for output
const (
	green = "\033[32m"
	red   = "\033[31m"
	reset = "\033[0m"
)

var (
	removeColor bool
)

// printWithColor prints text with optional color
func printWithColor(color, text string) {
	if removeColor {
		fmt.Println(text)
	} else {
		fmt.Println(color + text + reset)
	}
}

// validateAndRetrieveToken checks if a GitHub token is valid, retrieves user info, scopes, and optionally private repositories
func validateAndRetrieveToken(token, endpoint string, retrieveInfo, retrieveRepos bool) (*User, []string, []Repository, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, nil, nil, false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, nil, nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, nil, false
	}

	var user User
	if retrieveInfo {
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			fmt.Printf("Error decoding user response: %v\n", err)
			return nil, nil, nil, false
		}
	}

	// Retrieve scopes from the X-OAuth-Scopes header
	scopesHeader := resp.Header.Get("X-OAuth-Scopes")
	scopes := strings.Split(scopesHeader, ",")
	for i := range scopes {
		scopes[i] = strings.TrimSpace(scopes[i])
	}

	// Retrieve private repositories if requested
	var repositories []Repository
	if retrieveRepos {
		reposEndpoint := strings.Replace(endpoint, "/user", "/user/repos?visibility=private", 1)
		reposReq, err := http.NewRequest("GET", reposEndpoint, nil)
		if err != nil {
			fmt.Printf("Error creating repository request: %v\n", err)
			return &user, scopes, nil, true
		}
		reposReq.Header.Set("Authorization", "Bearer "+token)

		reposResp, err := client.Do(reposReq)
		if err != nil {
			fmt.Printf("Error retrieving repositories: %v\n", err)
			return &user, scopes, nil, true
		}
		defer reposResp.Body.Close()

		if reposResp.StatusCode == http.StatusOK {
			if err := json.NewDecoder(reposResp.Body).Decode(&repositories); err != nil {
				fmt.Printf("Error decoding repository response: %v\n", err)
				return &user, scopes, nil, true
			}
		}
	}

	return &user, scopes, repositories, true
}

// processTokensWithConcurrency validates tokens using multiple goroutines
func processTokensWithConcurrency(tokens []string, endpoint string, goroutines int, showInfo, retrieveRepos, showInvalid bool) {
	var wg sync.WaitGroup
	tokensChan := make(chan string, len(tokens))

	for _, token := range tokens {
		tokensChan <- token
	}
	close(tokensChan)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for token := range tokensChan {
				user, scopes, repos, valid := validateAndRetrieveToken(token, endpoint, showInfo, retrieveRepos)
				if valid {
					printWithColor(green, fmt.Sprintf("Valid token: %s", token))
					if showInfo && user != nil {
						fmt.Printf("User Info: %+v\n", *user)
						fmt.Printf("Scopes: %v\n", scopes)
					}
					if retrieveRepos {
						fmt.Println("Private Repositories:")
						for _, repo := range repos {
							fmt.Printf("- %s (%s)\n", repo.Name, repo.URL)
						}
					}
				} else {
					if showInvalid {
						printWithColor(red, fmt.Sprintf("Invalid token: %s", token))
					}
				}
			}
		}()
	}

	wg.Wait()
}

func main() {
	// Command-line flags
	keyFlag := flag.String("k", "", "GitHub OAuth2 token to validate")
	fileFlag := flag.String("f", "", "File containing tokens (one per line)")
	numGoroutines := flag.Int("t", 1, "Number of goroutines to use")
	showInvalid := flag.Bool("d", false, "Show invalid tokens")
	removeColorFlag := flag.Bool("remove-color", false, "Remove color from output")
	showInfo := flag.Bool("i", false, "Retrieve and display detailed information and scopes about valid tokens")
	retrieveRepos := flag.Bool("p", false, "Retrieve and display private repositories for valid tokens")
	endpointFlag := flag.String("e", "https://api.github.com/user", "GitHub API endpoint to use for validation")

	flag.Parse()

	removeColor = *removeColorFlag

	if *keyFlag != "" {
		// Validate a single token
		user, scopes, repos, valid := validateAndRetrieveToken(*keyFlag, *endpointFlag, *showInfo, *retrieveRepos)
		if valid {
			printWithColor(green, fmt.Sprintf("Valid token: %s", *keyFlag))
			if *showInfo && user != nil {
				fmt.Printf("User Info: %+v\n", *user)
				fmt.Printf("Scopes: %v\n", scopes)
			}
			if *retrieveRepos {
				fmt.Println("Private Repositories:")
				for _, repo := range repos {
					fmt.Printf("- %s (%s)\n", repo.Name, repo.URL)
				}
			}
		} else {
			if *showInvalid {
				printWithColor(red, fmt.Sprintf("Invalid token: %s", *keyFlag))
			}
		}
	} else if *fileFlag != "" {
		// Validate tokens from a file
		file, err := os.Open(*fileFlag)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var tokens []string
		for scanner.Scan() {
			tokens = append(tokens, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		processTokensWithConcurrency(tokens, *endpointFlag, *numGoroutines, *showInfo, *retrieveRepos, *showInvalid)
	} else {
		fmt.Println("Please specify a token (-k) or a file (-f) containing tokens.")
		flag.Usage()
	}
}
