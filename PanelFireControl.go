/*
    The Fire Control panel is how the player handles weapons
*/

package main
import "termbox-go"
import "fmt"

type PanelFireControl struct { }

func (p *PanelFireControl) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    if (g.ThePlayer.FireControlSettingAim) {
        return p.processInputSettingAim(g,ch,key)
    } else if (g.ThePlayer.FireControlSettingTarget) {
        return p.processInputSettingTarget(g,ch,key)
    } else { 
        return p.processInputNormal(g,ch,key)
    }    
}

func (p *PanelFireControl) processInputSettingTarget(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }
    //w := g.PlayerShip.Weapon
    
    w := g.ThePlayer.FireControlSelectedWeapon
    handledRune := true
    switch ch {
        case 'd':
            g.ThePlayer.FireControlSettingTarget = false            
            
        case 'c':
            if (g.ThePlayer.FireControlTarget.IsValidTarget == false) { return nil }
            w.TargetShip = g.ThePlayer.FireControlTarget.DesiredTarget.(*Ship)            
            w.TargetType = TargetTypeTarget
            g.ThePlayer.FireControlSettingTarget = false
            result.TicksToPass = 1
    
        default:
            handledRune = false
    }
    if (handledRune) { return result }
    
    return g.ThePlayer.FireControlTarget.ProcessInput(g, ch, key)
}

func (p *PanelFireControl) processInputSettingAim(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }

    w := g.ThePlayer.FireControlSelectedWeapon
    switch ch {
        case 'd':
            g.ThePlayer.FireControlSettingAim = false
            
        case 'c':
            w.FiringAngle = g.ThePlayer.FireControlAngle
            w.TargetType = TargetTypeManual
            g.ThePlayer.FireControlSettingAim = false
            result.TicksToPass = 1
            
        case 'a':
            g.ThePlayer.FireControlAngleLargeCounter(w)
        case 'z':
            g.ThePlayer.FireControlAngleSmallCounter(w)        
        case 's':
            g.ThePlayer.FireControlAngleLargeClockwise(w)
        case 'x':
            g.ThePlayer.FireControlAngleSmallClockwise(w)
        
            
        default:
            return nil
    }
    
    return result
}

func (p *PanelFireControl) processInputNormal(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }

    w := g.ThePlayer.FireControlSelectedWeapon
    switch ch {
        case 'a':
            g.ThePlayer.FireControlSelectedWeaponPrevious()
        case 'z':
            g.ThePlayer.FireControlSelectedWeaponNext()
    
        case 's':
            w.AutoFire = !w.AutoFire
        
        case 'x':
            g.ThePlayer.FireControlSettingAim = true
            g.ThePlayer.FireControlAngle = w.FiringAngle
            
        case 'c':
            g.ThePlayer.FireControlSettingTarget = true
            
            
        default:
            return nil
    }
    
    return result
}

func (p *PanelFireControl) Display(g *Game, r *ConsoleRange) {
    r.SetAttributes(termbox.ColorRed, termbox.ColorBlack)
    r.SetBorder()
    r.SetTitle("Fire Control")
    
    if (g.ThePlayer.FireControlSettingAim) {
        p.displaySettingAim(g,r)
    } else if (g.ThePlayer.FireControlSettingTarget) { 
        p.displaySettingTarget(g,r)
    } else { 
        p.displayNormal(g,r)
    }    
}
    
func (p *PanelFireControl) displaySettingTarget(g *Game, r *ConsoleRange) {
    black := termbox.ColorBlack
    green := termbox.ColorGreen | termbox.AttrBold
    
    r.DisplayText("Setting Target", 2, 1)
    r.Com("[d]", " Cancel", 2, 2, black, green)
 
    if (g.ThePlayer.FireControlTarget.IsValidTarget) {
        r.Com("[c]", " Set", 2, 9, black, yellow)
    }
    
 
    g.ThePlayer.FireControlTarget.Display(r, black, green)
}

func (p *PanelFireControl) displaySettingAim(g *Game, r *ConsoleRange) {
    black := termbox.ColorBlack
    green := termbox.ColorGreen | termbox.AttrBold
    yellow := termbox.ColorYellow | termbox.AttrBold

    r.DisplayText("Setting Aim", 2, 1)
    r.Com("[d]", " Cancel", 2, 2, black, green)
    r.Com("[c]", " Set", 2, 9, black, yellow)
        
    w := g.ThePlayer.FireControlSelectedWeapon
    
    p.writeNameWithReload(g,r,w,19,1)
    r.DisplayText(w.DesignName, 19, 2)
    // ---- 
    
    r.DisplayText(fmt.Sprintf("%03.f   Firing Arc   %03.f", w.FiringArcStart, w.FiringArcEnd), 8, 3)
    r.Com("[a][z]","",6,4, black, green)
    r.Com("[x][s]","",26,4, black, green)
    
    // Display relative to the ship
    r.DisplayText(fmt.Sprintf("%05.1f (ship relative)", g.ThePlayer.FireControlAngle), 3, 5)        
    compassX := 4
    compassY := 7
    compassFg := termbox.ColorRed
    compassBg := termbox.ColorWhite
    
    var firingArcPoints []CompassPoint
    firingArcPoints = Compass_GetArcPoints(w.FiringArcStart, w.FiringArcEnd)
    for i := range firingArcPoints {
        point := firingArcPoints[i]
        r.DisplayTextWithColor(" ", compassX + point.X, compassY + point.Y, compassFg, compassBg)
    }
    

    shipHeadingIcon := Compass_GetShipHeadingIcon(g.ThePlayer.FireControlAngle)
    r.DisplayTextWithColor(shipHeadingIcon, compassX, compassY, compassFg, compassBg)
    
    icon, x, y := Compass_GetNearDirectionArrow(g.ThePlayer.FireControlAngle)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, compassFg, compassBg)
    
    // Display in absolute terms too, as that may helpful when the ship is spinning various directions
    absoluteAngle := AddAngles(g.ThePlayer.FireControlAngle, g.PlayerShip.ShipHeadingInDegrees)
    compassX = 32
    compassY = 6
    
    adjustedStart := AddAngles(w.FiringArcStart, g.PlayerShip.ShipHeadingInDegrees)
    adjustedEnd := AddAngles(w.FiringArcEnd, g.PlayerShip.ShipHeadingInDegrees)    
    firingArcPoints = Compass_GetArcPoints(adjustedStart, adjustedEnd)
    for i := range firingArcPoints {
        point := firingArcPoints[i]
        r.DisplayTextWithColor(" ", compassX + point.X, compassY + point.Y, compassFg, compassBg)
    }
    
    shipHeadingIcon = Compass_GetShipHeadingIcon(absoluteAngle)
    r.DisplayTextWithColor(shipHeadingIcon, compassX, compassY, compassFg, compassBg)
    
    icon, x, y = Compass_GetNearDirectionArrow(absoluteAngle)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, compassFg, compassBg)
   
    r.DisplayText(fmt.Sprintf("(absolute) %05.1f", absoluteAngle), 20, 8)
}


func (p *PanelFireControl) displayNormal(g *Game, r *ConsoleRange) {
    brightRed := termbox.ColorRed | termbox.AttrBold
    black := termbox.ColorBlack
    green := termbox.ColorGreen | termbox.AttrBold
    white := termbox.ColorWhite | termbox.AttrBold
        
    r.Com("[a]", "", 2, 1, black, green )
    r.Com("[z]", "", 2, 9, black, green )
    
    // ---- at some point, this will break, and will need to be changed
    //      to SCROLL.
    
    // Display the list of weapons
    index := 0
    for we := g.PlayerShip.Weapons.Front(); we != nil; we = we.Next() {
        w := we.Value.(*ShipWeapon)
        
        if (w == g.ThePlayer.FireControlSelectedWeapon) {
            r.DisplayTextWithColor(">", 1, 2 + index, brightRed, black)
        }
        p.writeNameWithReload(g,r,w,2,2 + index)
        
        index++
    }
    
    // Display the selected weapon
    w := g.ThePlayer.FireControlSelectedWeapon
    r.DisplayText(w.DesignName, 19, 2)
    r.Com("[s]"," Auto Fire",19,3,black, green)
    if (w.AutoFire) {
        r.DisplayTextWithColor(" ON ", 33,3, black, white)
    } else {    
        r.DisplayText("OFF ", 33,3)
    }   
    
    var ast string    
    ast = " "; if w.TargetType == TargetTypeManual { ast = "*" }
    r.DisplayTextWithColor(ast, 18, 4, brightRed, black)
    r.Com("[x]", " Manual Aim", 19,4, black, green)
    r.DisplayTextWithColor(ast, 33, 4, brightRed, black)
    if (w.TargetType == TargetTypeManual) {
        // display the selected fire angle
        absoluteAngle := AddAngles(w.FiringAngle, g.PlayerShip.ShipHeadingInDegrees)
        r.DisplayText(fmt.Sprintf("R:%05.1f  A:%05.1f", w.FiringAngle, absoluteAngle), 20,5)
    }
   
    
    ast = " "; if w.TargetType == TargetTypeTarget { ast = "*" }
    r.DisplayTextWithColor(ast, 18, 6, brightRed, black)
    r.Com("[c]", " Set Target" + ast, 19, 6, black, green)
    r.DisplayTextWithColor(ast, 33, 6, brightRed, black)
    if (w.TargetType == TargetTypeTarget && w.TargetShip != nil) { 
        r.DisplayText(w.TargetShip.GetName(), 20, 7)
    }
   
/*   
    ast = " "; if w.TargetType == TargetTypeFireAtWill { ast = "*" }
    r.DisplayText(ast, 18, 8)
    r.Com("[d]", " Fire At Will" + ast, 19, 8, black, yellow)
  */  
    
}

func (p *PanelFireControl) writeNameWithReload(g *Game, r *ConsoleRange, w *ShipWeapon, x,y int) {
    darkGrey := termbox.ColorBlack | termbox.AttrBold

    fg := termbox.ColorRed
    if (w.AutoFire) {
        fg = termbox.ColorRed | termbox.AttrBold
    }

    maxLength := 15
    denominator := w.DesignCycle
    if (denominator == 0) { denominator = 1 }
    
    percentage := float64(denominator - w.CurrentCycle) / float64(denominator)
    length := int(float64(maxLength) * percentage)        
    fullString := TrimOrPadToLength(w.EmplacementName, maxLength)
    
    firstPart := fullString[0:length]
    secondPart := fullString[length:]
    r.DisplayTextWithColor(firstPart, x, y, fg, darkGrey)    
    r.DisplayTextWithColor(secondPart, x + length, y, fg, black)  
}
