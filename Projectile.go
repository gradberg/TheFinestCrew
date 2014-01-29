
package main

// ---- Perhaps at some point this should only have a pointer to the weapon that fired it. The only reason I don't want to
//      is because this way it does not actually need a real weapon to fire it (in case, i don't know, its debris?)


type Projectile struct {
    Point Point // location on map
    Heading float64 // heading.    
    WeaponType WeaponTypeEnum // the type of weapon that fired this. Lasers only exist momentarily (and are displayed the following tick)
    OriginShip *Ship // Ship that fired. Prevents projectiles from hitting ship that shot on their first turn

    DesignSpeed float64 // For Guns, the speed at which the projectile initially moves. For lasers, the distance it fires too.
    DesignDamage float64 // The maximum damage caused. For guns, this damage rating decreases as the projectile slows
    
    // ## Only used for GUN projectiles
    DesignDrag float64 // How much each turn a projectile slows down
    Speed float64 // remaining speed that the projectile is still going (since it slows down due to the DesignDrag)
}

func (p *Projectile) GetFuturePoint() Point {
    return p.Point.AddVector(p.Heading, p.Speed).Round()
}
func (p *Projectile) DoMovement() {
    p.Point = p.GetFuturePoint()
    // Do Drag
    p.Speed = Round(p.Speed - p.DesignDrag, 1)
    if (p.Speed < 0.0) { p.Speed = 0.0 }
}
