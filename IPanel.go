/*
    This interface defines the operations that all command panels should have in order to plug-and-play into
    the display and input systems.
*/

package main

import "termbox-go"

// Contains any information to pass between accepting input and processing the turn
type InputResult struct {
//    Handled bool     // indicates this panel did handle the input
    Exit bool        // indicates the user has chosen to exit
    TicksToPass uint // indicates how many ticks should pass given user input.
}

type IPanel interface {       
    // Non-nil return value indicates it handled the input
    ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult
    
    Display(g *Game, r *ConsoleRange)
    
}
