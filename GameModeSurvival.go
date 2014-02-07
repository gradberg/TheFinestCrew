/*
    Contains setup methods for supporting the 'survival' game mode.
*/

package main

func (g *Game) setupForStart() {
    g.createPlayerShip()
}

// returns a boolean indicating if the next round should start.
func (g *Game) checkForNextRound() bool {    
    isNextRound := true
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)
        
        if (s != g.PlayerShip && s.IsDestroyed() == false) {
            isNextRound = false
        }
    }
    
    // If the last enemey was JUST destroyed, delay one turn so that
    // the player can see that it was destroyed just prior to 
    if (isNextRound && !g.waitATurn) {
        g.waitATurn = true
        isNextRound = false
    }
    
    if (isNextRound) {
        g.totalDestroyedShips()
        g.waitATurn = false
    }    
    return isNextRound
}

// returns a boolean indicating that the player lost the game. (Presumably their
// ship was destroyed)
func (g *Game) checkForLost() bool {    
    isDead := g.PlayerShip != nil && g.PlayerShip.IsDestroyed()
    
    // Just like round-over, wait one turn after the player is destroyed so that they can
    // see that their ship is destroyed. In the future, this will maybe blank out their
    // controls or whatnot, but for now its fine.
    // (Or better yet, their various panels would show random destruction on them. It'd
    //  be cool)
    if (isDead && !g.waitATurn) {
        g.waitATurn = true
        isDead = false
    }    
    
    if (isDead) {
        g.totalDestroyedShips()
        g.waitATurn = false
    }
    return isDead
}

func (g *Game) totalDestroyedShips() {
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)
        
        if (s != g.PlayerShip && s.IsDestroyed()) {
            g.kills++
        }
    }
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
    aiShip := createPirateFighter("Pirate Fighter")
    aiShip.Point = NewPoint(g.Rand.Float64() * 800.0 - 400.0, g.Rand.Float64() * 800.0 - 400.0).Round()
    aiShip.CrewMembers.PushBack(NewCrewMember("Unknown", "Scoundrel", NewAiPiratePilot(g), CrewRolePilot))
    return aiShip
}