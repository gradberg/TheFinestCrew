/*
    Contains all of the display-specific knowledge for displaying the game on
    a given console or terminal library. 
*/

package main

func (g *Game) Display() {
    Console_DrawStart()
    
    // Determine what to display
    g.displayPanel(&PanelTitleScreen{})
    g.displayPanel(&PanelDeathScreen{})
    g.displayPanel(&PanelNextRound{})
    
    g.displayPanel(&PanelPersonnel{})
    
    g.displayPanel(&PanelTacticalMap{})
    g.displayPanel(&PanelTacticalAnalysis{})
    g.displayPanel(&PanelFireControl{})
    
    g.displayPanel(&PanelHelmStatus{})
    g.displayPanel(&PanelHelmControl{})
    
    g.displayPanel(&PanelOverlay{})
    
    Console_DrawEnd()
}

func (g *Game) displayPanel(panel IPanel) {
    //result := panel.GetDisplayStatus(g)
    result := g.DeterminePanelArrangement(panel)
    if (result.DisplayStatus == DISPLAY_NONE) { return }
    
    var r *ConsoleRange
    if (result.DisplayStatus == DISPLAY_NORMAL) {
        x := (result.X * PANEL_WIDTH) + 1
        y := (result.Y * PANEL_HEIGHT) + 1
        r = Console_NewRange(x, y, x + PANEL_WIDTH - 1, y + PANEL_HEIGHT - 1)
    } else if (result.DisplayStatus == DISPLAY_FULL) {
        r = Console_NewRange(1, 1, Width - 2, Height - 3) 
    } else if (result.DisplayStatus == DISPLAY_OVERLAY) {
        r = Console_EntireRange()
    }
    panel.Display(g, r)
}

const PANEL_WIDTH int = 39
const PANEL_HEIGHT int = 11
