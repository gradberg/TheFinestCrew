/*
    A helmsman AI crewmember how can take normal orders and fly the ship around
*/

package main
import "fmt"
import "math"

type AiHelmsmanBasic struct {       
    // Basic helmsman alerts the command/captain if he has not been tasked with some mission.
    // Have some nice information here is both good for theme, *AND* it serves a purpose of
    // providing a variety of information to the captain that he would otherwise not be able to see.

    ticksUntilNextWait int

    standingOrders *CrewMessage // holds the message the represents the helmsman's last orders.
    
    evasiveWeaveCounter int // used to dodge back and forth during evasive actions.
    evasiveWeaveAngle float64 // additional angle added to weaving back and forth
        
    // ---- need to define some sort of pluggable personality/action system,so that to add various
    // capabilities (the ability to say something when there are no orders) you simply reference
    // additionally objects that get passed the state on DoAction and may or may not be the AI's turn.
    
    // ---- for the snarky responses with new information/orders received, I can make an array of 
    // function pointers to the various checks to send unique responses back. That at least can reduce 
    // the clutter in the main processing function
}

func (ai *AiHelmsmanBasic) ClearEphemeralState() {
    ai.ticksUntilNextWait = 0
    ai.standingOrders = nil    
}

func (ai *AiHelmsmanBasic) DoAction(s *Ship, g *Game, cm *CrewMember) int {
    if (ai.standingOrders == nil) {
        if (ai.ticksUntilNextWait <= 0) {
            ai.randomizeNextMessage(g)
            
            // Tell the commander that i've got nothing to do.
            commanderRole := s.GetCommandingOfficerRole()
            g.EnqueueMessage(NewCrewMessage(
                cm, 
                s.GetCrewMemberForRole(commanderRole), 
                //fmt.Sprintf("Awaiting orders %s. I am the very model of a modern major general. I've information vegetable, animal, and mineral.", commanderRole.ToString()),
                fmt.Sprintf("Awaiting orders %s.", commanderRole.ToString()),
            ))
        }
        ai.ticksUntilNextWait --
    } else {
        // set this to zero, so when he runs out of orders (if that is posible), then he will
        // immediately inform the commanding officer        
        ai.ticksUntilNextWait = 0 
    }
    
    var newOrders *CrewMessage
    // Go through received messages, and act on any that are relevant.
    for e := cm.GetFirstNewMessage(g); e != nil; e = e.Next() {
        m := e.Value.(*CrewMessage)
    
        switch (m.Information) {
            // assign as new orders
            case InformationFullStop, InformationSetDestination, InformationSetCourse: newOrders = m
            case InformationEvasiveAction:
                // ---- clear evasive action state, as it is starting over
                newOrders = m
            default: // discard message. Its useless to the AI
        }
    }
    
    // if received new orders, then acknowledge that to the sender.
    if (newOrders != nil) {        
        // ---- should this take up a tick? It would suck if every crew member took up time before doing their
        // orders just to acknowledge them. But then again, people have the ability to speak AND do things concurrently,
        // so it would make sense for informative messages to not take up time.
        
        // Sarcastic responses, generally due to a duplicate order    
        shouldSendNormalResponse := true
        if (ai.standingOrders != nil) {
            switch (ai.standingOrders.Information) {
                case InformationFullStop:
                    if (newOrders.Information == InformationFullStop) {
                        g.EnqueueMessage(NewCrewMessage(
                                cm, 
                                newOrders.From, 
                                fmt.Sprintf("We are already stopping, %s.", newOrders.From.CrewRole.ToString()),
                        ))                        
                        shouldSendNormalResponse = false
                    }
                
                case InformationSetCourse:
                    if (newOrders.Information == InformationSetCourse &&
                        newOrders.Course == ai.standingOrders.Course &&
                        newOrders.Speed == ai.standingOrders.Speed) {
                        g.EnqueueMessage(NewCrewMessage(
                                cm, 
                                newOrders.From, 
                                fmt.Sprintf("We have already set that course, %s.", newOrders.From.CrewRole.ToString()),
                        ))      
                        shouldSendNormalResponse = false
                    }
                
                case InformationSetDestination:
                    if (newOrders.Information == InformationSetDestination && newOrders.Target == ai.standingOrders.Target) {
                        g.EnqueueMessage(NewCrewMessage(
                                cm, 
                                newOrders.From, 
                                fmt.Sprintf("We have already set destination to %s.", ai.standingOrders.Target.GetName()),
                        ))                        
                        shouldSendNormalResponse = false
                    }
                    
                case InformationEvasiveAction:
                    if (newOrders.Information == InformationEvasiveAction) {
                        g.EnqueueMessage(NewCrewMessage(
                                cm, 
                                newOrders.From, 
                                fmt.Sprintf("I am trying to %s.", newOrders.From.CrewRole.ToString()),
                        ))                        
                        shouldSendNormalResponse = false
                    }
            }
        }
        // More snarky responses
        if (shouldSendNormalResponse && newOrders.Information == InformationFullStop && s.SpeedInUnitsPerTick == 0) {
            g.EnqueueMessage(NewCrewMessage(
                    cm, 
                    newOrders.From, 
                    fmt.Sprintf("We already stopped %s.", newOrders.From.CrewRole.ToString()),
            ))                                
            newOrders = nil
            shouldSendNormalResponse = false
        } else if (shouldSendNormalResponse && newOrders.Information == InformationEvasiveAction && ai.areThereAnyEnemiesLeft(s,g,cm) == false) {            
            g.EnqueueMessage(NewCrewMessage(
                    cm, 
                    newOrders.From, 
                    fmt.Sprintf("There are no threats to evade %s.", newOrders.From.CrewRole.ToString()),
            ))                                
            newOrders= nil
            shouldSendNormalResponse = false
        }
                
        // "NORMAL" responses (as in just acknowledging what the person said)
        if (shouldSendNormalResponse) {
            switch (newOrders.Information) {
                case InformationEvasiveAction:                    
                    g.EnqueueMessage(NewCrewMessage(
                            cm, 
                            newOrders.From, 
                            "Yes sir!",
                    ))
        
                default:
                    g.EnqueueMessage(NewCrewMessage(
                            cm, 
                            newOrders.From, 
                            fmt.Sprintf("Aye aye, %s.", newOrders.From.CrewRole.ToString()),
                    ))
            }
        }
        
        ai.standingOrders = newOrders
    }
        
    // Act on any standing orders (which may or may not be new at this point)
    if (ai.standingOrders != nil) {
        sh := s.Helm
        switch (ai.standingOrders.Information) {
            case InformationSetCourse:
                sh.AutoPilotPower = 3 // full power
                
                
                if (sh.AutoPilotMode == AutoPilotModeCourse &&
                    sh.AutoPilotDesiredCourse == ai.standingOrders.Course &&
                    sh.AutoPilotDesiredSpeed == ai.standingOrders.Speed) {
                    // do nothing if already set
                    return 0
                }
                
                sh.AutoPilotMode = AutoPilotModeCourse
                sh.AutoPilotDesiredCourse = ai.standingOrders.Course
                sh.AutoPilotDesiredSpeed = ai.standingOrders.Speed                
                return 1    
        
            case InformationSetDestination:
                sh.AutoPilotPower = 3 // full power
            
                // check if the destination is already set to this. If so, do nothing.                 
                if (sh.AutoPilotMode == AutoPilotModeTarget && sh.AutoPilotDesiredTarget == ai.standingOrders.Target) {
                    // do nothing
                    return 0
                }                
                
                // otherwise set the target, which takes up a tick
                sh.AutoPilotMode = AutoPilotModeTarget
                sh.AutoPilotDesiredTarget = ai.standingOrders.Target                
                return 1
        
            case InformationFullStop:              
                sh.AutoPilotPower = 3 // full power

                // if the ship has stopped, inform the captain and clear orders
                if (sh.Ship.SpeedInUnitsPerTick == 0.0) {
                    g.EnqueueMessage(NewCrewMessage(
                            cm, 
                            ai.standingOrders.From, 
                            fmt.Sprintf("%s, reached full stop.", ai.standingOrders.From.CrewRole.ToString()),
                    ))                    
                    ai.standingOrders = nil
                    ai.randomizeNextMessage(g)
                    return 1 // 1, because this represents useful information to the captain
                }                
                
                // Check if full stop is already set. If so, do nothing (can report to captain that ship stopped)
                if (sh.AutoPilotMode == AutoPilotModeFullStop) {
                    // do nothing
                    return 0
                }
                
                sh.AutoPilotMode = AutoPilotModeFullStop
                return 1
        
            case InformationEvasiveAction: return ai.doEvasiveAction(s, g, cm)
        
            default: // Do nothing. Unrecognized :P
                return 0
        }    
    }
    
    // if it got here... it did 
    return 0
}

func (ai *AiHelmsmanBasic) randomizeNextMessage(g *Game) {
    ai.ticksUntilNextWait = 10 + int(g.Rand.Float64() * 100.0)    
}

func (ai *AiHelmsmanBasic) doEvasiveAction(s *Ship, g *Game, cm *CrewMember) int {
    // If there are no enemies, inform the captain and clear orders.
    if (ai.areThereAnyEnemiesLeft(s, g, cm) == false) {
        g.EnqueueMessage(NewCrewMessage(
            cm, 
            ai.standingOrders.From, 
            fmt.Sprintf("%s, no more threats to evade.", ai.standingOrders.From.CrewRole.ToString()),
        ))      
        ai.standingOrders = nil
        return 0
    }
    
    // Otherwise, figure out where all the enemies are
    evasiveAngle := ai.computeEvasiveHeading(s, g, cm)
        
    // Update the weave counter (which causes the ship to go back and forth in a weaving fashion to avoid shots
    ai.evasiveWeaveCounter--
    if (ai.evasiveWeaveCounter <= 0) {
        ai.evasiveWeaveCounter = 10 + int(g.Rand.Float64() * 20)        
        if (ai.evasiveWeaveAngle > 0.0) {
            ai.evasiveWeaveAngle = -20.0 + g.Rand.Float64() * -65.0
        } else {
            ai.evasiveWeaveAngle = +20.0 + g.Rand.Float64() * +65.0        
        }
    }    
    
    // Orient the ship appropriately
    sh := s.Helm
    sh.AutoPilotMode = AutoPilotModeCourse
    sh.AutoPilotDesiredCourse = AddAngles(evasiveAngle, ai.evasiveWeaveAngle)
    sh.AutoPilotDesiredSpeed = 1000.0 // something ludriciously high for full power
    sh.AutoPilotPower = 3

    return 1
}

func (ai *AiHelmsmanBasic) areThereAnyEnemiesLeft(s *Ship, g *Game, cm *CrewMember) bool {
    for e := g.Ships.Front(); e != nil; e = e.Next() {
        ship := e.Value.(*Ship)        
        // ---- check for friend vs foe?
        
        if (ship.HitPoints > 0.0 && ship != s) {
            // Must not be this ship
            return true
        }        
    }
    return false
}

// 
// Loops over all ships, assigning them a threat level and sticking that in their heading bucket.
// It then determines which part of the circle (0 to 360) has the worst threat and returns an 
// angle opposite of that.
func (ai *AiHelmsmanBasic) computeEvasiveHeading(s *Ship, g *Game, cm *CrewMember) float64 {
    // Loop over all the possible angles, and for each one, compute the threat presented by other
    // ships. At the end, whichever one at the least threat is the 'winner'
    minimumAngle := 0.0
    minimumThreat := 1000000.0
    for a := 0; a < 360; a += 3 {
        totalThreat := 0.0
        
        for e := g.Ships.Front(); e != nil; e = e.Next() {
            ss := e.Value.(*Ship)
            if (ss == s || ss.IsDestroyed()) { continue }
                        
            // Take the cosine of the angle difference (which is what vector dot products use)
            // If that value is negative, it means the ship is on the other side and doesn't count.
            // Take that value, and multiply it by the threat rating and add it up
            angleInDegrees, distance := ss.Point.Subtract(s.Point).ToVector()
            radians := AddAngles(float64(a), -angleInDegrees) * math.Pi / 180.0
            multiplier := math.Cos(radians) + 0.5 // Adding 0.5 means it takes into account the everything but the back 120 degrees
            if (multiplier < 0.0) { continue } // discard as it is on the wrong side
            
            shipDanger := 1.0 // ---- will be used in the future to compute based on the ship design
            threat := multiplier * shipDanger / distance
            totalThreat += threat
        }
        
        if (totalThreat < minimumThreat) {
            minimumAngle = float64(a)
            minimumThreat = totalThreat
        }
    }    
    return minimumAngle
}
