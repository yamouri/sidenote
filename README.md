# Sidenote
![Platform](https://img.shields.io/badge/platform-linux-blue)
![Type](https://img.shields.io/badge/type-desktop--app-lightgrey)
![Language](https://img.shields.io/github/languages/top/yamouri/sidenote)
![Release](https://img.shields.io/github/v/release/yamouri/sidenote)
![License](https://img.shields.io/github/license/yamouri/sidenote)

A minimalist, auto-saving note-taking app for Linux.

![Image.png](image.png)

## Features

- **Zero buttons, zero menus** - Just type
- **Auto-save on every keystroke** - Wow

## Installation

### From Source
```bash
git clone https://github.com/yamouri/sidenote.git
cd sidenote
go mod tidy
go build -o sidenote
sudo cp sidenote /usr/local/bin/
```

## Dependencies

- Go 1.21+
- Fyne 

## Usage

```bash
# Launch the app
sidenote

# Show help
sidenote --help
sidenote -h

# Backup current notes and create new file
sidenote --clean
sidenote -c

# More info
- Your notes are automatically saved to ~/.config/Sidenote/sidenote.txt
```

