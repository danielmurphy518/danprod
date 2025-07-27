package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// --- Configuration ---
// A map of short names to directory paths or web URLs.
var projects = map[string]string{
	// Local Directories
	"docs": "~/Documents",
	"dev":  "~/Developer",
	"dl":   "~/Downloads",

	// Web URLs
	"gh":   "https://github.com",
	"news": "https://news.ycombinator.com",
}

// expandPath replaces the ~ character with the user's home directory.
func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	return filepath.Join(home, path[1:]), nil
}

// openNewTab uses osascript to open a new Terminal.app tab.
func openNewTab(path string) error {
	script := fmt.Sprintf(`
	tell application "Terminal"
		if not (exists window 1) then reopen
		activate
		tell application "System Events" to keystroke "t" using command down
		delay 0.2
		do script "cd '%s' && clear" in window 1
	end tell
	`, path)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

func main() {
	// --- Argument Parsing ---
	listFlag := flag.Bool("l", false, "List all available aliases.")
	tabFlag := flag.Bool("t", false, "Open the folder in a new Terminal.app tab.")
	flag.Parse()

	alias := flag.Arg(0)

	// --- Logic ---
	if *listFlag {
		fmt.Println("Available aliases:")
		for a, p := range projects {
			fmt.Printf("  - %s: %s\n", a, p)
		}
		return
	}

	if alias == "" {
		flag.Usage()
		os.Exit(1)
	}

	path, found := projects[alias]
	if !found {
		fmt.Fprintf(os.Stderr, "Error: Alias '%s' not found.\n", alias)
		os.Exit(1)
	}

	// Check if the path is a URL or a local file path
	isURL := strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")

	if isURL {
		if *tabFlag {
			fmt.Fprintln(os.Stderr, "Error: The --tab flag can only be used with local directory paths.")
			os.Exit(1)
		}
		fmt.Printf("Opening URL for '%s': %s\n", alias, path)
		cmd := exec.Command("open", path)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening URL: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Handle as a local file path
		targetPath, err := expandPath(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Directory not found: %s\n", targetPath)
			os.Exit(1)
		}

		if *tabFlag {
			fmt.Printf("Opening '%s' in a new tab...\n", alias)
			err = openNewTab(targetPath)
		} else {
			fmt.Printf("Opening '%s' in Finder...\n", alias)
			err = exec.Command("open", targetPath).Run()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to execute open command: %v\n", err)
			os.Exit(1)
		}
	}
}


