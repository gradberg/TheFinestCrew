
package main
import "fmt"
import "termbox-go"

type PanelHelmControl struct { }

func (p *PanelHelmControl) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    if g.PlayerShip.Helm.IsDirectPilot {
        return p.processInputForDirectPilot(g, ch, key)
    } else if g.ThePlayer.HelmControlSettingCourse {
        return p.processInputForSettingCourse(g, ch, key)
    } else if g.ThePlayer.HelmControlSettingTarget {
        return p.processInputForSettingTarget(g, ch, key)
    } else {
        return p.processInputForAutoPilot(g, ch, key)
    }
}

func (p *PanelHelmControl) processInputForSettingTarget(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }
    
    handledRune := true
    switch ch {
        case 'm':
            g.ThePlayer.HelmControlSettingTarget = false
        case '\\':            
            g.ThePlayer.HelmControlSettingTarget = false
            g.PlayerShip.Helm.AutoPilotMode = AutoPilotModeTarget
            g.PlayerShip.Helm.AutoPilotDesiredTarget = g.ThePlayer.HelmControlTarget.DesiredTarget
            result.TicksToPass = 1 
        default:
            handledRune = false
    }
    if (handledRune) { return result }
    
    return g.ThePlayer.HelmControlTarget.ProcessInput(g, ch, key)
/*
    switch key {
    
        default:
            return nil
    }

    return result
    */
}

func (p *PanelHelmControl) processInputForSettingCourse(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }

    handledRune := true
    switch ch {
        case 'm':
            g.ThePlayer.HelmControlSettingCourse = false
        case '\\':            
            g.ThePlayer.HelmControlSettingCourse = false
            g.PlayerShip.Helm.AutoPilotMode = AutoPilotModeCourse
            g.PlayerShip.Helm.AutoPilotDesiredCourse = g.ThePlayer.HelmControlCourseSetter.Course
            g.PlayerShip.Helm.AutoPilotDesiredSpeed = g.ThePlayer.HelmControlCourseSetter.Speed
            result.TicksToPass = 1
        default:
            handledRune = false
    }
    if (handledRune) { return result }
    
    return g.ThePlayer.HelmControlCourseSetter.ProcessInput(g, ch, key)
/*
    switch key {                
        case termbox.KeyArrowLeft:
            g.ThePlayer.HelmControlDesiredCourseCounter()
            
        case termbox.KeyArrowRight:
            g.ThePlayer.HelmControlDesiredCourseClockwise()
            
        case termbox.KeyArrowUp:
            g.ThePlayer.HelmControlDesiredSpeedUp()
            
        case termbox.KeyArrowDown:
            g.ThePlayer.HelmControlDesiredSpeedDown()
    
        default:
            return nil
    }

    return result
    */
}

func (p *PanelHelmControl) processInputForAutoPilot(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }

    handledRune := true
    switch ch {
        case 'm':
            g.PlayerShip.Helm.IsDirectPilot = true
        case ',':
            g.PlayerShip.Helm.AutoPilotMode = AutoPilotModeOff
            result.TicksToPass = 1
        case '.':
            g.PlayerShip.Helm.AutoPilotMode = AutoPilotModeFullStop
            result.TicksToPass = 1            
        case '/':
            g.ThePlayer.HelmControlSettingCourse = true
        case '\'':
            g.ThePlayer.HelmControlSettingTarget = true
        default:
            handledRune = false
    }
    if (handledRune) { return result }

    switch key {
        case termbox.KeyHome:
            g.PlayerShip.Helm.AutoPilotPowerUp()
        case termbox.KeyEnd:
            g.PlayerShip.Helm.AutoPilotPowerDown()
        default:
            return nil
    }

    return result
}

func (p *PanelHelmControl) processInputForDirectPilot(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }

    handledRune := true
    switch ch {
        case '/':
            g.PlayerShip.Helm.AutoThrust = ! g.PlayerShip.Helm.AutoThrust
            g.PlayerShip.Helm.ClearAutoPilot()
        case ',':
            g.PlayerShip.Helm.ThrustersDown()
            g.PlayerShip.Helm.ClearAutoPilot()
        case '.':
            g.PlayerShip.Helm.ThrustersUp()
            g.PlayerShip.Helm.ClearAutoPilot()
        case 'm':
            g.PlayerShip.Helm.IsDirectPilot = false
        default:
            handledRune = false
    }
    if (handledRune) { return result }

    switch key {
        case termbox.KeyArrowUp:
            if g.PlayerShip.Helm.AutoThrust == false {
                g.PlayerShip.Helm.PilotIntent = PilotIntentThrust
                result.TicksToPass = 1
            }
            g.PlayerShip.Helm.ClearAutoPilot()

        case termbox.KeyHome:
            g.PlayerShip.Helm.ThrottleUp()
            g.PlayerShip.Helm.ClearAutoPilot()
        case termbox.KeyEnd:
            g.PlayerShip.Helm.ThrottleDown()
            g.PlayerShip.Helm.ClearAutoPilot()
            
        case termbox.KeyArrowLeft:            
            g.PlayerShip.Helm.PilotIntent = PilotIntentSpinCounter
            result.TicksToPass = 1
            g.PlayerShip.Helm.ClearAutoPilot()
            
        case termbox.KeyArrowRight:
            g.PlayerShip.Helm.PilotIntent = PilotIntentSpinClockwise
            result.TicksToPass = 1
            g.PlayerShip.Helm.ClearAutoPilot()
            
        default:
            return nil
    }

    return result
}


func (p *PanelHelmControl) Display(g *Game, r *ConsoleRange) {
    r.SetAttributes(termbox.ColorBlue, termbox.ColorWhite)
    r.SetBorder()
    r.SetTitle("Helm Control") 
    
    if g.PlayerShip.Helm.IsDirectPilot {
        p.displayForDirectPilot(g, r)
    } else if g.ThePlayer.HelmControlSettingCourse {
        p.displayForSettingCourse(g, r)
    } else if g.ThePlayer.HelmControlSettingTarget {
        p.displayForSettingTarget(g, r)
    } else {
        p.displayForAutoPilot(g, r)
    }
}

func (p *PanelHelmControl) displayForSettingTarget(g *Game, r *ConsoleRange) {
    r.DisplayText("Setting Target", 2, 1)
    r.DisplayText("[m] Cancel", 2, 2)
    r.DisplayTextWithColor("[m]", 2, 2, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    
    g.ThePlayer.HelmControlTarget.Display(r,termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    
    if (g.ThePlayer.HelmControlTarget.IsValidTarget) {
        r.DisplayText("[\\] Engage", 2, 9)
        r.DisplayTextWithColor("[\\]", 2, 9, termbox.ColorBlue, termbox.ColorYellow | termbox.AttrBold)
    }
}

func (p *PanelHelmControl) displayForSettingCourse(g *Game, r *ConsoleRange) {
    r.DisplayText("Setting Course", 2, 1)
    r.DisplayText("[m] Cancel", 2, 2)
    r.DisplayTextWithColor("[m]", 2, 2, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    r.DisplayText("[\\] Engage", 2, 9)
    r.DisplayTextWithColor("[\\]", 2, 9, termbox.ColorBlue, termbox.ColorYellow | termbox.AttrBold)
    
    g.ThePlayer.HelmControlCourseSetter.Display(
        r, 
        termbox.ColorBlue , termbox.ColorGreen | termbox.AttrBold, 
        termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold,
    )        
}

func (p *PanelHelmControl) displayForAutoPilot(g *Game, r *ConsoleRange) {
    w, h := r.GetSize()
    
    r.DisplayText("Auto Piloting", 2, 1)
    r.DisplayText("[m] Switch to Direct Piloting", 2, 2)
    r.DisplayTextWithColor("[m]", 2, 2, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)

    // ---- Display asterix around the setting that is currently set
    var ast string
    
    ast = " "; if g.PlayerShip.Helm.AutoPilotMode == AutoPilotModeOff { ast = "*" }
    r.DisplayText(fmt.Sprintf("%s[,] Off%s", ast, ast), 3, 3)
    r.DisplayTextWithColor("[,]", 4, 3, termbox.ColorBlue, termbox.ColorYellow | termbox.AttrBold)
    
    ast = " "; if g.PlayerShip.Helm.AutoPilotMode == AutoPilotModeFullStop { ast = "*" }
    r.DisplayText(fmt.Sprintf("%s[.] Full Stop%s", ast, ast), 3, 4)
    r.DisplayTextWithColor("[.]", 4, 4, termbox.ColorBlue, termbox.ColorYellow | termbox.AttrBold)
    
    ast = " "; if g.PlayerShip.Helm.AutoPilotMode == AutoPilotModeCourse { ast = "*" }
    r.DisplayText(fmt.Sprintf("%s[/] Set Course%s", ast, ast), 3, 5)
    r.DisplayTextWithColor("[/]", 4, 5, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
        
    ast = " "; if g.PlayerShip.Helm.AutoPilotMode == AutoPilotModeTarget { ast = "*" }
    r.DisplayText(fmt.Sprintf("%s['] Set Target%s", ast, ast), 3, 6)
    r.DisplayTextWithColor("[']", 4, 6, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    
    //r.DisplayText(fmt.Sprintf("%s[]] Plan Maneuver%s", ast, ast), 3, 7)
    //r.DisplayText("[]] Plan Maneuver", 3, 7)
        
    //r.DisplayText(fmt.Sprintf("%s[\\] Execute Planned Maneuver%s", ast, ast), 3, 8)
    //6r.DisplayText("[\\] Execute Planned Maneuver", 3, 8)
    
    // maneuever creation, and target selection, go to full screen?
    
    
    r.DisplayVerticalText("1 2 4 8", w -2, 3)
    r.DisplayVerticalText("/ / / /", w -3, 3)
    r.DisplayVerticalText("1 1 1 1", w -4, 3)
    r.DisplayText("POWER", w - 6, 1)
    r.DisplayText(">", w-5, 3 + (3 -g.PlayerShip.Helm.AutoPilotPower) * 2)
    
    r.DisplayTextWithColor("[HOME]", w - 7, 0, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    r.DisplayTextWithColor("[END]", w - 6, h - 1, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
}

func (p *PanelHelmControl) displayForDirectPilot(g *Game, r *ConsoleRange) {
    r.DisplayText("Direct Piloting", 2, 1)
    r.DisplayText("[m] Switch to Auto Piloting", 2, 2)
    r.DisplayTextWithColor("[m]", 2, 2, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    
    x := 2
    
    r.DisplayText("[/] Auto-Thrust: ", x, 3)    
    r.DisplayTextWithColor("[/]", x, 3, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    onOff := "OFF"
    if g.PlayerShip.Helm.AutoThrust { onOff = "ON " }
    r.DisplayText(onOff, x + 17, 3)
    
    r.DisplayText(         " Spin        Spin      ", x, 5)    
    if g.PlayerShip.Helm.AutoThrust == false {
        r.DisplayText(         "      Thrust              ", x, 4)
        r.DisplayTextWithColor(       "[UP]", x + 7, 5, termbox.ColorBlue, termbox.ColorYellow | termbox.AttrBold)
    }
    r.DisplayTextWithColor(   "[LEFT][RIGHT]", x + 3, 6, termbox.ColorBlue, termbox.ColorYellow | termbox.AttrBold)
    
    // Display the throttle
    w, h := r.GetSize()
    throttleMid := ((h - 2) * 2 / 3)
    for i := 1; i <= h - 2; i++ {
        var icon string
        switch i {
            case 1: icon = "╗"
            case h - 2: icon = "╝"
            case throttleMid: icon = "╣"
            default: icon = "║"
        }
        r.DisplayText(icon, w - 3, i)
    }
    
    t := g.PlayerShip.Helm.ThrottlePercentage
    var throttleLocation int
    if t == 0 {
        throttleLocation = throttleMid
    } else if t == 100 {
        throttleLocation = 1
    } else if t == -100 {
        throttleLocation = h - 2
    } else if t > 0 {
        a := float64(throttleMid - 2) * float64(g.PlayerShip.Helm.ThrottlePercentage) / 100.0
        throttleLocation = throttleMid - 1 - int(Round64(a))        
    } else if g.PlayerShip.Helm.ThrottlePercentage < 0.0 {
        a := float64(h - 3 - throttleMid) * float64(g.PlayerShip.Helm.ThrottlePercentage) / -100.0
        throttleLocation = throttleMid + 1 + int(Round64(a))                
    }
    r.DisplayText(">", w - 4, throttleLocation)
    r.DisplayText(formatPercentage(g.PlayerShip.Helm.ThrottlePercentage), w - 8, throttleLocation)
    
    r.DisplayVerticalText("THROTTLE", w-2, 1)
    r.DisplayTextWithColor("[HOME]", w - 7, 0, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
    r.DisplayTextWithColor("[END]", w - 6, h - 1, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
        
    st := 7
    r.DisplayText("╚", st,  h-2)
    for i := st + 1; i < st + 10; i ++ {
        r.DisplayText("═", i,  h-2)
    }
    r.DisplayText("╝", st + 10,  h-2)
    displacement := int(g.PlayerShip.Helm.ThrustersPercentage / 10.0)    
    r.DisplayText(fmt.Sprintf("%03dv", g.PlayerShip.Helm.ThrustersPercentage), st + displacement - 3, h-3)    
    r.DisplayText("Rotational", st + 12, h -4)
    r.DisplayText("Thrusters",st + 13 , h- 3)
    r.DisplayTextWithColor("[,][.]", st + 16, h-2, termbox.ColorBlue, termbox.ColorGreen | termbox.AttrBold)
}

func formatPercentage(value int) string {
    if (value == 100) {
        return " 100"
    } else if (value == -100) {
        return "-100"
    } else if (value == 0) {
        return " 000"
    } else if (value > 0) {
        return fmt.Sprintf(" %03d", value)
    } else {
        return fmt.Sprintf("-%03d", value * -1)
    }
}

