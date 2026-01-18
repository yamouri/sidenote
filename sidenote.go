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
)

var (
	savePath string
	textArea *widget.Entry
)

func main() {
	// Parse command line arguments
	helpFlag := flag.Bool("h", false, "Show help")
	helpFlagLong := flag.Bool("help", false, "Show help")
	cleanFlag := flag.Bool("c", false, "Clean - backup current file and create new one")
	cleanFlagLong := flag.Bool("clean", false, "Clean - backup current file and create new one")
	flag.Parse()

	if *helpFlag || *helpFlagLong {
		printHelp()
		return
	}

	configDir, _ := os.UserConfigDir()
	savePath = filepath.Join(configDir, folderName, fileName)
	os.MkdirAll(filepath.Dir(savePath), 0755)

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

	window.Resize(fyne.NewSize(600, 500))
	window.SetPadded(true)

	textArea = widget.NewMultiLineEntry()
	textArea.Wrapping = fyne.TextWrapWord
	textArea.TextStyle = fyne.TextStyle{Monospace: false}
	textArea.SetPlaceHolder("Type your notes here...")

	loadContent()

	textArea.OnChanged = func(text string) {
		saveContent()
	}

	scrollContainer := container.NewScroll(textArea)
	scrollContainer.SetMinSize(fyne.NewSize(580, 480))

	content := container.NewBorder(nil, nil, nil, nil, scrollContainer)
	window.SetContent(content)
	window.SetIcon(theme.DocumentIcon())

	window.SetCloseIntercept(func() {
		saveContent()
		window.Close()
	})

	window.ShowAndRun()
}

func loadContent() {
	content, err := os.ReadFile(savePath)
	if err == nil {
		textArea.SetText(string(content))
	}
}

func saveContent() {
	text := textArea.Text
	err := os.WriteFile(savePath, []byte(text), 0644)
	if err != nil {
		fmt.Printf("Save error: %v\n", err)
	}
}

func backupCurrentFile() {
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		fmt.Println("No existing file to backup.")
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := savePath + "." + timestamp + ".bkp"

	err := os.Rename(savePath, backupPath)
	if err != nil {
		fmt.Printf("Backup failed: %v\n", err)
		return
	}

	// Create new empty file
	os.WriteFile(savePath, []byte(""), 0644)
}

func printHelp() {
	fmt.Println("Sidenote - Minimalist note taking app")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  sidenote                    Launch the app")
	fmt.Println("  sidenote -h, --help         Show this help")
	fmt.Println("  sidenote -c, --clean        Backup current notes and create new file")
	fmt.Println()
	fmt.Println("Notes are auto-saved to: ~/.config/Sidenote/sidenote.txt")
}
