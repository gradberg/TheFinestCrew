/*
    Helper methods so that multiple panels can share the same target-setting code.
*/

package main
import "fmt"
import "termbox-go"

// Holds all selection information, so that the player
// can maintain multiples of these objects at once.
type CourseSetter struct {
    g *Game // pointer to the game for each of use.
        
    Course float64 // the course the player wants to set
    Speed float64 // the desired speed
}
func NewCourseSetter(g *Game) *CourseSetter {
    return &CourseSetter { g: g}
}

// Processes input. If returns nil, the user did not press a key that this control uses.
func (cs *CourseSetter) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }
    
    switch key {                
        case termbox.KeyArrowLeft: cs.courseCounter()
        case termbox.KeyArrowRight: cs.courseClockwise()
        case termbox.KeyArrowUp: cs.speedUp()
        case termbox.KeyArrowDown: cs.speedDown()
        default: return nil
    }

    return result
}

func (cs *CourseSetter) Display(
    r *ConsoleRange, 
    //screenX, screenY int, 
    hotkeyFg, hotkeyBg termbox.Attribute,
    compassFg, compassBg termbox.Attribute,
    ) {
    
    
    //w, _ := r.GetSize()
    compassX := 25 //w / 2.0
    compassY := 5
    
    r.DisplayTextWithColor("   ", compassX - 1, compassY - 1, compassFg, compassBg)
    r.DisplayTextWithColor("   ", compassX - 1, compassY + 0, compassFg, compassBg)
    r.DisplayTextWithColor("   ", compassX - 1, compassY + 1, compassFg, compassBg)
    
    r.DisplayTextWithColor("[LEFT]", compassX - 7, compassY - 1, hotkeyFg, hotkeyBg)
    r.DisplayTextWithColor("[RIGHT]", compassX + 2, compassY - 1, hotkeyFg, hotkeyBg)
        
    // Display the desired heading
    r.DisplayText("Course", 22,2)
    r.DisplayText(fmt.Sprintf("%03.f", cs.Course), 24, 3)
    
    shipHeadingIcon := Compass_GetShipHeadingIcon(cs.Course)
    r.DisplayTextWithColor(shipHeadingIcon, compassX, compassY, compassFg, compassBg)
    
    icon, x, y := Compass_GetNearDirectionArrow(cs.Course)
    r.DisplayTextWithColor(icon, compassX + x, compassY + y, compassFg, compassBg)
    
    // Display the desired speed
    r.DisplayText(fmt.Sprintf("Speed: %4.f", cs.Speed), 3, 5)
    r.DisplayTextWithColor("[UP]", 10, 4, hotkeyFg, hotkeyBg)
    r.DisplayTextWithColor("[DOWN]", 8, 6, hotkeyFg, hotkeyBg)
    
    
}


func (cs *CourseSetter) courseClockwise() {
    cs.Course = Round(cs.Course + 1.0, 0)
    if (cs.Course >= 360.0) { cs.Course -= 360.0 }
}
func (cs *CourseSetter) courseCounter() {
    cs.Course = Round(cs.Course - 1.0, 0)
    if (cs.Course < 0.0) { cs.Course += 360.0 }
}

func (cs *CourseSetter) speedUp() {
    cs.Speed = Round(cs.Speed + 1.0, 0)
}
func (cs *CourseSetter) speedDown() {
    cs.Speed = Round(cs.Speed - 1.0, 0)
}
