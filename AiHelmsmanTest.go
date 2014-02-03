/*
    A non-production helmsman AI crewmember for piloting ships around
*/

package main

type AiHelmsmanTest struct {       
    desiredHeading float64
    ticksUntilNextTurn int
}

func (ai *AiHelmsmanTest) ClearEphemeralState() {
    ai.ticksUntilNextTurn = 0
}

func (ai *AiHelmsmanTest) DoAction(s *Ship, g *Game, cm *CrewMember) int {
    if (ai.ticksUntilNextTurn == 0) {
        // ---- pick a new heading
        ai.ticksUntilNextTurn = 10 + int(g.Rand.Float64() * 100.0)
        ai.desiredHeading = g.Rand.Float64() * 360.0
    }
    ai.ticksUntilNextTurn --
    
    // turn the ship if it is not already on the right heading
    s.Helm.TurnTowardsAngle(ai.desiredHeading, 100)
    
    // full auto-thrust, all of the time.
  //  s.Helm.AutoThrust = true 

    return 1
}
