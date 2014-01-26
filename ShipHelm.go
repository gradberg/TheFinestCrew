/*
    Holds all status for the Ship's Helm. None of these directly change where the ship is, but when a tick
    occurs, the ship uses this information to determine how its position can change.
*/

package main
import "math"
    
type AutoPilotModeEnum int; const (
    AutoPilotModeOff AutoPilotModeEnum = iota
    AutoPilotModeFullStop
    AutoPilotModeCourse
    AutoPilotModeTarget
)

type PilotIntentEnum int; const ( 
    PilotIntentNone PilotIntentEnum = iota
    PilotIntentThrust
    PilotIntentSpinClockwise
    PilotIntentSpinCounter
)

type ShipHelm struct {
    Ship *Ship
    
    IsDirectPilot bool // If the ship is running on auto pilot or manual        
    ThrottlePercentage int // the throttle, from 100 to -100
    AutoThrust bool // Indicates that every tick the ship engages thrust, regardless of whatever other action.
    ThrustersPercentage int // Spin Thrusters throttle, from 000 to 100  
        
    AutoPilotMode AutoPilotModeEnum // Which setting the auto-pilot is in    
    AutoPilotPower int // 0-3 Power setting, which indicates the maximum spin and thrust autopilot can use
    
    AutoPilotDesiredCourse float64 // desired heading for autopilot-course
    AutoPilotDesiredSpeed float64 // desired speed for autopilot-course
    
    AutoPilotDesiredTarget ISpaceObject // chosen target by the player
    
    // Movement Intent (so all physical state changes occur simulatenously) 
    PilotIntent PilotIntentEnum
    
    //status indicators as to what the ship DID do last turn (thrusted forwards, backwards, spun)
    UsedForwardThrusters bool
    UsedBackwardThrusters bool
    UsedClockwiseThrusters bool
    UsedCounterThrusters bool

}

func NewShipHelm(s *Ship) *ShipHelm {
    return &ShipHelm {
        Ship: s,
        ThrottlePercentage: 100,
        ThrustersPercentage: 100,
        AutoPilotPower: 3,
    }
}


func (h *ShipHelm) ThrottleUp() {
    h.ThrottlePercentage += 10
    if h.ThrottlePercentage > 100 { h.ThrottlePercentage = 100 }
}
func (h *ShipHelm) ThrottleDown() {
    h.ThrottlePercentage -= 10
    if h.ThrottlePercentage < -100 { h.ThrottlePercentage = -100 }
}

func (h *ShipHelm) ThrustersUp() {
    h.ThrustersPercentage += 10
    if h.ThrustersPercentage > 100 { h.ThrustersPercentage = 100 }
}
func (h *ShipHelm) ThrustersDown() {
    h.ThrustersPercentage -= 10
    if h.ThrustersPercentage < 0  { h.ThrustersPercentage = 0 }
}

func (h *ShipHelm) AutoPilotPowerUp() {
    h.AutoPilotPower++
    if (h.AutoPilotPower > 3) { h.AutoPilotPower = 3 }
}
func (h *ShipHelm) AutoPilotPowerDown() {
    h.AutoPilotPower--
    if (h.AutoPilotPower < 0) { h.AutoPilotPower = 0 }
}

// Turns off autopilot and remoes any settings
func (h *ShipHelm) ClearAutoPilot() {
    h.AutoPilotMode = AutoPilotModeOff
}


// If the auto-pilot is on, this sets the throttles and turns/thrusts as necessary to achieve the desired goal
func (h *ShipHelm) DoAutoPilot() {
    // ---- At some point autopilot will probably have to be stateful, and just clear out its state when it gets changed
    
    switch h.AutoPilotMode {
        case AutoPilotModeFullStop: h.doAutoPilotFullStop()
        case AutoPilotModeCourse: h.doAutoPilotCourse()
        case AutoPilotModeTarget: h.doAutoPilotTarget()
        default: // do nothing
    }
}

func (h *ShipHelm) doAutoPilotFullStop() {
    // If the ship is stopped, turn all thrusters to max. And do nothing else.
    if (h.Ship.SpeedInUnitsPerTick == 0) {
        h.ThrustersPercentage = 100
        h.ThrottlePercentage = 100        
        h.AutoThrust = false
    
        // Full Stop mode stays on even if the ship is stopped, as it can be used to stop the ship when hit by weapons.
        return
    }
    
    // Otherwise, if the ship is moving, determine if the ship should spin before thrusting
    // This is going to be basic logic to start? It will always turn to the exact angle before 
    desiredAngle := GetOppositeDegrees(h.Ship.MovementHeadingInDegrees)
    h.TurnTowardsAngle(desiredAngle, h.maxAutoPilotPower())
    
    // The ship has turned if it needs to. If it is at the desired angle, set the
    // Throttle as necessary and turn on auto-thrust. 
    // NOTE this does mean the ship will thrust on the tick it gets to the desired angle
    if math.Abs(float64(int(h.Ship.ShipHeadingInDegrees) - int(desiredAngle))) <= 1.0 {
        if  h.Ship.SpeedInUnitsPerTick > h.Ship.MaxForwardThrust {
            h.ThrottlePercentage = 100
        } else {
            h.ThrottlePercentage = int(h.Ship.SpeedInUnitsPerTick * 100.0 / h.Ship.MaxForwardThrust)
        }
        if h.maxAutoPilotPower() < h.ThrottlePercentage { h.ThrottlePercentage = h.maxAutoPilotPower() }
        h.AutoThrust = true
    } else {
        h.AutoThrust = false
    }
}

func (h *ShipHelm) doAutoPilotCourse() {
    // ---- todo
    h.TurnTowardsAngle(h.AutoPilotDesiredCourse, h.maxAutoPilotPower())
    
    // ---- replace, as this implementation rather sucks
    // if the speed it below desired, fire the engines 100%. If it above desired, turn them off.
    if (h.AutoPilotDesiredSpeed > h.Ship.SpeedInUnitsPerTick) {
        h.ThrottlePercentage = 100
        h.AutoThrust = true
    } else {
        h.AutoThrust = false
    }
}

func (h *ShipHelm) doAutoPilotTarget() {

    // This should adapt its behavior based on how far away the target is.
    // if it is very far, then just aim at the ship and gun the engines.
    // if it is close, be smarter about that (whatever that means???)
    // To define 'close' and 'far', base that off of how fast the ship
    // is going towards the enemy versus how close far it is from the target
    sp := h.Ship.Point
    tp := h.AutoPilotDesiredTarget.GetPoint()

    modifierAngleInDegrees, distance := 0.0, 0.0 // if nothing selected, the player ship displays relative to 0.0
    modifierAngleInDegrees, distance = tp.Subtract(sp).ToVector()
    
    distance = distance // GO!    
    
    h.TurnTowardsAngle(modifierAngleInDegrees, h.maxAutoPilotPower())
    h.AutoThrust = true
    
    
    // ---- Lets start with the very naive solution of turn towards the target, and keep
    // the engines running the whole time.
    
    
}

// Sets the helm to turn towards the given angle in the shortest direction, which the
// correect amount of thrust not exceeding the provided maximum.
func (h *ShipHelm) TurnTowardsAngle(desiredAngle float64, maxThrusters int) {
    // If the angles already match (as close as they can) then do not bother turning
    // and instead shut off the turn flags
    if int(h.Ship.ShipHeadingInDegrees) == int(desiredAngle) {
        h.ThrustersPercentage = 100
        h.PilotIntent = PilotIntentNone
        return
    }

    // Turn the ship to adjust
    degreesToTurn, isClockwise := GetShortestTurn(h.Ship.ShipHeadingInDegrees, desiredAngle)
    
    // Determine the Spin Thruster amount
    if degreesToTurn > h.Ship.MaxRotation {
        h.ThrustersPercentage = 100
    } else {
        h.ThrustersPercentage = int(degreesToTurn * 100.0 / h.Ship.MaxRotation)
    }
        
    // Cap speed based on the provided maximum
    if (maxThrusters < h.ThrustersPercentage) { h.ThrustersPercentage = maxThrusters }
        
    if isClockwise {
        h.PilotIntent = PilotIntentSpinClockwise
    } else {
        h.PilotIntent = PilotIntentSpinCounter
    }
}

func (h *ShipHelm) maxAutoPilotPower() int {
    switch h.AutoPilotPower {
        case 0: return 13
        case 1: return 25
        case 2: return 50
        default: return 100
    }
}
