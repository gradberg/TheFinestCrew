/*
    Records the location and status of a given ship
*/

package main

import "container/list"

const MinimumSpeedDecay = 0.1
const SpeedDecayMultiplier = 0.02
const maxPastLocations = 10 // save a bunch, as the tactical map can just filter out how many it wants.

type Ship struct {
    Name string
    
    // There are alot of attributes that could be added, but I'll keep it to the 
    // ones that are actually implemented and have a use first.
    
    // Crew member list
    CrewMembers *list.List
    
    
    
    // Current status information    
    //X, Y float64  // Ships location on the map    
    Point Point // Location on the map
    MovementHeadingInDegrees float64 // The direction that the ship is currently moving
    SpeedInUnitsPerTick float64 // the speed of the ship going in that direction        
    ShipHeadingInDegrees float64 // The direction that the ship is FACING (and thus thrust affects)
    
    PastLocations *list.List // list of the past X locations this ship was at. This allows some maps to show 'trails' from moving ships

    Helm *ShipHelm
    //FireControl *ShipFireControl
    
    // Need a weapons-status object. Probably just have a list of weapon structs, each of which hold their own state.
    
    //Weapon *ShipWeapon // currently just a placeholder for developing the weapons functionality
    Weapons *list.List
    
    // Ship design information, like size, weapons, etc? 
    // Thrust, etc.
    
    DesignName string
    MaxForwardThrust float64 // how much the velocity can be changed each tick accelerating
    MaxBackwardThrust float64 // how much the velocity can be changed each tick deccelerating
    MaxRotation float64 // how many degrees the heading can be changed each tick
    
    // Status, like damage
    
    HitPoints float64 // 
    HitSize float64 // a 'sphere' representing how large of a ship this is, and thus how easy is it to
                    // hit with weapons. 
    
    WasHit bool // indicates the ship was hit by weapons fire in the last tick
}

func NewShip() *Ship {
    s := &Ship { }
    s.Helm = NewShipHelm(s)
    //s.FireControl = NewShipFireControl(s)
    s.CrewMembers = list.New()
    s.PastLocations = list.New()
    s.Weapons = list.New()
    return s
}

// Satisfy the ISpaceObject interface
func (s *Ship) GetPoint() Point { return s.Point }
func (s *Ship) GetCourse() float64 { return s.MovementHeadingInDegrees } // The direction that the ship is currently moving
func (s *Ship) GetSpeed() float64 { return s.SpeedInUnitsPerTick } // the speed of the ship going in that direction        
func (s *Ship) GetHeading() float64 { return s.ShipHeadingInDegrees } // The direction that the ship is FACING (and thus thrust affects)
func (s *Ship) GetName() string { return s.Name }

// This actually moves the ship based on its current momentum
func (s *Ship) DoMovement() {    
    switch s.Helm.PilotIntent {
        case PilotIntentThrust: 
            s.thrust()
        case PilotIntentSpinClockwise:
            s.rotateClockwise()            
            if s.Helm.AutoThrust { s.thrust() }    
        case PilotIntentSpinCounter:
            s.rotateCounter()            
            if s.Helm.AutoThrust { s.thrust() }    
        case PilotIntentNone:        
            if s.Helm.AutoThrust { s.thrust() }    
    }
    
    // Handle speed decay
    decay := s.SpeedInUnitsPerTick * SpeedDecayMultiplier
    if (decay < MinimumSpeedDecay) { decay = MinimumSpeedDecay }
    if (s.SpeedInUnitsPerTick < decay) { decay = s.SpeedInUnitsPerTick }
    s.SpeedInUnitsPerTick -= Round(decay,1)

    // save the current location to the list of past locations
    s.PastLocations.PushFront(s.Point)
    if (s.PastLocations.Len() > maxPastLocations) {
        s.PastLocations.Remove(s.PastLocations.Back())
    }
    
    // Move the ship
    s.Point = s.Point.AddVector(s.MovementHeadingInDegrees, s.SpeedInUnitsPerTick).Round()
    
    // Set the historical information so that it can be displayed on the GUI    
    switch s.Helm.PilotIntent {
        case PilotIntentThrust:
            s.Helm.UsedForwardThrusters = s.Helm.ThrottlePercentage > 0
            s.Helm.UsedBackwardThrusters = s.Helm.ThrottlePercentage < 0
            s.Helm.UsedClockwiseThrusters = false
            s.Helm.UsedCounterThrusters = false
        case PilotIntentSpinClockwise:
            s.Helm.UsedForwardThrusters = s.Helm.AutoThrust && s.Helm.ThrottlePercentage > 0
            s.Helm.UsedBackwardThrusters = s.Helm.AutoThrust && s.Helm.ThrottlePercentage < 0
            s.Helm.UsedClockwiseThrusters = s.Helm.ThrustersPercentage > 0
            s.Helm.UsedCounterThrusters = false
        case PilotIntentSpinCounter:
            s.Helm.UsedForwardThrusters = s.Helm.AutoThrust && s.Helm.ThrottlePercentage > 0
            s.Helm.UsedBackwardThrusters = s.Helm.AutoThrust && s.Helm.ThrottlePercentage < 0
            s.Helm.UsedClockwiseThrusters = false
            s.Helm.UsedCounterThrusters = s.Helm.ThrustersPercentage > 0
        default:
            s.Helm.UsedForwardThrusters = s.Helm.AutoThrust && s.Helm.ThrottlePercentage > 0
            s.Helm.UsedBackwardThrusters = s.Helm.AutoThrust && s.Helm.ThrottlePercentage < 0
            s.Helm.UsedClockwiseThrusters = false
            s.Helm.UsedCounterThrusters = false
    }
    
    // Clear pilot intent as the ship was moved
    s.Helm.PilotIntent = PilotIntentNone
}

func (s *Ship) rotateClockwise() {
    s.ShipHeadingInDegrees += s.MaxRotation * float64(s.Helm.ThrustersPercentage) / 100.0
    if (s.ShipHeadingInDegrees >= 360.0) { s.ShipHeadingInDegrees -= 360.0 }
}
func (s *Ship) rotateCounter() {
    s.ShipHeadingInDegrees -= s.MaxRotation * float64(s.Helm.ThrustersPercentage) / 100.0
    if (s.ShipHeadingInDegrees < 0) { s.ShipHeadingInDegrees += 360.0 }
}


func (s *Ship) thrust() {    
    var thrust float64 = 0.0
    if (s.Helm.ThrottlePercentage > 0) {
        thrust = s.MaxForwardThrust * float64(s.Helm.ThrottlePercentage) / 100.0
    } else if (s.Helm.ThrottlePercentage < 0) {
        thrust = s.MaxBackwardThrust * float64(s.Helm.ThrottlePercentage) / 100.0
    }
    s.DoAcceleration(s.ShipHeadingInDegrees, thrust)
}
func (s *Ship) DoAcceleration(angle, magnitude float64) {
    ax, ay := VectorToXy(angle, magnitude)
    mx, my := VectorToXy(s.MovementHeadingInDegrees, s.SpeedInUnitsPerTick)
    h, sp := XyToVector(ax + mx, ay + my)
    s.MovementHeadingInDegrees = Round(h    ,1)
    s.SpeedInUnitsPerTick = Round(sp,1)
}

// Each ship has one commanding officer. The particular role name can vary based on the ship's
// size, so to send messages to the right person, all of the crew members need to use this to 
// figure out who they're looking for.
func (s *Ship) GetCommandingOfficerRole() CrewRoleEnum {
    for e := s.CrewMembers.Front(); e != nil; e = e.Next() {
        cm := e.Value.(*CrewMember)
        cr := cm.CrewRole
        switch cr {
            case CrewRoleCommander: return cr
            case CrewRolePilot: return cr
            // default falls through to check the next crew member
        }
    }
    return CrewRoleError
}

// Retrieves the crew member for the given role. If one isn't found, the caller will have to deal
// with the nil
func (s *Ship) GetCrewMemberForRole(crewRole CrewRoleEnum) *CrewMember {
    for e := s.CrewMembers.Front(); e != nil; e = e.Next() {
        cm := e.Value.(*CrewMember)
        if (cm.CrewRole == crewRole) { return cm }
    }
    return nil
}

func (s *Ship) GetPreviousWeapon(sw *ShipWeapon) *ShipWeapon {
    if (sw == nil) {   
        first := s.Weapons.Back()
        return s.GetWeaponFromElement(first)
    } else {
        current := s.GetElementFromWeapon(sw)
        prev := current.Prev()
        return s.GetWeaponFromElement(prev)
    }
}

func (s *Ship) GetNextWeapon(sw *ShipWeapon) *ShipWeapon {
    if (sw == nil) {
        first := s.Weapons.Front()
        return s.GetWeaponFromElement(first)
    } else {
        current := s.GetElementFromWeapon(sw)
        next := current.Next()
        return s.GetWeaponFromElement(next)
    }
}
func (s *Ship) GetWeaponFromElement(e *list.Element) *ShipWeapon {
    if (e == nil) { 
        return nil 
    } else { 
        return e.Value.(*ShipWeapon) 
    }
}
func (s *Ship) GetElementFromWeapon(sw *ShipWeapon) *list.Element {
    for e := s.Weapons.Front(); e != nil; e = e.Next() {
        if (sw == e.Value.(*ShipWeapon)) { return e }
    }
    
    panic ("Could not find provided ShipWeapon in ship's Weapons list")
}



