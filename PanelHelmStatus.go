
package main
import "fmt"
import "termbox-go"

type PanelHelmStatus struct { }

func (p *PanelHelmStatus) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    return nil
}

func (p *PanelHelmStatus) Display(g *Game, r *ConsoleRange) {
    r.SetAttributes(termbox.ColorBlue, termbox.ColorWhite)
    r.SetBorder()
    r.SetTitle("Helm Status")
    
    ps := g.PlayerShip  
    
    // Create a white-highlighting where the heading/compass is displayed
    
    
    //w, _ := r.GetSize()
    compassX := 19 //w / 2.0
    compassY := 5
    
    r.DisplayTextWithColor("     ", compassX - 2, compassY - 2, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY - 1, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY + 0, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY + 1, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY + 2, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    
    // Display Ship's Heading
    r.DisplayText("Ship Heading", 14, 1)
    r.DisplayText(fmt.Sprintf("%03.f", ps.ShipHeadingInDegrees), 18, 2)
    
    shipHeadingIcon := Compass_GetShipHeadingIcon(ps.ShipHeadingInDegrees)
    r.DisplayTextWithColor(shipHeadingIcon, compassX, compassY, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    
    icon, x, y := Compass_GetNearDirectionArrow(ps.ShipHeadingInDegrees)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
    
    // Display the Track the ship is actually moving along
    r.DisplayTextWithColor("Ship Course", 14, 8, termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue)
    r.DisplayTextWithColor(fmt.Sprintf("%03.f", ps.MovementHeadingInDegrees), 18, 9, termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue)    
    
    icon, x, y = Compass_GetFarDirectionArrow(ps.MovementHeadingInDegrees)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue)
    
    r.DisplayTextWithColor("Speed", 4,8, termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue)
    r.DisplayTextWithColor(fmt.Sprintf("%7.1f", ps.SpeedInUnitsPerTick), 3,9, termbox.ColorWhite | termbox.AttrBold, termbox.ColorBlue)

    // Display engine status
    r.DisplayText("Engine", 30, 1)
    r.DisplayText("Status", 30, 2)
    r.DisplayText("⌂", 32, 4)
    r.DisplayText("║", 32, 5)
    r.DisplayText("╩", 32, 6)
    r.DisplayText("╞", 31, 6)
    r.DisplayText("╡", 33, 6)
    if (g.PlayerShip.Helm.UsedForwardThrusters) {
        r.DisplayTextWithColor("^^^", 31, 7, termbox.ColorYellow | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)
    }
    if (g.PlayerShip.Helm.UsedBackwardThrusters) {
        r.DisplayTextWithColor("V", 32, 3, termbox.ColorYellow | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)
    }
    if (g.PlayerShip.Helm.UsedClockwiseThrusters) {    
        r.DisplayTextWithColor(">", 31, 4, termbox.ColorYellow | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)
        r.DisplayTextWithColor("<", 34, 6, termbox.ColorYellow | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)
    }
    if (g.PlayerShip.Helm.UsedCounterThrusters) {
        r.DisplayTextWithColor("<", 33, 4, termbox.ColorYellow | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)        
        r.DisplayTextWithColor(">", 30, 6, termbox.ColorYellow | termbox.AttrBold, termbox.ColorRed | termbox.AttrBold)    
    
    }
}

