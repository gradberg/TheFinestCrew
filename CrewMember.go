

package main

import "container/list"

type CrewRoleEnum int
const (
    CrewRoleError CrewRoleEnum = iota // should not ever be used. Means something was unintiliazed.
    CrewRolePilot // a pilot is the sole member of a fighter
    CrewRoleHelmsman // Runs the helm of any ship that has multiple crew members
    CrewRoleCommander // Combination captain and tactical officer of tiny ships
)
func (e CrewRoleEnum) ToString() string {
    switch (e) {
        case CrewRolePilot: return "Pilot"
        case CrewRoleHelmsman: return "Helmsman"
        case CrewRoleCommander: return "Commander"
        default: return "ERROR"
    }
}

type CrewMember struct {
    FirstName string // 
    LastName string // used in captain-to-crew communications
    // gender?
    IsPlayer bool // indicates this is actually the player
    Ai IAiCrew // stores the AI that this crew member follows, if they aren't the player
    CrewRole CrewRoleEnum // used to send orders to the correct crew member on a ship
    ReceivedMessages *list.List
    // sent messages
    // overheard messages
}

func NewCrewMember(firstName string, lastName string, ai IAiCrew, crewRole CrewRoleEnum) *CrewMember { return &CrewMember { 
    FirstName: firstName,
    LastName: lastName,
    Ai: ai,
    CrewRole: crewRole,
    ReceivedMessages: list.New(),
} }

func (c *CrewMember) GetFullName() string {
    return c.FirstName + " " + c.LastName
}

// Returns a list element containing the first message received on the past tick.
// That element can then be used to loop through all received messages.
func (c *CrewMember) GetFirstNewMessage(g *Game) *list.Element {
    var firstElement *list.Element
    for e := c.ReceivedMessages.Back(); e != nil; e = e.Prev() {
        m := e.Value.(*CrewMessage)
        if (m.TickReceived == g.tick - 1) {
            firstElement = e
        } else {
            // break if this message is older
             break
        }
    }
    return firstElement
}