package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	appName    = "Sidenote"
	fileName   = "sidenote.txt"
	folderName = "Sidenote"

	windowWidth  = 600
	windowHeight = 500

	saveDelay = 500 * time.Millisecond
)

var (
	savePath  string
	textArea  *widget.Entry
	saveTimer *time.Timer
)

func main() {
	helpFlag := flag.Bool("h", false, "Show help")
	helpFlagLong := flag.Bool("help", false, "Show help")
	cleanFlag := flag.Bool("c", false, "Clean - backup current file and create new one")
	cleanFlagLong := flag.Bool("clean", false, "Clean - backup current file and create new one")
	flag.Parse()

	if *helpFlag || *helpFlagLong {
		printHelp()
		return
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to determine config directory: %v\n", err)
		os.Exit(1)
	}

	savePath = filepath.Join(configDir, folderName, fileName)

	if err := os.MkdirAll(filepath.Dir(savePath), 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create config directory: %v\n", err)
		os.Exit(1)
	}

	if *cleanFlag || *cleanFlagLong {
		backupCurrentFile()
		fmt.Println("Clean operation completed. Old file backed up.")
		return
	}

	runApp()
}

func runApp() {
	myApp := app.NewWithID("com.sidenote.app")
	window := myApp.NewWindow(appName)

	window.Resize(fyne.NewSize(windowWidth, windowHeight))
	window.SetPadded(true)
	window.SetIcon(theme.DocumentIcon())

	textArea = widget.NewMultiLineEntry()
	textArea.Wrapping = fyne.TextWrapWord
	textArea.SetPlaceHolder("Type your notes here...")

	loadContent()

	textArea.OnChanged = func(_ string) {
		debouncedSave()
	}

	scrollContainer := container.NewScroll(textArea)
	content := container.NewBorder(nil, nil, nil, nil, scrollContainer)
	window.SetContent(content)

	window.SetCloseIntercept(func() {
		saveContent()
		window.Close()
	})

	window.ShowAndRun()
}

func loadContent() {
	content, err := os.ReadFile(savePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		}
		return
	}
	textArea.SetText(string(content))
}

func debouncedSave() {
	if saveTimer != nil {
		saveTimer.Stop()
	}
	saveTimer = time.AfterFunc(saveDelay, saveContent)
}

func saveContent() {
	text := textArea.Text
	dir := filepath.Dir(savePath)
	tmpPath := filepath.Join(dir, "."+fileName+".tmp")

	if err := os.WriteFile(tmpPath, []byte(text), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Save failed: %v\n", err)
		return
	}

	if err := os.Rename(tmpPath, savePath); err != nil {
		fmt.Fprintf(os.Stderr, "Replace failed: %v\n", err)
	}
}

func backupCurrentFile() {
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		fmt.Println("No existing file to backup.")
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := savePath + "." + timestamp + ".bkp"

	if err := os.Rename(savePath, backupPath); err != nil {
		fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
		return
	}

	if err := os.WriteFile(savePath, []byte(""), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create new file: %v\n", err)
	}
}

func printHelp() {
	fmt.Println("Sidenote - Minimalist note taking app")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sidenote                    Launch the app")
	fmt.Println("  sidenote -h, --help         Show this help")
	fmt.Println("  sidenote -c, --clean        Backup current notes and create new file")
	fmt.Println()
	fmt.Println("Notes are saved to: ~/.config/Sidenote/sidenote.txt")
}
