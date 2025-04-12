package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	ModeNormal  = "NORMAL"
	ModeInsert  = "INSERT"
	ModeCommand = "COMMAND"
)

type Editor struct {
	fileName      string
	buffer        []string
	cursorX       int
	cursorY       int
	offsetY       int
	mode          string
	commandBuffer string
	lastKey       rune // for handling dd command
	quit          bool
	message       string
	screen        tcell.Screen
}

func NewEditor(fileName string) (*Editor, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	return &Editor{
		fileName: fileName,
		mode:     ModeNormal,
		message:  "",
		screen:   screen,
	}, nil
}

func (e *Editor) loadFile() error {
	data, err := ioutil.ReadFile(e.fileName)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	e.buffer = lines
	return nil
}

func (e *Editor) saveFile() error {
	output := strings.Join(e.buffer, "\n")
	return ioutil.WriteFile(e.fileName, []byte(output), 0644)
}

func (e *Editor) draw() {
	e.screen.Clear()
	w, h := e.screen.Size()
	hContent := h - 1

	style := tcell.StyleDefault

	// Draw file buffer
	for i := 0; i < hContent; i++ {
		lineIndex := e.offsetY + i
		if lineIndex >= len(e.buffer) {
			break
		}
		line := e.buffer[lineIndex]
		for j, ch := range line {
			if j >= w {
				break
			}
			e.screen.SetContent(j, i, ch, nil, style)
		}
	}

	// Draw status line
	status := fmt.Sprintf(" %s | %s | Ln %d, Col %d ", e.fileName, e.mode, e.cursorY+1, e.cursorX+1)
	if e.mode == ModeCommand {
		status += " :" + e.commandBuffer
	}
	if e.message != "" {
		status += " | " + e.message
	}

	statusStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	for i := 0; i < w; i++ {
		e.screen.SetContent(i, h-1, ' ', nil, statusStyle)
	}
	for i, ch := range status {
		if i >= w {
			break
		}
		e.screen.SetContent(i, h-1, ch, nil, statusStyle)
	}

	e.screen.ShowCursor(e.cursorX, e.cursorY-e.offsetY)
	e.screen.Sync()
}

func (e *Editor) processNormalMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
	case tcell.KeyRight:
		if e.cursorX < len(e.currentLine()) {
			e.cursorX++
		}
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
			e.adjustOffset()
			if e.cursorX > len(e.currentLine()) {
				e.cursorX = len(e.currentLine())
			}
		}
	case tcell.KeyDown:
		if e.cursorY < len(e.buffer)-1 {
			e.cursorY++
			e.adjustOffset()
			if e.cursorX > len(e.currentLine()) {
				e.cursorX = len(e.currentLine())
			}
		}
	case tcell.KeyCtrlC:
		e.quit = true
	default:
		switch ev.Rune() {
		case 'h':
			if e.cursorX > 0 {
				e.cursorX--
			}
		case 'l':
			if e.cursorX < len(e.currentLine()) {
				e.cursorX++
			}
		case 'k':
			if e.cursorY > 0 {
				e.cursorY--
				e.adjustOffset()
				if e.cursorX > len(e.currentLine()) {
					e.cursorX = len(e.currentLine())
				}
			}
		case 'j':
			if e.cursorY < len(e.buffer)-1 {
				e.cursorY++
				e.adjustOffset()
				if e.cursorX > len(e.currentLine()) {
					e.cursorX = len(e.currentLine())
				}
			}
		case 'i':
			e.mode = ModeInsert
			e.message = "insert mode"
		case ':':
			e.mode = ModeCommand
			e.commandBuffer = ""
		case 'x':
			line := e.currentLine()
			if e.cursorX < len(line) && len(line) > 0 {
				e.buffer[e.cursorIndex()] = line[:e.cursorX] + line[e.cursorX+1:]
			}
		case 'd':
			if e.lastKey == 'd' {
				idx := e.cursorIndex()
				e.buffer = append(e.buffer[:idx], e.buffer[idx+1:]...)
				if e.cursorY >= len(e.buffer) && e.cursorY > 0 {
					e.cursorY--
				}
				e.lastKey = 0
			} else {
				e.lastKey = 'd'
				go func() {
					time.Sleep(300 * time.Millisecond)
					e.lastKey = 0
				}()
			}
		default:
			e.lastKey = 0
		}
	}
}

func (e *Editor) processInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		e.mode = ModeNormal
		e.message = ""
	case tcell.KeyEnter:
		line := e.currentLine()
		newLine := ""
		if e.cursorX < len(line) {
			newLine = line[e.cursorX:]
			e.buffer[e.cursorIndex()] = line[:e.cursorX]
		}
		idx := e.cursorIndex() + 1
		e.buffer = append(e.buffer[:idx], append([]string{newLine}, e.buffer[idx:]...)...)
		e.cursorY++
		e.cursorX = 0
		e.adjustOffset()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if e.cursorX > 0 {
			line := e.currentLine()
			e.buffer[e.cursorIndex()] = line[:e.cursorX-1] + line[e.cursorX:]
			e.cursorX--
		}
	default:
		if ev.Rune() != 0 {
			line := e.currentLine()
			runes := []rune(line)
			if e.cursorX > len(runes) {
				runes = append(runes, ev.Rune())
			} else {
				runes = append(runes[:e.cursorX], append([]rune{ev.Rune()}, runes[e.cursorX:]...)...)
			}
			e.buffer[e.cursorIndex()] = string(runes)
			e.cursorX++
		}
	}
}

func (e *Editor) processCommandMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		e.mode = ModeNormal
		e.commandBuffer = ""
		e.message = ""
	case tcell.KeyEnter:
		e.executeCommand()
		e.commandBuffer = ""
		if e.mode == ModeCommand {
			e.mode = ModeNormal
		}
	default:
		if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			if len(e.commandBuffer) > 0 {
				e.commandBuffer = e.commandBuffer[:len(e.commandBuffer)-1]
			}
		} else if ev.Rune() != 0 {
			e.commandBuffer += string(ev.Rune())
		}
	}
}

func (e *Editor) executeCommand() {
	cmd := strings.TrimSpace(e.commandBuffer)
	switch cmd {
	case "w":
		err := e.saveFile()
		if err != nil {
			e.message = "Error saving file: " + err.Error()
		} else {
			e.message = "File saved"
		}
	case "q":
		e.quit = true
	case "wq":
		err := e.saveFile()
		if err != nil {
			e.message = "Error saving file: " + err.Error()
		} else {
			e.quit = true
		}
	case "q!":
		e.quit = true
	default:
		e.message = "Unknown command: " + cmd
	}
}

func (e *Editor) currentLine() string {
	if e.cursorIndex() < len(e.buffer) {
		return e.buffer[e.cursorIndex()]
	}
	return ""
}

func (e *Editor) cursorIndex() int {
	return e.cursorY
}

func (e *Editor) adjustOffset() {
	_, h := e.screen.Size()
	hContent := h - 1
	if e.cursorY < e.offsetY {
		e.offsetY = e.cursorY
	} else if e.cursorY >= e.offsetY+hContent {
		e.offsetY = e.cursorY - hContent + 1
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vim-clone <file>")
		os.Exit(1)
	}

	fileName := os.Args[1]
	editor, err := NewEditor(fileName)
	if err != nil {
		fmt.Println("Error creating editor:", err)
		os.Exit(1)
	}

	err = editor.loadFile()
	if err != nil {
		fmt.Println("Error loading file:", err)
		os.Exit(1)
	}

	editor.draw()
	for {
		ev := editor.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch editor.mode {
			case ModeNormal:
				editor.processNormalMode(ev)
			case ModeInsert:
				editor.processInsertMode(ev)
			case ModeCommand:
				editor.processCommandMode(ev)
			}
			editor.draw()
			if editor.quit {
				editor.screen.Fini()
				return
			}
		case *tcell.EventResize:
			editor.screen.Sync()
		}
	}
}
