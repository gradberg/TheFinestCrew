/*
    The Game class represents the game being played, and contains the main
    game loop and all relevant state for playing the game.
    
    Any display-specific knowledge should be kept in the GameDisplayAdapter.
*/

package main

import "container/list"
import "math/rand"
import "time"

type GameStatusEnum int; const (
    GameStatusTitleScreen GameStatusEnum = iota
    GameStatusPickShip
    GameStatusPickRole
    GameStatusNamePlayer
    GameStatusInstructions
    
    // ---- long term, most of these statuses would be pulled in to a game mode
    //      object which would manage all state related to game mode itself. That means
    //      the interface would need a method for the GameDisplayAdapter and
    //      GameInputAdapter to find out if various panels should be displayed based
    //      on the mode state.
    
    
    GameStatusDeathScreen
    
    GameStatusNextRound
    
    GameStatusPlaying
)

// Represents the game being played. 
type Game struct {
    GameStatus GameStatusEnum
    
    
    // ---- these will be extracted out to a GameMode object in the future,
    //      so that the map/universe can remain separate from the mechanics of
    //      how different game modes progress
    round int // used to keep track of the round. Probably should be
    // stored somewhere in the GameModeSurvival.go file. Needs a IGameMode interface?
    // that way game modes can be modular and have a good way to save all of their
    // relevant data throughout the mode?
    kills int // number of ships the player destroyed.
    waitATurn bool // used to pause a single turn after the round is over before moving on.

    GameSetup GameSetup
    
    tick int    
    
    ThePlayer *Player
    PlayerShip *Ship    
    
    Ships *list.List
    
    message string
    
    Planets *list.List
    
    Rand *rand.Rand
    
    Projectiles *list.List // Projectiles that are still flying around
    LaserProjectiles *list.List // List of lasers fired last tick (which is displayed on the map)
    
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
        LaserProjectiles: list.New(),
    }    
    g.ThePlayer = NewPlayer(g)
    
    return g
}

// Removes items from all lists, but leaves the game mode and such intact
// (Used for transitioning between rounds)
func (g *Game) ClearEphemeralState() {
    // ---- should all objects have their own 'ClearEphemeralState' method, so that
    //      they can be responsible for it themselves?
    
    g.ThePlayer.ClearEphemeralState(g)

    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)
        s.ClearEphemeralState()
    }

    g.Planets.Init()
    g.Ships.Init()
    g.PendingMessages.Init()
    g.Projectiles.Init()
    g.LaserProjectiles.Init()
}

// Blocking call which contains the main game loop, and exits when the user has chosen to quit.
func (g *Game) Run() {
    AsyncInput_Init()
    defer AsyncInput_Close()
    
    // Development setup
    //g.GameStatus = GameStatusPlaying
    //g.doSetupForDevelopment()

    for {
        // ---- check for change in rounds?
        if (g.GameStatus == GameStatusPlaying && g.checkForLost())  { 
            g.GameStatus = GameStatusDeathScreen
        }
        
        if (g.GameStatus == GameStatusPlaying) {
            if (g.checkForNextRound()) {
                //g.setupForNextRound()
                g.GameStatus = GameStatusNextRound
            }
        }
    
        g.Display()
        
        g.message = ""
        var ir *InputResult = g.DoInput()
        if ir.Exit { return }
        
        // This loop is used for the entire game lifecycle, so only
        // process ticks if the game is actually being played.
        if (g.GameStatus == GameStatusPlaying) {        
            g.processTick(ir)
        }
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
        g.processLasers()        
        g.processProjectiles()    
        
        // Last process ship movement, as they are essentially the slowest objects.
        g.processShipMovement()        
        
        // do this HERE, because 
        g.processMessages()
        
        g.tick++
    }       
}

func (g *Game) processAiCharacters() {
    // go through each ship's crewmembers. 
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship) 
        
        if (s.IsDestroyed() == false) {        
            // Go through each crew member and process their turn
            for ce := s.CrewMembers.Front(); ce != nil; ce = ce.Next() {
                c := ce.Value.(*CrewMember)
                
                if (c.IsPlayer) { continue }
               
                // ---- use the ticks number to control how many turns this crew member loses?                
                // ---- convert these parameters to a struct so that anything else in the future
                //      is easy to add (and passing them around internally in the ai structure
                //      is easy to write too)
                c.Ai.DoAction(s,g,c)
            }
            
            // do other per-ship tasks
            s.Helm.DoAutoPilot()
        }
        
        s.WasHit = false
    }   
    
}

func (g *Game) processMessages() {
    // go through each message and enqueue it in its receipients lists        
    for me := g.PendingMessages.Front(); me != nil; me = me.Next() {
        m := me.Value.(*CrewMessage)
        m.TickReceived = g.tick

        // This is sent to a particular crewmember
        if (m.To != nil) {        
            m.To.ReceivedMessages.PushBack(m)
            
            // ---- Loop through the remaining crew members and have them overhear it.
        

            // also add this to the sending player's messages
            if (m.From != nil) {
                m.From.ReceivedMessages.PushBack(m)
            }        
        }
        
        // this is meant to be 'heard' by everyone on the ship
        if (m.ToShip != nil) {
            for ce := m.ToShip.CrewMembers.Front(); ce != nil; ce = ce.Next() {
                c := ce.Value.(*CrewMember)
                c.ReceivedMessages.PushBack(m)
            }
        }
    }
    g.PendingMessages.Init() // clear the list
}

func (g *Game) processWeapons() {
    // clear the list of lasers fired last turn
    g.LaserProjectiles.Init()
    
    // Loop over each ship and each weapon, creating fired lasers and projectiles as appropriate
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)        
        if s.IsDestroyed() { continue }
        
        for we := s.Weapons.Front(); we != nil; we = we.Next() {
            w := we.Value.(*ShipWeapon)
            g.processShipWeapon(s, w)
        }        
    }   
}

func (g *Game) processShipWeapon(s *Ship, w *ShipWeapon) {
    // if its not ready to fire, just cycle it.
    if (w.CurrentCycle > 0) {
        // Only cycle if there is ammo, otherwise it doesn't make logical sense
        // that it is loaded with no ammo
        if (w.Ammunition > 0 || w.WeaponType == WeaponTypeLaser) {
            w.CurrentCycle --
        }
        return
    }
    
    // If the weapon is not set to fire... don't fire!
    if (!w.AutoFire) { return }    
    
    w.SetFiringAngle(s)
    shotAngle := AddAngles(w.FiringAngle, s.ShipHeadingInDegrees)            
    switch (w.WeaponType) {
        case WeaponTypeLaser:
            // This just SPAWNS the lasers. They impact just like projectiles. This is to prevent
            // a ship getting hit with a laser fired by an 'early' ship from being destroyed and unable
            // to fire its lasers. It makes it fair.                   
            p := &Projectile {
                Point: s.Point,
                Heading: shotAngle,                        
                OriginShip: s,
                Speed: w.DesignDistance,
                DesignDamage: w.DesignDamage,
            }
            g.LaserProjectiles.PushBack(p)                    
            
            LogCalc("Laser fired from ship. sHeading [%f], pHeading [%f]", s.ShipHeadingInDegrees, p.Heading)
            
            // weapon must now reload
            w.CurrentCycle = w.DesignCycle
        
        case WeaponTypeGun:
            // if the weapon is projectile based and out of ammunition, don't fire
            if (w.Ammunition == 0) { return }
            
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
            LogCalc("Projectile fired from ship. sHeading [%f], pHeading [%f]", s.ShipHeadingInDegrees, p.Heading)
            
            // Set the weapon back to its reload setting
            w.CurrentCycle = w.DesignCycle
            w.Ammunition --            
    } 
}

func (g *Game) processLasers() {
    // ---- todo (after adding in the tactical map logic to DISPLAY lasers)
    
    for pe := g.LaserProjectiles.Front(); pe != nil; pe = pe.Next() {
        p := pe.Value.(*Projectile)
        
        // Get a list of any ships it would impact
        startPoint := p.Point
        endPoint := p.GetFuturePoint()
        impactedShips := make([]*Ship, 0, 10)
        for se := g.Ships.Front(); se != nil; se = se.Next() {
            s := se.Value.(*Ship)
            if (p.OriginShip == s) { continue }
            
            impactDistance := DistanceFromLineSegment(s.Point, startPoint, endPoint)
            if (impactDistance <= s.HitSize) {
                impactedShips = append(impactedShips, s)
                LogCalc("Laser *MIGHT HIT* ship %s at %f,%f by laser going from %f,%f to %f,%f",
                    s.Name, s.Point.X(), s.Point.Y(), 
                    startPoint.X(), startPoint.Y(), 
                    endPoint.X(), endPoint.Y(),
                )
            }
        }
        
        // if it hit nothing, just continue
        if len(impactedShips) == 0 { continue }
        
        // go through that list of ships, and impact the one closest to
        // where the laser started (since in 99% of situations that is
        // the one it would impact, and the other 1% would be too ambigious
        // for anyone to care)
        closestShip := impactedShips[0]
        closestAngle, closestDistance := closestShip.Point.Subtract(startPoint).ToVector()
        for i := 1; i < len(impactedShips); i++ {            
            s := impactedShips[i]
            angle, distance := s.Point.Subtract(startPoint).ToVector()
            if (distance < closestDistance) {
                closestShip = s
                closestDistance = distance
                closestAngle = angle
            } 
        }        
        
        // Set the laser's end-point to match the ship's loaction 
        // and damage the ship
        s := closestShip
        damage := p.DesignDamage * (p.Speed - closestDistance) / p.Speed
        s.HitPoints = Round(s.HitPoints - damage, 1)        
        LogCalc("Laser *HIT* Ship %s at %f, %f for damage %f",
            s.Name, s.Point.X(), s.Point.Y(), damage,
        )                    
        s.WasHit = true
        
        p.Speed = closestDistance
        p.Heading = closestAngle

        // Enqueue message for the ship being damaged
        message := ""
        if (s.HitPoints / s.DesignHitPoints) > 0.333 {
            message = PickRandomString(g,
                "The ship was hit by laser fire.", 
                "You hear a faint buzzing sound as the ship encounters laser fire.",
                "The lights flare brightly as the ship is hit by a laser.",
            )
        } else {
            message = PickRandomString(g,
                "Gaping holes are burned into the hull by laser fire.", 
                "You hear a loud buzzing sound in the ship's electrical system as it is stuck by a laser.",
                "Several light fixtures flare brightly and burn out as the ship is hit by a laser.",
                "Computer consoles flash static as the ship is hit by laser fire.",
            )
        }                
        g.PendingMessages.PushBack(NewShipStatusMessage(message, s))
        
    }
}

func (g *Game) processProjectiles() {
    // Do projectile hit detection
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
                
                LogCalc("Projectile  *HIT* Ship %s at %f, %f by projectile at %f, %f, pHeading %f, pSpeed %f, distance from path %f, damage %f",
                    s.Name, s.Point.X(), s.Point.Y(), p.Point.X(), p.Point.Y(), p.Heading, p.Speed, distance, damage,
                )                    
                
                // Transfer impact to the ship's momentum IF IT IS A PROJECTILE OR MISSILE
                s.DoAcceleration(p.Heading, damage / 5.0)
                s.WasHit = true
                didImpact = true

                // Enqueue message for the ship being damaged
                message := ""
                if (s.HitPoints / s.DesignHitPoints) > 0.333 {
                    message = PickRandomString(g,
                        "The ship was hit by a projectile.", 
                        "The ship shakes as it is hit by a projectile.",
                        "The lights flicker as a projectile strikes the ship.",
                    )
                } else {
                    message = PickRandomString(g,
                        "The ship shudders violently, being hit by a projectile.", 
                        "A projectile strikes the ship, breaching several compartments.",
                        "The lights cut out for a moment after a projectile strikes the ship.",
                        "Sparks fly from computer consoles as a projectile slams into the ship.",
                    )
                }                
                g.PendingMessages.PushBack(NewShipStatusMessage(message, s))
                
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

        // The ship does not try to change its momentum if it is destroyed
        if (s.IsDestroyed() == false) {
            s.DoThrust()
        }
        
        // However, destroyed ships will keep on moving until the wreckage
        // drags to a halt
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