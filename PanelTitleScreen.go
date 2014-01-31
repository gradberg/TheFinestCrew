/*
    The Overlay Panel is the one the display any additional inforjmation above and below
    the panels
*/

package main
//import "fmt"
import "termbox-go"

type PanelTitleScreen struct { }

func (p *PanelTitleScreen) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }
    
    switch ch {
        // Quit
        case 's':
            g.doSetupForDevelopment()
            g.GameStatus = GameStatusPlaying            
            
        case 'q': 
            result.Exit = true
            
        default:
            return nil
    }
    
    return result
}

func (p *PanelTitleScreen) Display(g *Game, r *ConsoleRange) {
    title := []string {
        // 0
        "                                                                                ",
        "  Here's to                                                                     ",
        "                                                                                ",
        "    TTTTT  H   H  EEEEE          FFFFF  IIIII  N   N  EEEEE   SSSS  TTTTT       ",        
        "      T    H   H  E              F        I    NN  N  E      S        T         ",
        "      T    HHHHH  EEE            FFF      I    N N N  EEE     SSS     T         ",
        "      T    H   H  E              F        I    N  NN  E          S    T         ",
        "      T    H   H  EEEEE          F      IIIII  N   N  EEEEE  SSSS     T         ",
        "                                                                                ",
        "                                                                                ",
        "            CCCC  RRRR   EEEEE  W   W                      --▄▄                 ",
        "           C      R   R  E      W W W              \\▄    ▄██████     █          ",
        "           C      RRRR   EEE    W W W           ▄█████████████████████          ",
        "           C      R  R   E       W W           ██▀   /▀     ▀\\    ████          ",
        "            CCCC  R   R  EEEEE   W W                                            ",
        "                                                                                ",
        "                         in the fleet.                                          ",
        "                                                                                ",
        "                                                                                ",
        "                                                                                ",
        // 20
    }        
    for i, s := range title {
        r.DisplayText(s, 0, i)
    }
    r.DisplayTextWithColor("▄█", 48,12, termbox.ColorBlue | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayVerticalTextWithColor("███", 70, 11, termbox.ColorRed, termbox.ColorBlack)
    r.DisplayTextWithColor("██", 71, 11, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayTextWithColor("███", 71, 12, termbox.ColorYellow | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayTextWithColor("██", 71, 13, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack)
    
    r.Com("[s]"," Start Game", 5, 19, termbox.ColorWhite | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)
    r.Com("[q]"," Quit", 5, 20, termbox.ColorBlack, termbox.ColorWhite | termbox.AttrBold) 
    
    
    
    
}

 