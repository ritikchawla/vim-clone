# A Vim Clone in Go

A lightweight text editor inspired by Vim, implemented in Go using the tcell library for terminal handling.

## Features

### Core Features (MVP)
- 📝 **Open and Display Files**
  - Load and render file contents
  - Scroll through files larger than terminal height

- 🧭 **Cursor Movement**
  - Arrow keys for navigation
  - Vim-style h, j, k, l movement keys
  - Viewport scrolling to keep cursor on screen

- ⌨️ **Multiple Modes**
  - Normal Mode (default)
    - Navigation
    - Delete characters
    - Switch to Insert mode
  - Insert Mode
    - Direct text input
    - Live editing
    - Esc to return to Normal mode
  - Command Mode
    - Accessed via ':'
    - File operations

- 💾 **File Operations**
  - `:w` - Save current file
  - `:q` - Quit
  - `:wq` - Save and quit
  - `:q!` - Force quit without saving

- ❌ **Basic Editing**
  - `i` - Enter Insert mode
  - `x` - Delete character under cursor
  - `dd` - Delete current line

- 🧠 **Status Line**
  - Current mode indicator
  - File name
  - Cursor position
  - Command input display
  - Status messages

## Installation

1. Ensure you have Go installed on your system
2. Clone this repository:
   ```bash
   git clone https://github.com/ritikchawla/vim-clone.git
   cd vim-clone
   ```
3. Build the project:
   ```bash
   go build
   ```

## Usage

Run Gim with a file name:
```bash
./vim-clone <filename>
```

### Basic Commands
- `i` - Enter Insert mode
- `Esc` - Return to Normal mode
- `:w` - Save file
- `:q` - Quit
- `:wq` - Save and quit
- `:q!` - Force quit without saving
- Arrow keys or `h,j,k,l` - Move cursor
- `x` - Delete character under cursor
- `dd` - Delete current line

## Dependencies
- [tcell](https://github.com/gdamore/tcell) - Terminal handling library