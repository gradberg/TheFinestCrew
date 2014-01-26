/*
    A crew message represents any order or information passed between crew members. 
*/

package main
import "fmt"

type CrewMessage struct {
    From *CrewMember
    To *CrewMember // should not be nil, as that indicates an Ai (or player?) tried to send something to someone who doesn't exist.
    
    Information InformationEnum
    // information enum // Used by AI to 'understand' a message the player or another AI sent.
    // information object pointer // used to transmit whatever data necessary for advanced messages
    
    Message string // display message that players see
    
    Target ISpaceObject // used for any messages that specify a destination or target space object
    Course float64 // used for setting course
    Speed float64 // used for setting course
    
    TickReceived int // the tick in which this message was received by the player (which is the END of the tick it was sent)
}

func NewCrewMessage(from, to *CrewMember, message string) *CrewMessage {
    return &CrewMessage {
        From: from,
        To: to,
        Message: message,
        Information: InformationNone,
    }
}

type InformationEnum int; const (
    InformationNone InformationEnum = iota
    
    // # To helmsman #
    InformationFullStop // indicate to stop the ship
    InformationSetCourse // indicate to set a course
    InformationSetDestination // indicate to set a destination for the ship
    InformationEvasiveAction 
    
    
    //InformationAwaitingOrders // indicates to the captain that a given crewmen is free for new orders
)


func NewMessageFullStop(from, to *CrewMember) *CrewMessage {
    m := NewCrewMessage(from, to,            
        fmt.Sprintf(
            "Mr. %s, full stop please.", 
            to.LastName, 
        ),
    )
    m.Information = InformationFullStop
    return m
}
func NewMessageSetCourse(from, to *CrewMember, courseSetter *CourseSetter) *CrewMessage {
    m := NewCrewMessage(from, to,            
        fmt.Sprintf(
            "Mr. %s, set course for %03f°, %4.f ticks.", 
            to.LastName, 
            courseSetter.Course,
            courseSetter.Speed,
        ),
    )
    m.Information = InformationSetCourse
    m.Course = courseSetter.Course
    m.Speed = courseSetter.Speed
    return m
}
func NewMessageSetDestination(from, to *CrewMember, targeter *Targeter) *CrewMessage {
    m := NewCrewMessage(from, to,            
        fmt.Sprintf(
            "Mr. %s, make our destination the %s %s.", 
            to.LastName, 
            targeter.TargetType.ToString(),
            targeter.DesiredTarget.GetName(),
        ),
    )
    m.Information = InformationSetDestination
    m.Target = targeter.DesiredTarget
    return m
}
func NewMessageEvasiveAction(from, to *CrewMember) *CrewMessage {
    m := NewCrewMessage(from, to,            
        fmt.Sprintf("%s get us out of here!", to.LastName ),
    )
    m.Information = InformationEvasiveAction
    return m
}








