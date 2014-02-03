/*
    Contains setup methods for supporting the 'survival' game mode.
*/

package main

func (g *Game) setupForStart() {
    g.createPlayerShip()
}

// returns a boolean indicating if the next round should start.
func (g *Game) checkForNextRound() bool {
    // ---- see if there are no non-destroyed non-player-ships left.
    return false
}

// returns a boolean indicating that the player lost the game. (Presumably their
// ship was destroyed)
func (g *Game) checkForLost() bool {
    return false
}

func (g *Game) setupForNextRound() {
    // Clear the state in the game
    g.ClearEphemeralState()
    
    // Increase the round number. Note, this does mean the player sees the round start at "1"
    g.round++

    // Create some planets for (pointless) flavor    
    g.Planets.PushBack(&Planet {
        Point: NewPoint(g.Rand.Float64() * 2000.0 - 1000.0, g.Rand.Float64() * 2000.0 - 1000.0).Round(),
        Name: "Jupiter",
    })
    g.Planets.PushBack(&Planet {
        Point: NewPoint(g.Rand.Float64() * 2000.0 - 1000.0, g.Rand.Float64() * 2000.0 - 1000.0).Round(),
        Name: "Mars",
    })    
    g.Planets.PushBack(&Planet {    
        Point: NewPoint(g.Rand.Float64() * 2000.0 - 1000.0, g.Rand.Float64() * 2000.0 - 1000.0).Round(),
        Name: "Saturn",
    })
    
    // Add the player's existing ship to the list
    g.Ships.PushBack(g.PlayerShip)
    
    // Add a number of pirate ships based on the round + 1
    for i := 0; i <= g.round; i++ {
        g.Ships.PushBack(g.createRandomAiPirateFighter())
    }    
}


func (g *Game) createRandomAiPirateFighter() *Ship {    
    aiShip := createPirateFighter("Pirate Fither")
    aiShip.Point = NewPoint(g.Rand.Float64() * 800.0 - 400.0, g.Rand.Float64() * 800.0 - 400.0).Round()
    aiShip.CrewMembers.PushBack(NewCrewMember("Unknown", "Scoundrel", NewAiPiratePilot(g), CrewRolePilot))
    return aiShip
}