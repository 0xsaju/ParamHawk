package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	targetDomain string
	collectFlag  bool
	paramFlag    bool
)

func init() {
	flag.StringVar(&targetDomain, "d", "", "Target domain")
	flag.BoolVar(&collectFlag, "c", false, "Collect all URLs from the target using waybackurls and store them")
	flag.BoolVar(&paramFlag, "p", false, "Generate param URLs")
	flag.Parse()
}

func displayBanner() {
	fmt.Println(`
   _____                          _    _                _    
  |  __\\                        | |  | |              | |   
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
# Version: v_1.0                                                                                                #
#################################################################################################################

-d      Target domain
-c      Collect all URLs from the target using waybackurls
-p      Generate param URLs

`)
}

func installWaybackURLs() error {
	cmd := exec.Command("go", "install", "github.com/tomnomnom/waybackurls@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func processURL(url string) string {
	// Check if the URL contains both '?' and '='
	if strings.Contains(url, "?") && strings.Contains(url, "=") {
		// Remove all characters after '='
		re := regexp.MustCompile(`=(.*?)(&|$)`)
		url = re.ReplaceAllString(url, "=")

		// Remove duplicate '=' characters
		url = strings.ReplaceAll(url, "==", "=")

		return url
	}

	// Remove the full line if '?' and '=' are not present
	return ""
}

func collectURLs(targetDomain string) error {
	// Install waybackurls as a dependency
	if err := installWaybackURLs(); err != nil {
		return fmt.Errorf("error installing waybackurls: %v", err)
	}

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

func getOutputFilePath(targetDomain, suffix string) string {
	// Replace invalid characters in the target domain with underscores
	targetDomain = strings.ReplaceAll(targetDomain, "://", "")
	targetDomain = strings.ReplaceAll(targetDomain, "/", "_")
	targetDomain = strings.ReplaceAll(targetDomain, ".", "_")

	return fmt.Sprintf("%s_%s", targetDomain, suffix)
}

func main() {
	displayBanner()

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
