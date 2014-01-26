/*
    The Overlay Panel is the one the display any additional inforjmation above and below
    the panels
*/

package main
import "fmt"
import "termbox-go"

type PanelOverlay struct { }

func (p *PanelOverlay) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }
    
    // The reatime up down can be changed whether or not this is in realtime mode
    if key == termbox.KeyPgup {
        g.ThePlayer.IncreaseRealTimeTicksPerSecond()
        return &InputResult {}
    } else if key == termbox.KeyPgdn {
        g.ThePlayer.DecreaseRealTimeTicksPerSecond()
        return &InputResult {}
    }
    
    // If realtime is enabled, any other inputs disable it
    if g.ThePlayer.RunningRealTime {
        // Any other key turns of realtie
        g.ThePlayer.RunningRealTime = false
        return &InputResult { }
    }
    
    if (ch == 0) {
        switch key {          
            case termbox.KeySpace:
                // User has chosen to pass a tick
                result.TicksToPass = 1
                            
            default:
                return nil
        }
        
    } else {        
        switch ch {
            // Quit
            case 'Q': 
                result.Exit = true
                
            default:
                return nil
        }
    }    
    
    return result
}

func (p *PanelOverlay) Display(g *Game, r *ConsoleRange) {
    r.DisplayText(fmt.Sprintf(
        "%s %s, %s aboard %s, a %s.", 
        g.ThePlayer.CrewMember.FirstName,
        g.ThePlayer.CrewMember.LastName,
        g.ThePlayer.CrewMember.CrewRole.ToString(),
        g.PlayerShip.Name, 
        g.PlayerShip.DesignName),
    1, 0)
    r.DisplayText(fmt.Sprintf("Tick: %6d", g.tick), 67, 0)
    
    r.DisplayText("[Space] Pass Time", 1, 23)
    r.DisplayTextWithColor("[Space]", 1, 23, termbox.ColorBlack , termbox.ColorYellow | termbox.AttrBold)
    
    r.DisplayText("[Enter] RealTime: ", 24,23)
    if (g.ThePlayer.RunningRealTime) {
        r.DisplayTextWithColor(" ON ", 24+18,23, termbox.ColorBlack, termbox.ColorWhite | termbox.AttrBold)    
    } else {
        r.DisplayText("OFF", 24+18,23)        
    }
    r.DisplayTextWithColor(fmt.Sprintf("[Enter]"), 24,23, termbox.ColorWhite | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)
    
    r.DisplayText(fmt.Sprintf("[PgUp][PgDn] %4.1f Ticks/Second", g.ThePlayer.RealTimeTicksPerSecond), 48,23)
    r.DisplayTextWithColor(fmt.Sprintf("[PgUp][PgDn]"), 48,23, termbox.ColorWhite | termbox.AttrBold, termbox.ColorGreen | termbox.AttrBold)
    
    r.Com("[Q]", " Quit", 70, 24, termbox.ColorBlack, termbox.ColorWhite | termbox.AttrBold) 
    
    r.DisplayText(g.message, 1, 24)
    
    
    
}

