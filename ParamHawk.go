package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	targetDomain string
	collectFlag  bool
	paramFlag    bool
	listFlag     bool
)

func init() {
	flag.StringVar(&targetDomain, "d", "", "Target domain")
	flag.BoolVar(&collectFlag, "c", false, "Collect all URLs from the target using waybackurls and store them")
	flag.BoolVar(&paramFlag, "p", false, "Generate param URLs")
	flag.BoolVar(&listFlag, "l", false, "Read list of domains from standard input")
	flag.Parse()
}

func displayBanner() {
	fmt.Println(`
   _____                          _    _                _    
  |  __ \                        | |  | |              | |   
  | |__) |_ _ _ __ __ _ _ __ ___ | |__| | __ ___      _| | __
  |  ___/ _' | '__/ _' | '_ ' _ \|  __  |/ _' \ \ /\ / / |/ /
  | |  | (_| | | | (_| | | | | | | |  | | (_| |\ V  V /|   < 
  |_|   \__,_|_|  \__,_|_| |_| |_|_|  |_|\__,_| \_/\_/ |_|\_\
                                                            
#################################################################################################################
# Tools Name: ParamHawk                                                                                         #
# Description: This script is designed to automate the discovery and extraction of parameters from target URLs. #
#              It removes duplicates, trims URLs, and organizes them for further analysis.                      #
# Author: 0xsaju                                                                                                #
# LinkedIn: https://linkedin.com/in/0xsaju                                                                      #
# Version: v_1.1                                                                                                #
#################################################################################################################

-d      Target domain
-c	Collect all URLs from the target using waybackurls
-p	Generate param URLs
-l	Read list of domains from standard input

`)
}

func processURL(url string) string {
	// Check if the URL contains both '?' and '='
	if strings.Contains(url, "?") && strings.Contains(url, "=") {
		// Split the URL into two parts: before and after the first '='
		parts := strings.SplitN(url, "=", 2)

		// Take only the part before the first '='
		url = parts[0] + "="

		return url
	}

	// Remove the full line if '?' and '=' are not present
	return ""
}

func getOutputFilePath(targetDomain, suffix string) string {
	// Replace invalid characters in the target domain with underscores
	targetDomain = strings.ReplaceAll(targetDomain, "://", "")
	targetDomain = strings.ReplaceAll(targetDomain, "/", "_")
	targetDomain = strings.ReplaceAll(targetDomain, ".", "_")

	return fmt.Sprintf("%s_%s", targetDomain, suffix)
}


func isWaybackURLsInstalled() bool {
	cmd := exec.Command("which", "waybackurls")
	err := cmd.Run()
	return err == nil
}

func installWaybackURLs() error {
	cmd := exec.Command("which", "waybackurls")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Installing waybackurls...")
		cmd = exec.Command("go", "install", "github.com/tomnomnom/waybackurls@latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	fmt.Println("waybackurls is already installed.")
	return nil
}

func collectURLs(targetDomain string) error {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Time taken to collect URLs: %s\n", time.Since(startTime))
	}()

	cmd := exec.Command("waybackurls", targetDomain)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running waybackurls: %v", err)
	}

	filePath := getOutputFilePath(targetDomain, "urls.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			fmt.Fprintln(file, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading waybackurls output: %v", err)
	}

	fmt.Printf("Collected URLs saved to %s\n", filePath)
	return nil
}

func generateParamURLs() error {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Time taken to generate param URLs: %s\n", time.Since(startTime))
	}()

	filePath := getOutputFilePath(targetDomain, "param_urls.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	urlsFilePath := getOutputFilePath(targetDomain, "urls.txt")
	urlsFile, err := os.Open(urlsFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer urlsFile.Close()

	// Maintain a set to track seen URLs
	seenURLs := make(map[string]bool)

	scanner := bufio.NewScanner(urlsFile)
	for scanner.Scan() {
		line := scanner.Text()
		modifiedURL := processURL(line)

		// Write modified URL to output file if not empty and not a duplicate
		if modifiedURL != "" && !seenURLs[modifiedURL] {
			fmt.Fprintln(file, modifiedURL)
			seenURLs[modifiedURL] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	fmt.Printf("Processing complete. Param URLs saved to %s\n", filePath)
	return nil
}

func readDomainsFromInput() ([]string, error) {
	var domains []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		if domain != "" {
			domains = append(domains, domain)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading domains from standard input: %v", err)
	}
	return domains, nil
}

func main() {
	displayBanner()

	if !isWaybackURLsInstalled() {
		if err := installWaybackURLs(); err != nil {
			fmt.Printf("Error installing waybackurls: %v\n", err)
			os.Exit(1)
		}
	}

	if listFlag {
		fmt.Println("Reading list of domains from standard input...")
		domains, err := readDomainsFromInput()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		for _, domain := range domains {
			targetDomain = domain

			if collectFlag {
				err := collectURLs(targetDomain)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}

			if paramFlag {
				err := generateParamURLs()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	} else {
		if targetDomain == "" {
			fmt.Println("Error: Please provide a target domain using the -d flag.")
			os.Exit(1)
		}

		if collectFlag {
			err := collectURLs(targetDomain)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}

		if paramFlag {
			err := generateParamURLs()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
