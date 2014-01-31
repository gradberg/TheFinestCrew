/* 
The Console functions are convenience wrappers around whatever terminal or console
library the game uses. It should have little to no game-specific functionality, as
that should reside in something more akin to a GameDisplayAdapter class (which
is analgous to a ViewModel in MVVM
*/

package main

import (
    "fmt"
    termbox "termbox-go"
)

const Height = 25
const Width = 80

func Console_Init() error {
    return termbox.Init()    
}

func Console_Validate() error {
    // need to change so that it SETS the terminal
    w, h := termbox.Size()
    
    if h < Height || w < Width {
        return fmt.Errorf("Need minimum terminal size of %d x %d", Height, Width)
    } else {
        return nil
    }
}

func Console_Close() {
    termbox.Close()
}

func Console_ReadKey() termbox.Event {    
    for {
        e := termbox.PollEvent()      
        // Filter out events that are irrelevant  
        if e.Type == termbox.EventKey { return e }
    }
}

// Display cycle functions
func Console_DrawStart() {
    termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
}
func Console_DrawEnd() {
    termbox.Flush()
}

// Returns a ConsoleRange that represents the entire display area.
func Console_EntireRange() *ConsoleRange {
    return Console_NewRange(0, 0, Width - 1, Height - 1)
}

// A range represents a rectangle of cells of a given size. By
// operating on a range as a single unit, various display operations
// become much easier.
type ConsoleRange struct {
    x1, x2 int // start and end x coordinates
    y1, y2 int // start and end y coordinates
    fg termbox.Attribute
    bg termbox.Attribute
    
}
func Console_NewRange(x1, y1, x2, y2 int) *ConsoleRange {
    return &ConsoleRange { 
        x1: x1, x2: x2, 
        y1: y1, y2: y2,
        fg: termbox.ColorWhite,
        bg: termbox.ColorBlack,
    }
}
func (r *ConsoleRange) GetSize() (int,int) {
    return r.x2 - r.x1 + 1, r.y2 - r.y1 + 1
}
func (r *ConsoleRange) SetAttributes(fg, bg termbox.Attribute) {
    r.fg = fg
    r.bg = bg
    
    for x := r.x1; x <= r.x2; x++ {
        for y := r.y1; y <= r.y2; y++ {
            termbox.SetCell(x, y, ' ', fg, bg)            
        }
    }    
}
func (r *ConsoleRange) SetBorder() {
    // Corners
    termbox.SetCell(r.x1, r.y1, '┌', r.fg, r.bg)
    termbox.SetCell(r.x2, r.y1, '┐', r.fg, r.bg)
    termbox.SetCell(r.x1, r.y2, '└', r.fg, r.bg)
    termbox.SetCell(r.x2, r.y2, '┘', r.fg, r.bg)
    
    // Top and Bottom sides
    for x := r.x1 + 1; x < r.x2; x++ {
        termbox.SetCell(x, r.y1, '─', r.fg, r.bg)
        termbox.SetCell(x, r.y2, '─', r.fg, r.bg)
    }
    
    // left and right sides
    for y := r.y1 + 1; y < r.y2; y++ {
        termbox.SetCell(r.x1, y, '│', r.fg, r.bg)
        termbox.SetCell(r.x2, y, '│', r.fg, r.bg)
    }
}

func (r *ConsoleRange) SetTitle(title string) {
    // The title takes up a specific amount of space in the range
    maxLength := r.x2 - r.x1 - 3
    runes := TrimLength([]rune(title), maxLength)
    
    r.writeText(runes, r.x1 + 2, r.y1, r.fg | termbox.AttrBold, r.bg)
}

// Displays text, truncated (not wrapped) using the default colors
func (r *ConsoleRange) DisplayText(text string, x int, y int) {
    r.DisplayTextWithColor(text, x, y, r.fg, r.bg)
}

func (r *ConsoleRange) DisplayVerticalText(text string, x int, y int) {
    r.DisplayVerticalTextWithColor(text, x, y, r.fg, r.bg)
}

func (r *ConsoleRange) DisplayTextWithColor(text string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
    r.DisplayRunesWithColor([]rune(text),x,y,fg,bg)
}

func (r *ConsoleRange) DisplayRunesWithColor(text []rune, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
    // if this is outside of the Y bounds, just don't display anything
    if y < 0 || y > r.y2 - r.y1 {
        return
    }
    
    maxLength := r.x2 - x + 1
    runes := TrimLength(text, maxLength)
    r.writeText(runes, r.x1 + x, r.y1 + y, fg, bg)
}

func (r *ConsoleRange) DisplayVerticalTextWithColor(text string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
    // if this is outside of the Y bounds, just don't display anything
    if y < 0 || y > r.y2 - r.y1 {
        return
    }
    
    maxLength := r.y2 - y + 1
    runes := TrimLength([]rune(text), maxLength)
    
    for i := 0; i < len(runes); i++ {
        r.writeText(runes[i:i+1], r.x1 + x, r.y1 + y + i, fg, bg)
    }
}

// DisplayCommand (shortened to take less typing
func (r *ConsoleRange) Com(hotkey string, description string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {
    hotkeyLength := len([]rune(hotkey))
    r.DisplayTextWithColor(hotkey, x, y, fg, bg)
    r.DisplayText(description, x + hotkeyLength, y)
}

func (r *ConsoleRange) writeText(text []rune, x int, y int, fg termbox.Attribute, bg termbox.Attribute) {    
    for index, runeValue := range text {
        termbox.SetCell(x + index, y, runeValue, fg, bg)
    }
}



