
package main
import "termbox-go"
import "fmt"

// enum with all possible statuses
type PersonnelStatusEnum int; const (
    PersonnelStatusNormal PersonnelStatusEnum = iota
    PersonnelStatusFullScreen
    PersonnelStatusSayRoot
    PersonnelStatusSayHelmsman
    PersonnelStatusSayHelmsmanCourse
    PersonnelStatusSayHelmsmanDestination
)


type PanelPersonnel struct { }

func (p *PanelPersonnel) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    switch (g.ThePlayer.PersonnelStatus) {
        case PersonnelStatusNormal: return p.processInputNormal(g, ch, key)
        default:         
            if (g.ThePlayer.CrewMember.CrewRole == CrewRoleHelmsman) {
                return p.processInputHelmsmanSay(g, ch, key)
            } else if (g.ThePlayer.CrewMember.CrewRole == CrewRoleCommander) {
                return p.processInputCommanderSay(g, ch, key)
            }  else {
                return nil
            }
    }
}

func (p *PanelPersonnel) processInputHelmsmanSay(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }
    
    ps := g.ThePlayer.PersonnelStatus
    switch ch {
        case '1': ps = PersonnelStatusNormal
        default: return nil
    }
    g.ThePlayer.PersonnelStatus = ps
    
    return result
}

func (p *PanelPersonnel) processInputCommanderSay(g *Game, ch rune, key termbox.Key) *InputResult {
    // Handle any special controls first
    if (g.ThePlayer.PersonnelStatus == PersonnelStatusSayHelmsmanDestination) {
        result := g.ThePlayer.PersonnelHelmsmanTarget.ProcessInput(g, ch, key)
        if (result != nil) { return result }
    }   
    if (g.ThePlayer.PersonnelStatus == PersonnelStatusSayHelmsmanCourse) {
        result := g.ThePlayer.PersonnelHelmsmanCourseSetter.ProcessInput(g, ch, key)
        if (result != nil) { return result }
    }
    
    result := &InputResult { Exit: false, TicksToPass: 0 }
    
    ps := g.ThePlayer.PersonnelStatus
    switch ch {
        case '1': 
            switch (g.ThePlayer.PersonnelStatus) {
                case PersonnelStatusSayHelmsman: ps = PersonnelStatusSayRoot
                case PersonnelStatusSayHelmsmanDestination: ps = PersonnelStatusSayHelmsman
                default: ps = PersonnelStatusNormal
            }            
            
        case '2':
            switch (g.ThePlayer.PersonnelStatus) {
                case PersonnelStatusSayRoot: ps = PersonnelStatusSayHelmsman      
                case PersonnelStatusSayHelmsman: // Full Stop
                    helmsman := g.PlayerShip.GetCrewMemberForRole(CrewRoleHelmsman)
                    g.EnqueueMessage(NewMessageFullStop(g.ThePlayer.CrewMember, helmsman))
                    ps = PersonnelStatusNormal
                    result.TicksToPass = 1 
                
                default: return nil
            }
                  
            
        case '3':
            switch (g.ThePlayer.PersonnelStatus) {
                case PersonnelStatusSayHelmsman: ps = PersonnelStatusSayHelmsmanCourse            
                case PersonnelStatusSayRoot: ps = PersonnelStatusSayHelmsman
                default: return nil
            }
            
        case '4':
            switch (g.ThePlayer.PersonnelStatus) {                
                case PersonnelStatusSayHelmsman: ps = PersonnelStatusSayHelmsmanDestination                
                default: return nil
            }
            
        case '7':
            switch (g.ThePlayer.PersonnelStatus) {                
                case PersonnelStatusSayHelmsman: // Evasive Action
                    helmsman := g.PlayerShip.GetCrewMemberForRole(CrewRoleHelmsman)
                    g.EnqueueMessage(NewMessageEvasiveAction(g.ThePlayer.CrewMember, helmsman))
                    ps = PersonnelStatusNormal
                    result.TicksToPass = 1 
                
                default: return nil
            }
        
            
        case '\\':    
            switch (g.ThePlayer.PersonnelStatus) {                
                case PersonnelStatusSayHelmsmanDestination:
                    // ---- this should be in some helper method that can choose from a randomized list
                    //      of phrases. Plus when the commander is Ai, he should reuse the same logic
                    // ---- and needs to be made gender aware (another good reason for helper functions)
                    // ---- so the from crew and to crew needs to be passed in (as well as the target)            
                    helmsman := g.PlayerShip.GetCrewMemberForRole(CrewRoleHelmsman)
                    g.EnqueueMessage(NewMessageSetDestinationTargeter(g.ThePlayer.CrewMember, helmsman, g.ThePlayer.PersonnelHelmsmanTarget))
                    
                    ps = PersonnelStatusNormal
                    result.TicksToPass = 1 

                case PersonnelStatusSayHelmsmanCourse:
                    helmsman := g.PlayerShip.GetCrewMemberForRole(CrewRoleHelmsman)
                    g.EnqueueMessage(NewMessageSetCourse(g.ThePlayer.CrewMember, helmsman, g.ThePlayer.PersonnelHelmsmanCourseSetter))
                    
                    ps = PersonnelStatusNormal
                    result.TicksToPass = 1 
                
                
                default: return nil
            }
            
        default:
            return nil
    }
    
    g.ThePlayer.PersonnelStatus = ps
    
    return result
}

func (p *PanelPersonnel) processInputNormal(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult { Exit: false, TicksToPass: 0 }
    
    //handledRune := true
    switch ch {
        case '1': 
            g.ThePlayer.PersonnelStatus = PersonnelStatusSayRoot
        default:
            //handledRune = false
            return nil
    }
    /* if (handledRune) { return result }

    switch key {
        case termbox.KeyArrowLeft:
        default:
            return nil
    }
*/
    return result
}


func (p *PanelPersonnel) Display(g *Game, r *ConsoleRange) {
    r.SetBorder()
    r.SetTitle("Personnel")
        
    switch (g.ThePlayer.PersonnelStatus) {
        case PersonnelStatusNormal: p.displayNormal(g,r)        
        default:            
            if (g.ThePlayer.CrewMember.CrewRole == CrewRoleHelmsman) {
                p.displayHelmsmanSay(g,r)
            } else if (g.ThePlayer.CrewMember.CrewRole == CrewRoleCommander) {
                p.displayCommanderSay(g,r)
            }        
    }
}

const black = termbox.ColorBlack
const green = termbox.ColorGreen | termbox.AttrBold
const yellow = termbox.ColorYellow | termbox.AttrBold

func (p *PanelPersonnel) displayHelmsmanSay(g *Game, r *ConsoleRange) {
    r.Com("[1]", " Go Back", 1,2,black, green)
    
}

func (p *PanelPersonnel) displayCommanderSay(g *Game, r *ConsoleRange) {
    r.Com("[1]", " Go Back", 1,2,black, green)
    
    // who else is there to send messages to?
    // what level is this at?
    
    // Just hard-code it right now. Until I work on it further, I cannot make informed decisions on this.
    
    helmsman := g.PlayerShip.GetCrewMemberForRole(CrewRoleHelmsman)
    
    switch (g.ThePlayer.PersonnelStatus) {
        case PersonnelStatusSayRoot:                
            r.DisplayText("> Say",1,1)
            r.Com("[2]", fmt.Sprintf(" Helmsman %s", helmsman.GetFullName()), 1, 3, black, green)
            
        case PersonnelStatusSayHelmsman:              
            r.DisplayText("> Say > Helmsman",1,1)        
            r.Com("[2]", " Bring us to a stop", 1, 3, black, yellow)
            r.Com("[3]", " Set course for ...", 1,4, black, green) 
            r.Com("[4]", " Set destination to ...", 1, 5, black, green)  
            // manuever
            // jump k
            r.Com("[7]", " Take evasive action", 1, 8, black, yellow)
            
            // 5 evasive action
            
        case PersonnelStatusSayHelmsmanCourse:
            r.DisplayText("> Say > Helmsman > Course", 1, 1)
            g.ThePlayer.PersonnelHelmsmanCourseSetter.Display(r, black, green, termbox.ColorBlue, termbox.ColorWhite | termbox.AttrBold)
            
            r.Com("[\\]", " Say", 2, 9, black, yellow)
            
        case PersonnelStatusSayHelmsmanDestination:
            r.DisplayText("> Say > Helmsman > Destination",1,1)        
            // display target selection just like helm control
            g.ThePlayer.PersonnelHelmsmanTarget.Display(r, black, green)
            
            if (g.ThePlayer.PersonnelHelmsmanTarget.IsValidTarget) {
                r.Com("[\\]", " Say", 2, 9, black, yellow)
                //r.DisplayText("[\\] Say", 2, 9)
                //r.DisplayTextWithColor("[\\]", 2, 9, termbox.ColorBlack, termbox.ColorYellow | termbox.AttrBold)
            }
            
        default: 
            r.DisplayText("ERROR", 2,3)
    }
}

func (p *PanelPersonnel) displayNormal(g *Game, r *ConsoleRange) {
    w, h := r.GetSize()
    
    // key for displaying in full screen mode to scroll through all previous messages
    
    // don't display in full screen mode
    r.Com("[1]", " Say Something", 2, 1, black, green)
    
    // ---- convert this to display many messages
    // ---- dim old messages
    // display the front message if there is one
    pcm := g.ThePlayer.CrewMember
    
    
    fd := NewFlowDocument(w - 2, h - 3)
    //stop := false
    //m := pcm.ReceivedMessages.Back()    
    
    for e := pcm.ReceivedMessages.Back(); e != nil && !fd.IsFull(); e = e.Prev() {
        m := e.Value.(*CrewMessage)
        
        // See if this message was FROM the player, TO the player, overheard, or a status message
        daString := ""
        fg := termbox.ColorBlack | termbox.AttrBold
        if (g.ThePlayer.CrewMember == m.To) {        
            if (m.TickReceived +1 == g.tick) {
                fg = termbox.ColorWhite | termbox.AttrBold
            } else if (m.TickReceived + 7 > g.tick) {
                fg = termbox.ColorWhite
            }            
            if (m.From != nil) {
                fcm := m.From
                //fd.AddParagraph(fmt.Sprintf("[FROM %s %s, %s]", fcm.FirstName, fcm.LastName, fcm.CrewRole.ToString()), fg, termbox.ColorBlack)            
                daString = fmt.Sprintf("%s %s, %s, says: ", fcm.FirstName, fcm.LastName, fcm.CrewRole.ToString())
            }
            
        } else if (g.ThePlayer.CrewMember == m.From) {
            daString = "You said: "
        
        } else if (m.IsStatusMessage) {            
            if (m.TickReceived +1 == g.tick) {
                fg = termbox.ColorBlue | termbox.AttrBold
            } else {
                fg = termbox.ColorBlue
            }
        
        } else { // overheard message
            daString = "You overheard: "
        }
        
        // In the normal view, it does not show the tick received. In the full view, it will.
        fd.AddParagraph(daString + m.Message, fg, termbox.ColorBlack)
    }
    
        
        

    fd.Write(r, 1, 2)
}
