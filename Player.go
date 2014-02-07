/*
    Stores all information relevant to the individual player's state,
    including player stats and display information.
*/

package main

type Player struct {
    Game *Game

    CrewMember *CrewMember // the crew member representing the player

    RunningRealTime bool // indicates the player is letting the game run
    RealTimeTicksPerSecond float64 // has many ticks per second run in real time mode

    TacticalMapScale float64 // view scale of the tactical map
    TacticalMapFullScreen bool
    
    TacticalAnalysisSelection *Ship // pointer to whatever ship the tactical analysis panel is looking at.
                           // if nil, then nothing is selected.
                           
    HelmControlSettingCourse bool // indicates the player is on the set-course screen
    HelmControlCourseSetter *CourseSetter
    
    HelmControlSettingTarget bool // indicates the player is currently setting a target to move to
    HelmControlTarget *Targeter
    
    PersonnelStatus PersonnelStatusEnum // indicates the personnel panel is in a certain visual state
    PersonnelHelmsmanTarget *Targeter // used to choose a target for the helmsman
    PersonnelHelmsmanCourseSetter *CourseSetter // used to choose course for the helmsman
 
    FireControlSelectedWeapon *ShipWeapon
 
    FireControlSettingAim bool
    FireControlAngle float64 // degrees of manual aim for the gun
    //FireControlCourseSetter *CourseSetter // used for aiming the gun?
    
    FireControlSettingTarget bool
    FireControlTarget *Targeter
}

func NewPlayer(g *Game) *Player {
    return &Player {
        Game: g,
        TacticalMapScale: 100, // How much distance around the map the player sees.
        RealTimeTicksPerSecond: 1.0,
        HelmControlTarget: NewTargeter(g),
        HelmControlCourseSetter: NewCourseSetter(g),
        PersonnelHelmsmanTarget: NewTargeter(g),
        PersonnelHelmsmanCourseSetter: NewCourseSetter(g),
        FireControlTarget: NewTargeter(g),
    }
}

func (p *Player) ClearEphemeralState(g *Game) {
    p.TacticalMapScale = 100
    
    // Targeters must be cleared out so that they do not reference nonexistant ships.
    p.HelmControlTarget = NewTargeter(g)
    p.HelmControlSettingTarget = false
    p.PersonnelHelmsmanTarget = NewTargeter(g)    
    p.FireControlTarget = NewTargeter(g)
    p.FireControlSettingTarget = false
    
    p.FireControlSettingAim = false
    
    // Course setters might as well be cleared out to present a clean slate
    p.HelmControlCourseSetter = NewCourseSetter(g)
    p.HelmControlSettingCourse = false
    p.PersonnelHelmsmanCourseSetter = NewCourseSetter(g)
    p.PersonnelStatus = PersonnelStatusNormal
    p.TacticalAnalysisSelection = nil
    
    p.RunningRealTime = false
    p.TacticalMapFullScreen = false
}

func (p *Player) DecreaseTacticalMapScale() {
    if p.TacticalMapScale >= 100000.0 { return }
    p.TacticalMapScale *= 2.0
}
func (p *Player) IncreaseTacticalMapScale() {    
    if p.TacticalMapScale <= 0.1 { return }
    p.TacticalMapScale /= 2.0
}
func (p *Player) ToggleTacticalMapFullScreen() {
    if p.TacticalMapFullScreen { 
        p.TacticalMapFullScreen = false
    } else {
        p.TacticalMapFullScreen = true
    }
}
func (p *Player) IncreaseRealTimeTicksPerSecond() {
    if p.RealTimeTicksPerSecond >= 32.0 { return }
    p.RealTimeTicksPerSecond *= 2.0
}
func (p *Player) DecreaseRealTimeTicksPerSecond() {
    if p.RealTimeTicksPerSecond <= 0.2 { return }
    p.RealTimeTicksPerSecond /= 2.0
}

func (p *Player) FireControlAngleLargeCounter(s *ShipWeapon) {
    p.FireControlAngle = Round(AddAngles(p.FireControlAngle, -1.0),1)
    if (s.IsInFiringArc(p.FireControlAngle) == true) { return } 
    p.FireControlAngle = s.FiringArcStart
}

func (p *Player) FireControlAngleSmallCounter(s *ShipWeapon) {
    p.FireControlAngle = Round(AddAngles(p.FireControlAngle, -0.1),1)
    if (s.IsInFiringArc(p.FireControlAngle) == true) { return } 
    p.FireControlAngle = s.FiringArcStart
}

func (p *Player) FireControlAngleLargeClockwise(s *ShipWeapon) {
    p.FireControlAngle = Round(AddAngles(p.FireControlAngle, 1.0),1)
    if (s.IsInFiringArc(p.FireControlAngle) == true) { return } 
    p.FireControlAngle = s.FiringArcEnd
}

func (p *Player) FireControlAngleSmallClockwise(s *ShipWeapon) {
    p.FireControlAngle = Round(AddAngles(p.FireControlAngle, 0.1),1)
    if (s.IsInFiringArc(p.FireControlAngle) == true) { return } 
    p.FireControlAngle = s.FiringArcEnd
}

func (p *Player) FireControlSelectedWeaponPrevious() {
    p.FireControlSelectedWeapon = p.Game.PlayerShip.GetPreviousWeapon(p.FireControlSelectedWeapon)
    if (p.FireControlSelectedWeapon == nil) {
        p.FireControlSelectedWeapon = p.Game.PlayerShip.GetPreviousWeapon(p.FireControlSelectedWeapon)
    }
}
func (p *Player) FireControlSelectedWeaponNext() {
    p.FireControlSelectedWeapon = p.Game.PlayerShip.GetNextWeapon(p.FireControlSelectedWeapon)
    if (p.FireControlSelectedWeapon == nil) {
        p.FireControlSelectedWeapon = p.Game.PlayerShip.GetNextWeapon(p.FireControlSelectedWeapon)
    }
}
         