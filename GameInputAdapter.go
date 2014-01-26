/*
    Contains all of the input-specific knowledge for the user interacting with the game on
    a given console or terminal library. 
*/

package main

import "fmt"
import termbox "termbox-go"

// Primary function responsible for accepting user input, using that to change the state of the game?
func (g *Game) DoInput() *InputResult {
    air := AsyncInput_ReceiveInput() 
    
    // if this was a timer event, then handle that special
    if air.IsTimerEvent {
        if g.ThePlayer.RunningRealTime {
            // Queue up another timer event 
            g.queueTimerEvent()
            return &InputResult { TicksToPass: 1 }
        } else {
            // if realtime is disabled now, then ignore any timer events
            return &InputResult { }    
        }
    }
    
    // unfortunately the async input it peculair enough that turning it on should run in this method
    if g.ThePlayer.RunningRealTime == false && air.Key == termbox.KeyEnter {
        g.ThePlayer.RunningRealTime = true
        g.queueTimerEvent()
        return &InputResult { }
    }
                
    
    var result *InputResult
    
    // Should go first, to catch any "Stop real time" button presses.
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelOverlay{})
    
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelPersonnel{})
    
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelTacticalMap{})
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelTacticalAnalysis{})
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelFireControl{})
    
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelHelmStatus{})
    result = g.processInputForPanel(result, air.Ch, air.Key, &PanelHelmControl{})
    
        
    if result == nil {    
        // Read normal keys, then try keycodes    
        if (air.Ch == 0) {
            g.message = fmt.Sprintf("Unrecognized keycode: %d", uint(air.Key))
        } else {        
            // Message about unrecognized command? 
            g.message = fmt.Sprintf("Unrecognized key: %c", air.Ch)
        }    
        result = &InputResult { }
    }
    
    return result
}

func (g *Game) processInputForPanel(previousResult *InputResult, ch rune, key termbox.Key, panel IPanel) *InputResult {
    // As these are chained together, if the previous call did have a non-nil result, then
    // just return that as the first panel to accept input wins.
    if previousResult != nil { return previousResult }
    
    // if this panel is disabled, completely skip it
    if g.IsThisPanelEnabled(panel) == false { return previousResult } 
    
    return panel.ProcessInput(g, ch, key)
}

func (g *Game) queueTimerEvent() {
    AsyncInput_EmptyTimerQueue()
    AsyncInput_QueueTimer(int64(1000.0 / g.ThePlayer.RealTimeTicksPerSecond))
}
