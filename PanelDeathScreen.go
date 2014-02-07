/*
    Displays information when the player is defeated.
*/

package main
import "fmt"
import "termbox-go"

type PanelDeathScreen struct { }

func (p *PanelDeathScreen) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }

    switch ch {
        // Quit
        case 'Q', 'q': 
            result.Exit = true
            
        default:
            return nil
    }
    
    return result
}

func (p *PanelDeathScreen) Display(g *Game, r *ConsoleRange) {
    white := termbox.ColorWhite | termbox.AttrBold
    grey := termbox.ColorWhite
    black := termbox.ColorBlack
    
    fd := NewFlowDocument(78, 23)
    fd.AddParagraph("You have been destroyed", grey, black)
    fd.AddParagraph(" ", grey, black)
    fd.AddParagraph(fmt.Sprintf("You destroyed %d ships before heading off to Spacey Jones Locker.", g.kills), grey, black)
    fd.Write(r, 1,1)
    
    r.Com("[Q]", " Quit", 1, 23, black, white)
    
    
}



 