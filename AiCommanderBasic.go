/*
    A commander AI that sends appropriate orders to other crew members and manages
    the weapons.
*/

package main


type AiCommanderBasic struct {        
    primaryTargetShip *Ship  // current enemy target
}

func (ai *AiCommanderBasic) ClearEphemeralState() {
    ai.primaryTargetShip = nil
}

func (ai *AiCommanderBasic) DoAction(s *Ship, g *Game, cm *CrewMember) int {
    result := 0
    
    // if the current enemy is destroyed, change it and tell the helmsman
    // to go to it.    
    result = ai.doChangePrimaryTarget(s, g, cm)
    if (result > 0) { return result }
    
    // if there is no primarytarget at this point, do nothing.
    if (ai.primaryTargetShip == nil) { return 0 }
    
    result = ai.doAimWeaponsAtPrimaryTarget(s,g,cm)
    if (result > 0) { return result }
    
    // if not, loop through the remaining enemy ships, see if any are in
    // the firing arc of a weapon, and if so, set it to fire.
    result = ai.doAimWeaponsAtRemainingEnemies(s,g,cm)
    if (result > 0) { return result }

    // ---- what about evasive action if the ship is taking too much damage
    
    return 0
}

// Determines if the AI needs to choose a new target enemy. If it does
// it tells the Helmsman to navigate to that ship.
func (ai *AiCommanderBasic) doChangePrimaryTarget(s *Ship, g *Game, cm *CrewMember) int {
    // If a valid target is still selected, stick with it.
    // ---- this could be smarter about switching to a closer one, etc.
    if !(ai.primaryTargetShip == nil || ai.primaryTargetShip.IsDestroyed()) { return 0 }
    
    var closestShip *Ship
    closestDistance := float64(MaxInt)
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        ts := se.Value.(*Ship)
        
        if (ts == s) { continue }
        if (ts.IsDestroyed()) { continue }
        // ---- verify it is an enemy ship
        
        distance := ts.GetPoint().DistanceFrom(s.GetPoint())
        
        if (distance < closestDistance) {
            closestShip = ts
            closestDistance = distance
        }
    }
    
    if (closestShip == nil) {
        return 0
    } else {
        ai.primaryTargetShip = closestShip
        
        LogAi(cm, "Chose new ship as target: [%s]", closestShip.GetName())
        
        //  send message to helmsman to target ship
        helmsman := s.GetCrewMemberForRole(CrewRoleHelmsman)
        g.EnqueueMessage(NewMessageSetDestination(cm, helmsman, closestShip, SettingTargetTypeShip))
        
        return 1
    }
}

func (ai *AiCommanderBasic) doAimWeaponsAtPrimaryTarget(s *Ship, g *Game, cm *CrewMember) int {
    usableWeapons := ai.findUsableWeapons(s, ai.primaryTargetShip)
    
    // first loop over any weapons that are already targeting the primary target
    // and make sure they are set to auto-fire. This does NOT take any ticks.
    for _, w := range usableWeapons {
        if (w.TargetType == TargetTypeTarget && w.TargetShip == ai.primaryTargetShip) {
            w.AutoFire = true
        }
    }
    
    // Next, loop over them again and if one is NOT set to it as the target, set it.
    // This DOES take a tick
    for _, w := range usableWeapons {
        if (w.TargetType == TargetTypeTarget && w.TargetShip == ai.primaryTargetShip) {
            continue
        }
        
        LogAi(cm, "Setting weapon %s to target ship %s", w.EmplacementName, ai.primaryTargetShip.Name)
        
        w.TargetType = TargetTypeTarget
        w.TargetShip = ai.primaryTargetShip
        return 1
    }    

    return 0
}

func (ai *AiCommanderBasic) doAimWeaponsAtRemainingEnemies(s *Ship, g *Game, cm *CrewMember) int {
    // Loop through each weapon from the top down, and make sure it is targeting the
    // closest enemy
    for we := s.Weapons.Front(); we != nil; we = we.Next() {
        w := we.Value.(*ShipWeapon)

        var closestShip *Ship
        closestDistance := float64(MaxInt)
        for se := g.Ships.Front(); se != nil; se = se.Next() {
            ts := se.Value.(*Ship)
            
            if (ts == s) { continue }
            if (ts.IsDestroyed()) { continue }
            // ---- verify it is an enemy ship
            if (ai.isWeaponUsableAgainstShip(w, s, ts) == false) { continue }
            
            distance := ts.GetPoint().DistanceFrom(s.GetPoint())            
            if (distance < closestDistance) {
                closestShip = ts
                closestDistance = distance
            }
        }
        
        // if there is no valid target, just turn off the weapon here. Note this means 
        // the AI does not turn of ALL invalid targets in one turn like they are technically
        // allowed to, but it shouldn't be a big deal at all.
        if (closestShip == nil) {
            w.AutoFire = false
            // this does not take a tick, so do not return here.
        } else if !(w.TargetType == TargetTypeTarget && closestShip == w.TargetShip) {
            w.AutoFire = true
            w.TargetType = TargetTypeTarget
            w.TargetShip = closestShip
            return 1
        } else {
            // do nothing. It wasn't cleared, it wasn't set. Move on to the next one.
        }
    }
    
    return 0
}

func (ai *AiCommanderBasic) findUsableWeapons(playerShip, targetShip *Ship) []*ShipWeapon {
    results := []*ShipWeapon { }
    
    for we := playerShip.Weapons.Front(); we != nil; we = we.Next() {
        w := we.Value.(*ShipWeapon)
        
        if ai.isWeaponUsableAgainstShip(w, playerShip, targetShip) {
            results = append(results, w)
        }
    }    
    
    return results
}

func (ai *AiCommanderBasic) isWeaponUsableAgainstShip(w *ShipWeapon, playerShip, targetShip *Ship) bool {     
    bearing, distance := targetShip.Point.Subtract(playerShip.Point).ToVector()
       
    // Skip if its out of cycle (as it'll be better to use a different gun
    if (w.CurrentCycle > 0) { return false }
    
    // skip if out of ammunition or out of range
    if (w.WeaponType == WeaponTypeGun) {
        if (w.Ammunition == 0) { return false }
        // Projectile weapons are only "within range" if they are within 2 turns of hitting 
        // the enemy
        if (distance > w.DesignSpeed * 2) { return false }
    } else if (w.WeaponType == WeaponTypeLaser) {
        if (distance > w.DesignDistance) { return false }
    } else {
        LogWarn("AiCommanderBasic Weapon Type not defined!")
    }
    
    // Check that it is in the firing angle
    weaponAngle := AddAngles(bearing, - playerShip.ShipHeadingInDegrees)
    return w.IsInFiringArc(weaponAngle)
}
    