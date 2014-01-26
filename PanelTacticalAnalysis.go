
package main
import "fmt"
import "termbox-go"

type PanelTacticalAnalysis struct { }

func (p *PanelTacticalAnalysis) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }

    switch ch {
        case '{':
            g.ThePlayer.TacticalAnalysisSelection = g.GetPreviousShip(g.ThePlayer.TacticalAnalysisSelection)
            if (g.ThePlayer.TacticalAnalysisSelection == g.PlayerShip) {
                g.ThePlayer.TacticalAnalysisSelection = g.GetPreviousShip(g.ThePlayer.TacticalAnalysisSelection)
            }
        case '}':
            g.ThePlayer.TacticalAnalysisSelection = g.GetNextShip(g.ThePlayer.TacticalAnalysisSelection)
            if (g.ThePlayer.TacticalAnalysisSelection == g.PlayerShip) {
                g.ThePlayer.TacticalAnalysisSelection = g.GetNextShip(g.ThePlayer.TacticalAnalysisSelection)
            }
        default:
            return nil
    }
    
    return result
}


func (p *PanelTacticalAnalysis) Display(g *Game, r *ConsoleRange) {
    w, _ := r.GetSize()
    r.SetAttributes(termbox.ColorRed, termbox.ColorBlack)
    r.SetBorder()
    r.SetTitle("Tactical Analysis")
    
    r.DisplayTextWithColor(g.PlayerShip.Name, 1,1, termbox.ColorGreen | termbox.AttrBold, termbox.ColorBlack | termbox.AttrBold)
    
    r.DisplayText("[{][}] Select Ship", 20, 0)
    r.DisplayTextWithColor("[{][}]", 20, 0, termbox.ColorRed, termbox.ColorGreen | termbox.AttrBold)    
    r.DisplayText("<--->", 17, 5)
    
    tas := g.ThePlayer.TacticalAnalysisSelection
    ps := g.PlayerShip
    
    modifierAngleInDegrees, distance := 0.0, 0.0 // if nothing selected, the player ship displays relative to 0.0
    if (tas != nil) {
        modifierAngleInDegrees, distance = tas.Point.Subtract(ps.Point).ToVector()
    }
    
    
    // Display player ship's information    
    compassX, compassY := 13, 5
    r.DisplayTextWithColor("     ", compassX - 2, compassY - 2, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY - 1, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY + 0, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY + 1, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    r.DisplayTextWithColor("     ", compassX - 2, compassY + 2, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    
    modifiedHeading := AddAngles(AddAngles(ps.ShipHeadingInDegrees, -modifierAngleInDegrees),90.0)
    shipHeadingIcon := Compass_GetShipHeadingIcon(modifiedHeading)
    r.DisplayTextWithColor(shipHeadingIcon, compassX, compassY, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    
    icon, x, y := Compass_GetNearDirectionArrow(modifiedHeading)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
    
    modifiedCourse := AddAngles(AddAngles(ps.MovementHeadingInDegrees, -modifierAngleInDegrees),90.0)
    icon, x, y = Compass_GetFarDirectionArrow(modifiedCourse)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, termbox.ColorBlack | termbox.AttrBold, termbox.ColorRed)
    
    r.DisplayText(fmt.Sprintf("HP: %5.1f", ps.HitPoints), 1, 3)
    r.DisplayTextWithColor(fmt.Sprintf("RH: %03.f", modifiedHeading), 1, 7, termbox.ColorRed, termbox.ColorBlack)
    r.DisplayTextWithColor(fmt.Sprintf("RC: %03.f", modifiedCourse), 1, 8, termbox.ColorBlack | termbox.AttrBold, termbox.ColorRed)
    r.DisplayTextWithColor(fmt.Sprintf("S: %6.1f", ps.SpeedInUnitsPerTick), 1, 9, termbox.ColorBlack | termbox.AttrBold, termbox.ColorRed)
    
    // Display the selected ship's information
    if (tas == nil) {
        r.DisplayText("-NONE-", w - 7, 1)
        
    } else {    
        nameLength := len(tas.Name)
        if (nameLength > 22) { nameLength = 22 }        
        
        // ---- color based on friend or foe
        r.DisplayTextWithColor(tas.Name, w - nameLength - 1, 1, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack | termbox.AttrBold)
        
        r.DisplayText("DIST", 18, 3)
        r.DisplayText(fmt.Sprintf("%5.f", distance), 17, 4) 
        
        bearing := AddAngles(-ps.ShipHeadingInDegrees, modifierAngleInDegrees)
        r.DisplayText("BEAR.", 17, 7)
        r.DisplayText(fmt.Sprintf("%05.1f", bearing), 17, 6)

        compassX, compassY = w - 14, 5
        r.DisplayTextWithColor("     ", compassX - 2, compassY - 2, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
        r.DisplayTextWithColor("     ", compassX - 2, compassY - 1, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
        r.DisplayTextWithColor("     ", compassX - 2, compassY + 0, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
        r.DisplayTextWithColor("     ", compassX - 2, compassY + 1, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
        r.DisplayTextWithColor("     ", compassX - 2, compassY + 2, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
  
        
    
        
        modifiedHeading := AddAngles(AddAngles(tas.ShipHeadingInDegrees, -modifierAngleInDegrees),90.0)
        shipHeadingIcon := Compass_GetShipHeadingIcon(modifiedHeading)
        r.DisplayTextWithColor(shipHeadingIcon, compassX, compassY, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
        
        icon, x, y := Compass_GetNearDirectionArrow(modifiedHeading)
        r.DisplayTextWithColor(icon, compassX + x, compassY + y, termbox.ColorRed, termbox.ColorBlack | termbox.AttrBold)
        
        modifiedCourse := AddAngles(AddAngles(tas.MovementHeadingInDegrees, -modifierAngleInDegrees), 90.0)
        icon, x, y = Compass_GetFarDirectionArrow(modifiedCourse)
        r.DisplayTextWithColor(icon, compassX + x, compassY + y, termbox.ColorBlack | termbox.AttrBold, termbox.ColorRed)
        
        dx := w - 8
        r.DisplayText(fmt.Sprintf("%5.1f :HP", tas.HitPoints), dx - 2, 3)
        r.DisplayText(fmt.Sprintf("%03.f :RH", modifiedHeading), dx, 7)
        r.DisplayTextWithColor(fmt.Sprintf("%03.f :RC", modifiedCourse), dx, 8, termbox.ColorBlack | termbox.AttrBold, termbox.ColorRed)
        r.DisplayTextWithColor(fmt.Sprintf("%6.1f :S", tas.SpeedInUnitsPerTick), dx - 2, 9, termbox.ColorBlack | termbox.AttrBold, termbox.ColorRed)
        
    }
}

