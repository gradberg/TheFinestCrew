/*
    Organizes game-setup code
*/


package main

type PlayableShipEnum int; const (
    PlayableShipMilitiaCorvette PlayableShipEnum = iota
    PlayableShipPirateFighter
)
func (e PlayableShipEnum) ToString() string {
    switch e {
        case PlayableShipMilitiaCorvette: return "Militia Corvette"
        case PlayableShipPirateFighter: return "Pirate Fighter"
        default: return "ERROR PLAYABLESHIPENUM"
    }
}

type GameSetup struct {
    CrewRole CrewRoleEnum
    Ship PlayableShipEnum
    PlayerName string
    //ShipName string
}

func (g *Game) createPlayerShip() {
    if (g.GameSetup.Ship == PlayableShipMilitiaCorvette) {
        g.PlayerShip = createMilitiaCorvette("Centauri II")    
        
        if (g.GameSetup.CrewRole == CrewRoleHelmsman) {
            playerCrew := NewCrewMember("Victor", "Snapes", nil, CrewRoleHelmsman)
            playerCrew.IsPlayer = true
            g.ThePlayer.CrewMember = playerCrew
            g.ThePlayer.FireControlSelectedWeapon = g.PlayerShip.Weapons.Front().Value.(*ShipWeapon)
            g.PlayerShip.CrewMembers.PushBack(playerCrew)            
            //g.PlayerShip.CrewMembers.PushBack(NewCrewMember("Roy", "Higgards", &AiHelmsmanBasic{}, CrewRoleHelmsman))
        } else if (g.GameSetup.CrewRole == CrewRoleCommander) {
            playerCrew := NewCrewMember("Victor", "Snapes", nil, CrewRoleCommander)
            playerCrew.IsPlayer = true
            g.ThePlayer.CrewMember = playerCrew
            g.ThePlayer.FireControlSelectedWeapon = g.PlayerShip.Weapons.Front().Value.(*ShipWeapon)
            g.PlayerShip.CrewMembers.PushBack(playerCrew)            
            g.PlayerShip.CrewMembers.PushBack(NewCrewMember("Roy", "Higgards", &AiHelmsmanBasic{}, CrewRoleHelmsman))
        } else {
            panic("Incorrect startup (unknown role for militia corvette")
        }
    } else if (g.GameSetup.Ship == PlayableShipPirateFighter) {
        g.PlayerShip = createPirateFighter("Black III")
        
        if (g.GameSetup.CrewRole == CrewRolePilot) {
            playerCrew := NewCrewMember("Victor", "Snapes", nil, CrewRolePilot)
            playerCrew.IsPlayer = true
            g.ThePlayer.CrewMember = playerCrew
            g.ThePlayer.FireControlSelectedWeapon = g.PlayerShip.Weapons.Front().Value.(*ShipWeapon)
            g.PlayerShip.CrewMembers.PushBack(playerCrew)            
        } else {
            panic("Incorrect startup (unknown role for pirate fighter")
        }
    } else {
        panic("Incorrect startup (unknown ship)")
    }
 
}

func createMilitiaCorvette(name string) *Ship {
    p := NewShip() 
    p.MaxForwardThrust = 1.0
    p.MaxBackwardThrust = 0.4
    p.MaxRotation = 18    
    p.Helm.IsDirectPilot = true
    p.Name = name
    p.DesignName = "militia corvette"
    p.HitSize = 3.7
    p.HitPoints = 50.0
    p.Weapons.PushBack(New1KgGun("Main Cannon", 330, 30, 40))
    p.Weapons.PushBack(New1MwLaser("Fore Laser", 300, 60))
    p.Weapons.PushBack(New1MwLaser("Port Laser", 180, 300))
    p.Weapons.PushBack(New1MwLaser("Strbd. Laser", 60, 180))
    return p 
}



func createPirateFighter(name string) *Ship {
    aiShip := NewShip()
    aiShip.MaxForwardThrust = 0.4
    aiShip.MaxBackwardThrust = 0.1
    aiShip.MaxRotation = 18 // manueverable little thing
    aiShip.HitSize = 1.8
    aiShip.HitPoints = 10.0
    aiShip.DesignName = "pirate fighter"    
    aiShip.Name = name
    aiShip.Weapons.PushBack(New1MwLaser("Main Weapon", 300, 60))    
    return aiShip
}


