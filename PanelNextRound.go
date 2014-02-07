/*
    Displays useful information between rounds.
*/

package main
import "fmt"
import "termbox-go"

type PanelNextRound struct { }

func (p *PanelNextRound) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }

    switch ch {
        // Quit
        case '1':
            g.setupForNextRound()
            g.GameStatus = GameStatusPlaying
        
        case 'Q', 'q': 
            result.Exit = true
            
        default:
            return nil
    }
    
    return result
}

func (p *PanelNextRound) Display(g *Game, r *ConsoleRange) {
    white := termbox.ColorWhite | termbox.AttrBold
    grey := termbox.ColorWhite
    black := termbox.ColorBlack
    
    fd := NewFlowDocument(78, 22)    
    fd.AddParagraph(fmt.Sprintf("You have destroyed %d ships.", g.kills), grey, black)
    fd.AddParagraph(" ", grey, black)
    fd.AddParagraph(fmt.Sprintf("Starting round %d", g.round+1), grey, black)
    fd.Write(r, 1,1)
    
    r.Com("[1]", " Continue", 1,22, black, white)
    r.Com("[Q]", " Quit", 1, 23, black, white)
    
    
}



 