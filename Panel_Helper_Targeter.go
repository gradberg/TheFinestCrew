/*
    Helper methods so that multiple panels can share the same target-setting code.
*/

package main
import "fmt"
import "termbox-go"

type SettingTargetTypeEnum int
const (
    SettingTargetTypeShip SettingTargetTypeEnum = iota
    SettingTargetTypePlanet    
)
func (e SettingTargetTypeEnum) ToString() string {
    switch (e) {
        case SettingTargetTypeShip: return "Ship"
        case SettingTargetTypePlanet: return "Planet"
        default: return "ERROR TARGET TYPE"
    }
}

// Holds all selection information, so that the player
// can maintain multiples of these objects at once.
type Targeter struct {
    DesiredTarget ISpaceObject // current selection the player is viewing (and if okay is hit, will be used)
    TargetType SettingTargetTypeEnum // current type of 'object' being selected
    IsValidTarget bool // stupid flag to indicate if the desired target is nil or not.
    
    DisablePlanets bool
    
    g *Game // pointer to the game for each of use.
}
func NewTargeter(g *Game) *Targeter {
    return &Targeter { g: g}
}

// Processes input. If returns nil, the user did not press a key that this control uses.
func (t *Targeter) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }
 
    switch key {
        case termbox.KeyArrowLeft:
            t.desiredTargetPrevious()
            
        case termbox.KeyArrowRight:
            t.desiredTargetNext()
            
        case termbox.KeyArrowUp:
            t.targetTypeUp()
            
        case termbox.KeyArrowDown:
            t.targetTypeDown()
    
        default:
            return nil
    }

    return result
}

func (t *Targeter) Display(
    r *ConsoleRange, 
    //screenX, screenY int, 
    hotkeyFg, hotkeyBg termbox.Attribute,
    ) {
    
    // up and down changes type (ship versus planet vs ?)
    // left and right changes actual selection
    
    typeString := "ERROR UNKNOWN TYPE"
    switch t.TargetType {
        case SettingTargetTypeShip: typeString = "Ships"
        case SettingTargetTypePlanet: typeString = "Planets"
    }
    
    r.DisplayText(fmt.Sprintf("Type: %s", typeString), 24,4)
    r.Com("[UP]", "", 30, 3, hotkeyFg, hotkeyBg)
    r.Com("[DOWN]", "", 30, 5, hotkeyFg, hotkeyBg)
        
    r.Com("[LEFT][RIGHT]", " Target", 2, 5, hotkeyFg, hotkeyBg)
    
    if (t.IsValidTarget) {
        target := t.DesiredTarget                
        r.DisplayText(target.GetName(), 4, 6)
        
        _, distance := t.g.PlayerShip.Point.Subtract(target.GetPoint()).ToVector()        
        r.DisplayText(fmt.Sprintf("Distance %7.1f", distance), 4, 7)
        if target.IsDestroyed() {
            r.DisplayText("**DESTROYED**", 4, 8)
        }
    }
}


func (t *Targeter) targetTypeUp() {
    switch (t.TargetType) {
        case SettingTargetTypeShip:   
            if (t.DisablePlanets == false) {
                t.TargetType = SettingTargetTypePlanet
            }
        case SettingTargetTypePlanet:  t.TargetType = SettingTargetTypeShip
    }        
    t.DesiredTarget = nil
    t.IsValidTarget = false
    t.desiredTargetNext()
}
func (t *Targeter) targetTypeDown() {
    switch (t.TargetType) {
        case SettingTargetTypeShip:           
            if (t.DisablePlanets == false) {
                t.TargetType = SettingTargetTypePlanet
            }
        case SettingTargetTypePlanet: t.TargetType = SettingTargetTypeShip
    }        
    t.DesiredTarget = nil
    t.IsValidTarget = false
    t.desiredTargetNext()
}

func (t *Targeter) desiredTargetNext() {
    switch (t.TargetType) {
        case SettingTargetTypeShip:   
            var ship *Ship
            if (t.IsValidTarget) {
                ship = t.DesiredTarget.(*Ship)
            }
            
            newShip := t.g.GetNextShip(ship)
            if (t.g.PlayerShip == newShip) {
                newShip = t.g.GetNextShip(newShip)
                t.DesiredTarget = newShip
                t.IsValidTarget = (newShip != nil)
            } else {
                t.DesiredTarget = newShip
                t.IsValidTarget = (newShip != nil)
            }
        case SettingTargetTypePlanet: 
            var planet *Planet
            if (t.IsValidTarget) {
                planet = t.DesiredTarget.(*Planet)
            }
            
            newPlanet := t.g.GetNextPlanet(planet)
            t.DesiredTarget = newPlanet
            t.IsValidTarget = (newPlanet != nil)
    }   
}
func (t *Targeter) desiredTargetPrevious() {
    switch (t.TargetType) {
        case SettingTargetTypeShip:   
            var ship *Ship
            if (t.IsValidTarget) {
                ship = t.DesiredTarget.(*Ship)
            }
            
            newShip := t.g.GetPreviousShip(ship)
            if (t.g.PlayerShip == newShip) {
                newShip = t.g.GetPreviousShip(newShip)
                t.DesiredTarget = newShip
                t.IsValidTarget = (newShip != nil)
            } else {
                t.DesiredTarget = newShip
                t.IsValidTarget = (newShip != nil)
            }
        case SettingTargetTypePlanet: 
            var planet *Planet
            if (t.IsValidTarget) {
                planet = t.DesiredTarget.(*Planet)
            }
            
            newPlanet := t.g.GetPreviousPlanet(planet)
            t.DesiredTarget = newPlanet
            t.IsValidTarget = (newPlanet != nil)
    }   
}