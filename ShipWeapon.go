/*
    Stores information about a single weapon on the ship
*/

package main

type WeaponTypeEnum int ; const (
    WeaponTypeGun WeaponTypeEnum = iota 
    WeaponTypeLaser
    WeaponTypeMissile
)
func (e WeaponTypeEnum) ToString() string {
    switch (e) {
        case WeaponTypeGun: return "Gun"
        case WeaponTypeLaser: return "Laser"
        case WeaponTypeMissile: return "Missile"
        default: return "ERROR WEAPONTYPEENUM"
    }
}

type TargetTypeEnum int ; const (
    TargetTypeManual TargetTypeEnum = iota
    TargetTypeTarget
    TargetTypeFireAtWill
)

type ShipWeapon struct {
    EmplacementName string // what the weapon is called on the ship (front turret)
    
    // ## Status fields (reload, remaining ammunition, etc)
    CurrentCycle int // how many turns until the ship can fire again
    Ammunition int // remaining ammunitions
    
    // ## Targetting fields (No real point in having a separate fire control object
    AutoFire bool    
    TargetType TargetTypeEnum    
    FiringAngle float64
    TargetShip *Ship

    
    // ## Common Blueprint fields
    DesignName string // what the weapon DESIGN is (10 kg gun)
    WeaponType WeaponTypeEnum
    DesignDamage float64 // The maximum damage caused. For guns, this damage rating decreases as the projectile slows
    DesignCycle int // the number of turns between shots
    //DesignPrecision float64 // total width that the gun randomly fires in (thus low is better)
    // These go CLOCKWISE. So it spins clockwise from **FiringArcStart** to **FiringArcEnd**.
    FiringArcStart float64 // degrees where the firing arc start
    FiringArcEnd float64 // degrees where the firing arc ends
    // Should they take time to turn? It'd be cool : P Especially if you can watcy them turn in the panel
    //FiringAngle float64 // Current Position where the gun is sitting
    //FiringArcSpeed float64 // number of degrees the gun can spin per turn
    
    
    // ## Gun-specific fields
    DesignSpeed float64 // For Guns, the speed at which the projectile initially moves
    DesignDrag float64 // How much each turn a projectile slows down    
    DesignAmmunition int // maximum (starting) amount of rounds a projectile or missile weapon has    
    //DesignDrift float64 // how much a projectile can drift to a different heading each tick.
    
    
    // ## Laser-specific fields
    DesignDistance float64 // For lasers, they fire out to a certain distance, with a linear dropoff in damage
    // power usage?        
}

// Set of common weapon definitions
func New1KgGun(emplacementName string, firingArcStart, firingArcEnd float64, ammunition int) *ShipWeapon {
    w := &ShipWeapon { }
    w.FiringArcStart = firingArcStart
    w.FiringArcEnd = firingArcEnd
    w.DesignAmmunition = ammunition
    w.Ammunition = ammunition
    
    // ---- intelligently find the middle-point?
    if (w.IsInFiringArc(0.0)) {
        w.FiringAngle = 0.0
    } else {
        w.FiringAngle = firingArcStart
    }
    w.EmplacementName = emplacementName
    
    w.WeaponType = WeaponTypeGun
    w.DesignName = "1 kg Gun"
    w.DesignSpeed = 50.0
    w.DesignDrag = 3.0
    w.DesignDamage = 8.0     
    w.DesignCycle = 2
    
    return w
}


func New1MwLaser(emplacementName string, firingArcStart, firingArcEnd float64) *ShipWeapon {
    w := &ShipWeapon { }
    w.FiringArcStart = firingArcStart
    w.FiringArcEnd = firingArcEnd
    
    // ---- intelligently find the middle-point?
    if (w.IsInFiringArc(0.0)) {
        w.FiringAngle = 0.0
    } else {
        w.FiringAngle = firingArcStart
    }
    w.EmplacementName = emplacementName
    
    w.WeaponType = WeaponTypeLaser
    w.DesignName = "1 MW Laser"
    w.DesignDistance = 180.0
    w.DesignDamage = 3.0     
    w.DesignCycle = 1
    
    return w
}

func (s *ShipWeapon) IsInFiringArc(testAngle float64) bool {
    // Test same angle
    if (s.FiringArcStart == s.FiringArcEnd) {
        return s.FiringArcStart == testAngle
    }
    
    // if the start angle is LESS than the end angle, then they do not cross the 360->0 threshold
    if (s.FiringArcStart < s.FiringArcEnd) {
        return testAngle >= s.FiringArcStart && testAngle <= s.FiringArcEnd
    } else {
        // If it does cross the 360->0 threshold, then it is in the arc if it is not between the two values
        return !(testAngle < s.FiringArcStart && testAngle > s.FiringArcEnd)
    }    
}

//
// If the ship is set to fire-at-will or a target, this attempts to set the firing angle
func (s *ShipWeapon) SetFiringAngle(ship *Ship) {
    switch (s.TargetType) {
        case TargetTypeTarget:
            // determine the bearing to the ship
            angle, _ := s.TargetShip.Point.Subtract(ship.Point).ToVector()
            
            // Use that beaQring minus the ship's heading to determine the firing a
            bearing := AddAngles(-ship.ShipHeadingInDegrees, angle)
            
            // attempt to set that (even if it probably wont be workable)
            if (s.IsInFiringArc(bearing)) {
                s.FiringAngle = bearing
            } else {
                startAngle, _ := GetShortestTurn(bearing, s.FiringArcStart)
                endAngle, _ := GetShortestTurn(bearing, s.FiringArcEnd)
                
                if (startAngle < endAngle) {
                    s.FiringAngle = s.FiringArcStart
                } else {
                    s.FiringAngle = s.FiringArcEnd
                }
            }
        
        case TargetTypeFireAtWill:
            // ---- unimplemented yet
    
        default:
            // do nothing
    }
}

