/*
    This class keeps track of all status arrangements of panels. This allows a single code file
    to be responsible for which panel is visible and when.
*/

package main

import "reflect"

type DisplayStatus uint16; const (
    DISPLAY_NONE DisplayStatus = iota  // doesn't display on screen
    DISPLAY_NORMAL   // Displays as normal on the screen    
    DISPLAY_FULL     // Displays as a full-screen panel
    DISPLAY_OVERLAY  // Displays even fuller
)

type DisplayResult struct {
    DisplayStatus DisplayStatus
    X int // Which location on the screen 'grid' this panel should be at
    Y int 
}

func (g *Game) DeterminePanelArrangement(panel interface{}) *DisplayResult {
    // ---- use individual panel bool flags for each panel to indicate that they are full-screen... stored on the Player?
    
    // First see if this is not in the playing-state, and display different
    // screens as a result.
    
    if (g.GameStatus == GameStatusDeathScreen) {
        // ---- to do
    } else if (g.GameStatus != GameStatusPlaying) {
        return overlayIfPanel(panel, &PanelTitleScreen{})
    }
    
    switch panel.(type) {   
        case *PanelOverlay:
            return newDisplayResult(DISPLAY_OVERLAY, 0, 0)
    }
    
    // These occur about the crew-role-check, because if the player does not have access to them,
    // then they should never even happen.
    if (g.ThePlayer.TacticalMapFullScreen) {
        return fullScreenIfPanel(panel, &PanelTacticalMap{})
    }
    
    // need some array map thing to store static arrange of panels? I probably don't need that currently and go with something
    // more simple.
    
    // nested switches... RUN!
    switch g.ThePlayer.CrewMember.CrewRole {
        case CrewRoleHelmsman:            
            switch panel.(type) {  
                case *PanelHelmControl:
                    return newDisplayResult(DISPLAY_NORMAL, 1, 1)
                case *PanelHelmStatus:
                    return newDisplayResult(DISPLAY_NORMAL, 0, 1)
                    
                case *PanelPersonnel:
                    return newDisplayResult(DISPLAY_NORMAL, 1, 0)
                    
                case *PanelTacticalMap:
                    return newDisplayResult(DISPLAY_NORMAL, 0, 0)
                case *PanelTacticalAnalysis:
                    return newDisplayResult(DISPLAY_NONE, 1, 0)
                case *PanelFireControl:
                    return newDisplayResult(DISPLAY_NONE, 0, 0) 
                    
                default:
                    return newDisplayResult(DISPLAY_NONE, 0, 0)
            }
            
        case CrewRoleCommander:
            switch panel.(type) {  
                case *PanelHelmControl:
                    return newDisplayResult(DISPLAY_NONE, 1, 0)
                case *PanelHelmStatus:
                    return newDisplayResult(DISPLAY_NONE, 1, 1)
                    
                case *PanelPersonnel:
                    return newDisplayResult(DISPLAY_NORMAL, 1, 0)
                    
                case *PanelTacticalMap:
                    return newDisplayResult(DISPLAY_NORMAL, 0, 0)
                case *PanelTacticalAnalysis:
                    return newDisplayResult(DISPLAY_NORMAL, 0, 1)
                case *PanelFireControl:
                    return newDisplayResult(DISPLAY_NORMAL, 1, 1) 
                    
                default:
                    return newDisplayResult(DISPLAY_NONE, 0, 0)
            }
            
        case CrewRolePilot:
            switch panel.(type) {  
                case *PanelHelmControl:
                    return newDisplayResult(DISPLAY_NORMAL, 1, 0)
                case *PanelHelmStatus:
                    return newDisplayResult(DISPLAY_NONE, 1, 1)
                    
                case *PanelPersonnel:
                    return newDisplayResult(DISPLAY_NONE, 1, 0)
                    
                case *PanelTacticalMap:
                    return newDisplayResult(DISPLAY_NORMAL, 0, 0)
                case *PanelTacticalAnalysis:
                    return newDisplayResult(DISPLAY_NORMAL, 0, 1)
                case *PanelFireControl:
                    return newDisplayResult(DISPLAY_NORMAL, 1, 1) 
                    
                default:
                    return newDisplayResult(DISPLAY_NONE, 0, 0)
            }
            
        
        default: panic("Player CrewRole undefined. Cannot run game.")
    }
    
}

func fullScreenIfPanel(testPanel interface{}, targetPanel interface{}) *DisplayResult {
    same := reflect.TypeOf(testPanel) == reflect.TypeOf(targetPanel)
    if same {
        return newDisplayResult(DISPLAY_FULL, 0, 0)
    } else {
        return newDisplayResult(DISPLAY_NONE, 0, 0)
    }
}

func overlayIfPanel(testPanel interface{}, targetPanel interface{}) *DisplayResult {
    same := reflect.TypeOf(testPanel) == reflect.TypeOf(targetPanel)
    if same {
        return newDisplayResult(DISPLAY_OVERLAY, 0, 0)
    } else {
        return newDisplayResult(DISPLAY_NONE, 0, 0)
    }
}

func newDisplayResult(status DisplayStatus, x int, y int) *DisplayResult {
    return &DisplayResult {                 
        DisplayStatus: status,
        X: x,
        Y: y,
    }
}

func (g *Game) IsThisPanelEnabled(panel interface{}) bool {
    return g.DeterminePanelArrangement(panel).DisplayStatus != DISPLAY_NONE
}
