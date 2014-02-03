/*
    A non-production helmsman AI crewmember for piloting ships around
*/

package main

type PiratePilotStateEnum int; const (
    PiratePilotStateCalm PiratePilotStateEnum = iota
    PiratePilotStateAlarmSetHelmTarget
    PiratePilotStateAlarmSetWeaponTarget
    PiratePilotStateAlarmEnd    
)

type AiPiratePilot struct {       

    State PiratePilotStateEnum
    
    alarmDistance float64

    // When calm, the ship flies around randomly    
    desiredHeading float64
    ticksUntilNextTurn int
}
func NewAiPiratePilot(g *Game) *AiPiratePilot { 
    ai := &AiPiratePilot { } 
    ai.alarmDistance = 100.0 + g.Rand.Float64() * 100.0
    return ai
}

func (ai *AiPiratePilot) ClearEphemeralState() {
    ai.ticksUntilNextTurn = 0
    ai.State = PiratePilotStateCalm
}

func (ai *AiPiratePilot) DoAction(s *Ship, g *Game, cm *CrewMember) int {
    ai.doStateTransitions(s,g,cm)
    
    switch ai.State {
        case PiratePilotStateCalm:
            return ai.doCalm(s,g,cm)
            
        case PiratePilotStateAlarmSetHelmTarget: 
            s.Helm.IsDirectPilot = false
            s.Helm.AutoPilotMode = AutoPilotModeTarget
            s.Helm.AutoPilotDesiredTarget = g.PlayerShip
            return 1
            
        case PiratePilotStateAlarmSetWeaponTarget: 
            w := s.Weapons.Front().Value.(*ShipWeapon)
            w.AutoFire = true
            w.TargetType = TargetTypeTarget
            w.TargetShip = g.PlayerShip
            return 1           
            
        default:
            // do nothing
            return 0
    }
}

func (ai *AiPiratePilot) doStateTransitions(s *Ship, g *Game, cm *CrewMember) {
    switch ai.State {
        case PiratePilotStateCalm:
            // if the enemy ship is close enough, then the pirate becomes alarmed
            _, distance := s.Point.Subtract(g.PlayerShip.Point).ToVector()
            if ((distance <= ai.alarmDistance) || s.WasHit) {
                ai.State = PiratePilotStateAlarmSetHelmTarget
            }
            
        case PiratePilotStateAlarmSetHelmTarget: 
            ai.State = PiratePilotStateAlarmSetWeaponTarget
        case PiratePilotStateAlarmSetWeaponTarget: 
            ai.State = PiratePilotStateAlarmEnd
            
    }
}

func (ai *AiPiratePilot) doCalm(s *Ship, g *Game, cm *CrewMember) int {
    if (ai.ticksUntilNextTurn == 0) {
        // ---- pick a new heading
        ai.ticksUntilNextTurn = 10 + int(g.Rand.Float64() * 100.0)
        ai.desiredHeading = g.Rand.Float64() * 360.0
    }
    ai.ticksUntilNextTurn --
    
    // full auto-thrust, all of the time.
    s.Helm.AutoThrust = true     
    
    // turn the ship if it is not already on the right heading
    s.Helm.TurnTowardsAngle(ai.desiredHeading, 100)
    return 1    
}

