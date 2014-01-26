/*
    The Game class represents the game being played, and contains the main
    game loop and all relevant state for playing the game.
    
    Any display-specific knowledge should be kept in the GameDisplayAdapter.
*/

package main

import "container/list"
import "math/rand"
import "time"

// Represents the game being played. 
type Game struct {
    tick int    
    
    ThePlayer *Player
    PlayerShip *Ship    
    
    Ships *list.List
    
    message string
    
    Planets *list.List
    
    Rand *rand.Rand    
    
    Projectiles *list.List
    
    PendingMessages *list.List
}

// Creates a new game object so that the actual game lifecycle can begin
func NewGame() *Game {
    // create random number generator
    source := rand.NewSource(int64(time.Now().Nanosecond()))    
    
    g := &Game { 
        tick: 0,
        message: "",
        Planets: list.New(),
        Ships: list.New(),
        Rand: rand.New(source),
        PendingMessages: list.New(),
        Projectiles: list.New(),
    }    
    g.ThePlayer = NewPlayer(g)
    
    return g
}

// Blocking call which contains the main game loop, and exits when the user has chosen to quit.
func (g *Game) Run() {
    AsyncInput_Init()
    defer AsyncInput_Close()

    // Development setup
    g.doSetupForDevelopment()

    for {
        g.Display()
        
        g.message = ""
        var ir *InputResult = g.DoInput()
        if ir.Exit { return }
        
        g.processTick(ir)
    }    
}



func (g *Game) processTick(ir *InputResult) {    
    // Loop over the number of ticks the input says it will take, and do the processing for each one?
    
    var tickCount uint
    for tickCount = 0; tickCount < ir.TicksToPass; tickCount++ {
        // first, any actions characters did should occur first. That way when players fire weapons, 
        // their choices will be immediately carried-out
        g.processAiCharacters()
        
        // Next, process weapons, which both generates projectiles, and has them hit ships. This occurs
        // before ship movement so that in their first turn, projectiles should be able to impact
        // a ship where it was at the start of a turn. After that it is a crap shoot
        g.processWeapons()
        
        // Last process ship movement, as they are essentially the slowest objects.
        g.processShipMovement()        
        
        g.tick++
    }       
}

func (g *Game) processAiCharacters() {
    // go through each ship's crewmembers. 
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)        
        
        // Go through each crew member and process their turn
        for ce := s.CrewMembers.Front(); ce != nil; ce = ce.Next() {
            c := ce.Value.(*CrewMember)
            
            if (c.IsPlayer) { continue }
            
            // ---- use the ticks number to control how many turns this crew member loses?
            c.Ai.DoAction(s,g,c)
        }
        
        // do other per-ship tasks
        s.Helm.DoAutoPilot()
        s.WasHit = false
    }   
    
    // go through each message and enqueue it in its receipients lists        
    for e := g.PendingMessages.Front(); e != nil; e = e.Next() {
        m := e.Value.(*CrewMessage)
        m.TickReceived = g.tick
        m.To.ReceivedMessages.PushBack(m)
    }
    g.PendingMessages.Init() // clear the list
}

func (g *Game) processWeapons() {
    // loop over each ship, firing any weapons as appropriate
    //   projectiles spawn projectiles
    //   lasers hit immediately.
    // go through each ship's crewmembers. 
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)
        
        for we := s.Weapons.Front(); we != nil; we = we.Next() {
            w := we.Value.(*ShipWeapon)
                        
            // if its not ready to fire, just cycle it.
            if (w.CurrentCycle > 0) {
                w.CurrentCycle --
                continue
            }
            
            // If the weapon is not set to fire... don't fire!
            if (!w.AutoFire) { continue }
            
            // create new projectile (if set to fire?)?
            // At some point, call the FireControl object or the ship and ask it to generate projectiles.
     
            w.SetFiringAngle(s)
            shotAngle := AddAngles(w.FiringAngle, s.ShipHeadingInDegrees)
            
            p := &Projectile {
                Point: s.Point,
                Heading: shotAngle,
                Speed: w.DesignSpeed,
                OriginShip: s,
                DesignSpeed: w.DesignSpeed,
                DesignDrag: w.DesignDrag,
                DesignDamage: w.DesignDamage,
            }
            g.Projectiles.PushBack(p)
            
            // ---- will be more useful once I add random deflection in
            LogCalc("Projectile Fired from ship. sHeading [%f], pHeading [%f]", s.ShipHeadingInDegrees, p.Heading)
            
            // Set the weapon back to its reload setting
            w.CurrentCycle = w.DesignCycle
            
        }
        
        
        
    }   

    projectilesToRemove := list.New()
    // loop over every projectile, moving it and having it impact ships as appropriate    
    for pe := g.Projectiles.Front(); pe != nil; pe = pe.Next() {
        p := pe.Value.(*Projectile)
        
        // see if it impacts
        projectilePoint1 := p.Point
        projectilePoint2 := p.GetFuturePoint()
        
        // see if it impacts any ships
        didImpact := false
        for se := g.Ships.Front(); se != nil; se = se.Next() {
            s := se.Value.(*Ship)
            if (p.OriginShip == s) { continue }
            
            // If the distance from the line segment to the ship is LESS than the ship's size, then it impacts            
            distance := DistanceFromLineSegment(s.Point, projectilePoint1, projectilePoint2)        
            if (distance <= s.HitSize) {
                // if so, damage the ship and add this to the remove list    
                damage := p.DesignDamage * p.Speed / p.DesignSpeed
                s.HitPoints = Round(s.HitPoints - damage, 1)
                
                LogCalc("Projectile  *HIT* Ship %s at %f, %f by projectile at %f, %f, pHeading %f, pSpeed %f, distance from path %f",
                    s.Name, s.Point.X(), s.Point.Y(), p.Point.X(), p.Point.Y(), p.Heading, p.Speed, distance,
                )                    
                
                // Transfer impact to the ship's momentum IF IT IS A PROJECTILE OR MISSILE
                s.DoAcceleration(p.Heading, damage / 5.0)
                
                s.WasHit = true
            
                didImpact = true
                break // Stop looking for other ships to impact
            } else {            
            
                LogCalc("Projectile missed Ship %s at %f, %f by projectile at %f, %f, pHeading %f, pSpeed %f, distance from path %f",
                    s.Name, s.Point.X(), s.Point.Y(), p.Point.X(), p.Point.Y(), p.Heading, p.Speed, distance,
                )                    
            }
        }
        
        if (!didImpact) {        
            // otherwise, move it and incur drag
            p.DoMovement()        
        }
        
        // If it has stopped, add it to the remove list.
        if (p.Speed <= 0.0 || didImpact) { projectilesToRemove.PushBack(pe) }
    }
    
    // remove any projectiles that impacted or timed-out
    for pe := projectilesToRemove.Front(); pe != nil; pe = pe.Next() {
        g.Projectiles.Remove(pe.Value.(*list.Element))
    }

}

func (g *Game) processShipMovement() {
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)

        // If the ship is destroyed, turn off all thrusters and everything else.
        // However, the wreckage will keep moving until it drags to a halt
        if (s.HitPoints <= 0.0) {
            // ---- disable autopilot and all other controls.
        }
        
        s.DoMovement()
    }   
}


func (g *Game) GetPreviousShip(currentShip *Ship) *Ship {
    if (currentShip == nil) {
        first := g.Ships.Back()
        return g.GetShipFromElement(first)
    } else {
        current := g.GetElementFromShip(currentShip)
        prev := current.Prev()
        return g.GetShipFromElement(prev)
    }
}

func (g *Game) GetNextShip(currentShip *Ship) *Ship {
    if (currentShip == nil) {
        first := g.Ships.Front()
        return g.GetShipFromElement(first)
    } else {
        current := g.GetElementFromShip( currentShip)
        next := current.Next()
        return g.GetShipFromElement(next)
    }
}

func (g *Game) GetShipFromElement(e *list.Element) *Ship {
    if (e == nil) { 
        return nil 
    } else { 
        return e.Value.(*Ship) 
    }
}

func (g *Game) GetElementFromShip(s *Ship) *list.Element {
    for e := g.Ships.Front(); e != nil; e = e.Next() {
        if (s == e.Value.(*Ship)) { return e }
    }
    
    panic ("Could not find provided ship in game's ship list")
}


func (g *Game) GetPreviousPlanet(currentPlanet *Planet) *Planet {
    if (currentPlanet == nil) {   
        first := g.Planets.Back()
        return g.GetPlanetFromElement(first)
    } else {
        current := g.GetElementFromPlanet(currentPlanet)
        prev := current.Prev()
        return g.GetPlanetFromElement(prev)
    }
}

func (g *Game) GetNextPlanet(currentPlanet *Planet) *Planet {
    if (currentPlanet == nil) {
        first := g.Planets.Front()
        return g.GetPlanetFromElement(first)
    } else {
        current := g.GetElementFromPlanet(currentPlanet)
        next := current.Next()
        return g.GetPlanetFromElement(next)
    }
}
func (g *Game) GetPlanetFromElement(e *list.Element) *Planet {
    if (e == nil) { 
        return nil 
    } else { 
        return e.Value.(*Planet) 
    }
}
func (g *Game) GetElementFromPlanet(p *Planet) *list.Element {
    for e := g.Planets.Front(); e != nil; e = e.Next() {
        if (p == e.Value.(*Planet)) { return e }
    }
    
    panic ("Could not find provided ship in game's ship list")
}

func (g *Game) EnqueueMessage(m *CrewMessage) {
    if (m.To == nil) {
        LogWarn("Discarding message from %s %s as it has no receipient.", m.From.FirstName, m.From.LastName)
    } else {
        LogInfo("Enqueueing message from %s %s to %s %s.", m.From.FirstName, m.From.LastName, m.To.FirstName, m.To.LastName)
        g.PendingMessages.PushBack(m)
    }
}


func (g *Game) doSetupForDevelopment() {    
    // First add the player's ship
    g.PlayerShip = NewShip() 
    g.PlayerShip.MaxForwardThrust = 1.0
    g.PlayerShip.MaxBackwardThrust = 0.4
    g.PlayerShip.MaxRotation = 18    
    g.PlayerShip.Helm.IsDirectPilot = true
    g.PlayerShip.Name = "Centauri II"
    g.PlayerShip.DesignName = "militia corvette"
    g.PlayerShip.HitSize = 1.0
    g.PlayerShip.HitPoints = 50.0
    g.PlayerShip.Weapons.PushBack(New1KgGun("Main Cannon", 300, 60))
    
    playerCrew := NewCrewMember("Victor", "Snapes", nil, CrewRoleCommander)
    playerCrew.IsPlayer = true
    g.ThePlayer.CrewMember = playerCrew
    g.ThePlayer.FireControlSelectedWeapon = g.PlayerShip.Weapons.Front().Value.(*ShipWeapon)
    g.PlayerShip.CrewMembers.PushBack(playerCrew)        
    
    g.PlayerShip.CrewMembers.PushBack(NewCrewMember("Roy", "Higgards", &AiHelmsmanBasic{}, CrewRoleHelmsman))
    
    g.Ships.PushBack(g.PlayerShip)    
    
    g.Ships.PushBack(g.createRandomAiPirateFighter())
    g.Ships.PushBack(g.createRandomAiPirateFighter())
    
   /*
    // Add any remaining AI ships
    aiShip := NewShip()
    aiShip.Point = NewPoint(30, -10)
    aiShip.ShipHeadingInDegrees = 90.0
    aiShip.MaxForwardThrust = 0.4
    aiShip.MaxBackwardThrust = 0.1
    aiShip.MaxRotation = 18 // manueverable little thing
    aiShip.CrewMembers.PushBack(&AiHelmsmanTest{})
    aiShip.Name = "Jesse's Bounty"
    aiShip.DesignName = "a pirate corvette"
    g.Ships.PushBack(aiShip)
   
    aiShip = NewShip()
    aiShip.Point = NewPoint(120, 50)
    aiShip.ShipHeadingInDegrees = 90.0
    aiShip.MaxForwardThrust = 0.4
    aiShip.MaxBackwardThrust = 0.1
    aiShip.MaxRotation = 18 // manueverable little thing
    aiShip.CrewMembers.PushBack(&AiHelmsmanTest{})
    aiShip.Name = "Carpe Diem"
    aiShip.DesignName = "a pirate corvette"
    g.Ships.PushBack(aiShip)
     */
    
    
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
    
}


func (g *Game) createRandomAiPirateFighter() *Ship {
    aiShip := NewShip()
    aiShip.Point = NewPoint(g.Rand.Float64() * 800.0 - 400.0, g.Rand.Float64() * 800.0 - 400.0).Round()
    aiShip.ShipHeadingInDegrees = 0.0
    aiShip.MaxForwardThrust = 0.4
    aiShip.MaxBackwardThrust = 0.1
    aiShip.MaxRotation = 18 // manueverable little thing
    aiShip.CrewMembers.PushBack(NewCrewMember("Unknown", "Scoundrel", &AiHelmsmanTest{}, CrewRolePilot))
    aiShip.HitSize = 0.5
    aiShip.HitPoints = 10.0
    aiShip.Name = "Pirate Fighter"
    aiShip.DesignName = "unknown design"    
    
    return aiShip
}