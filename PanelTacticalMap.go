package main

import "fmt"
import "termbox-go"

type PanelTacticalMap struct { }
    
const (
    map_BLANK = "·"
    map_PROJECTILE = "*"
    map_SHIPS = "»"
    map_PLANET = "○"
    map_PLANET_WITH_SHIPS = "%"
    map_PLANETS = "8"        
    map_DESTROYED_SHIP = "@"
)

type spotType uint16
const (
    spot_EMPTY spotType = iota    
    spot_TRAIL // displays like spot_EMPTY, except as a different color to indicate ship trails.
    spot_LASER // displays a laser (gets overwritten by projectiles)
    spot_PROJECTILE // displays a flying projectile    
    spot_SHIP
    spot_SHIPS
    spot_PLANET
    spot_PLANET_WITH_SHIPS
    spot_PLANETS    
    
)

type spot struct {
    spotType spotType
    // When this is just a single ship
    ship *Ship 
    laser *Projectile // if this ends up being a laser map spot, this is used to look up the angle to display
    
    // When this is just a single planet?    
    
    // This indicates that whatever in this spot was hit by weapons fire
    // last turn, and should be appropriately displayed.
    wasHit bool    
}

func (p *PanelTacticalMap) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }    

    switch ch {
        case 'q':
            g.ThePlayer.DecreaseTacticalMapScale()
        case 'w':
            g.ThePlayer.IncreaseTacticalMapScale()
        case 'e':
            g.ThePlayer.ToggleTacticalMapFullScreen()
                            
        default:
            return nil
    }

    return result
}
    

func (p *PanelTacticalMap) Display(g *Game, r *ConsoleRange) {
    w, h := r.GetSize()

    r.SetAttributes(termbox.ColorRed, termbox.ColorBlack)
    r.SetBorder()
    r.SetTitle("Tactical Map")
    
    // Create a 2D array of the cells that will be visible
    grid := createBlankGrid(r)
    g.fillGridWithInformation(r, grid)
    displayGrid(r, grid)
    
    // Display any overlay information?
    
    
    // Display Control information    
    r.DisplayText(fmt.Sprintf("Scale: %8.1f", g.ThePlayer.TacticalMapScale), 1, 1)
    r.DisplayTextWithColor("[q][w]", 1, 2, termbox.ColorRed, termbox.ColorGreen | termbox.AttrBold)       
    
    r.DisplayText("Full [e]", w - 9, 0)
    r.DisplayTextWithColor("[e]", w - 4, 0, termbox.ColorRed, termbox.ColorGreen | termbox.AttrBold)
    
    
    
    // Location
    r.DisplayText(fmt.Sprintf("X %10.1f", g.PlayerShip.Point.X()), 2 ,h -1)
    r.DisplayText(fmt.Sprintf("Y %10.1f", g.PlayerShip.Point.Y()), w - 14  ,h-1)
}

func createBlankGrid(r *ConsoleRange) [][]spot {
    w, h := r.GetSize()
    grid := make([][]spot, h - 2, h - 2)
    for rowIndex := 0; rowIndex < len(grid); rowIndex++ {
        grid[rowIndex] = make([]spot, w - 2, w - 2)
        
    }
    return grid
}

// Loops over all the objects on the map, figured out which cell in the visible grid they should be in
func (g *Game) fillGridWithInformation(r *ConsoleRange, grid [][]spot) {
    w, h  := r.GetSize()        
    var centerX int = (w - 2) / 2
    var centerY int = (h - 2) / 2
    
    // Loop over every relevant object to display and figure out which dot it should be in.    
    for e := g.LaserProjectiles.Front(); e != nil; e = e.Next() {
        p := e.Value.(*Projectile)
        
        // For lasers, this finds the start and end points, uses the Bresenham line algorithm to determine
        // how to plot the line across the grid, then loops over all those points and adds them to the
        // grid if they are in the visible window.        
        rX0, rY0 := determineRelativeMapSpot(g.PlayerShip.Point, p.Point, g.ThePlayer.TacticalMapScale)
        rX1, rY1 := determineRelativeMapSpot(g.PlayerShip.Point, p.GetFuturePoint(), g.ThePlayer.TacticalMapScale)
        points := BresenhamLine2(rX0 + centerX, rY0 + centerY, rX1 + centerX, rY1 + centerY)        
        for i := 0; i < len(points); i++ {
            point := points[i]
            if (point.X < 0 || point.X >= w-2 || point.Y < 0 || point.Y >= h - 2) { continue }
            grid[point.Y][point.X].spotType = combineSpotType(grid[point.Y][point.X].spotType, spot_LASER)
            grid[point.Y][point.X].laser = p
        }        
    }
    for e := g.Projectiles.Front(); e != nil; e = e.Next() {
        p := e.Value.(*Projectile)
        rX, rY := determineRelativeMapSpot(g.PlayerShip.Point, p.Point, g.ThePlayer.TacticalMapScale)
        aX := centerX + rX
        aY := centerY + rY
        if (aX < 0 || aX >= w-2 || aY < 0 || aY >= h-2) { continue }
        grid[aY][aX].spotType = combineSpotType(grid[aY][aX].spotType, spot_PROJECTILE)
    }
    for e := g.Planets.Front(); e != nil; e = e.Next() {
        // do something with e.Value
        p := e.Value.(*Planet)
        relativeX, relativeY := determineRelativeMapSpot(g.PlayerShip.Point, p.Point, g.ThePlayer.TacticalMapScale)
        
        // if the object's coordinates are off of the visible map, don't bother
        absoluteX := centerX + relativeX
        absoluteY := centerY + relativeY
        if (absoluteX < 0 || absoluteX >= w-2 || absoluteY < 0 || absoluteY >= h-2) { continue }
        grid[absoluteY][absoluteX].spotType = combineSpotType(grid[absoluteY][absoluteX].spotType, spot_PLANET)
    } 
    for se := g.Ships.Front(); se != nil; se = se.Next() {
        s := se.Value.(*Ship)        
        relativeX, relativeY := determineRelativeMapSpot(g.PlayerShip.Point, s.Point, g.ThePlayer.TacticalMapScale)
        
        absoluteX := centerX + relativeX
        absoluteY := centerY + relativeY
        
        if !(absoluteX < 0 || absoluteX >= w-2 || absoluteY < 0 || absoluteY >= h-2) { 
            grid[absoluteY][absoluteX].spotType = combineSpotType(grid[absoluteY][absoluteX].spotType, spot_SHIP)           
            
            // assign the ship to this spot. If there are multiple ships, then the tac map
            // will not display any of them
            grid[absoluteY][absoluteX].ship = s        
            grid[absoluteY][absoluteX].wasHit = grid[absoluteY][absoluteX].wasHit || s.WasHit
        }
        
        for pp := s.PastLocations.Front(); pp != nil; pp = pp.Next() {
            point := pp.Value.(Point)
            relativeX, relativeY := determineRelativeMapSpot(g.PlayerShip.Point, point, g.ThePlayer.TacticalMapScale)
            
            absoluteX := centerX + relativeX
            absoluteY := centerY + relativeY
            if !(absoluteX < 0 || absoluteX >= w-2 || absoluteY < 0 || absoluteY >= h-2) {  
                grid[absoluteY][absoluteX].spotType = combineSpotType(grid[absoluteY][absoluteX].spotType, spot_TRAIL)
            }
        }
        
        
    }
}

// Determines which map grid the given point should be assigned to, given
// where that point is, where the center is (presumably the player), and what
// scale the map is being displayed at.
func determineRelativeMapSpot(center Point, test Point, mapSpotSize float64) (int, int) {
    x := int(Round((test.X() - center.X()) / mapSpotSize, 0))
    y := int(Round((test.Y() - center.Y()) / mapSpotSize, 0))
    return x, y
}

// ---- this is going to need to be re-written to handle a nice flexible precedence
//      system for displaying items to the map
func combineSpotType(current, additional spotType) spotType {
    if additional == spot_TRAIL {
        if (current == spot_EMPTY) {        
            // if the spot is blank, then set it as a trail, 
            return spot_TRAIL
        } else {
            // otherwise don't overwrite any more useful information on the map
            return current        
        }
    }
    if additional == spot_LASER {
        if (current == spot_EMPTY || current == spot_TRAIL) {        
            // if the spot is blank, then set it as a trail, 
            return spot_LASER
        } else {
            // otherwise don't overwrite any more useful information on the map
            return current        
        }
    }
    if additional == spot_PROJECTILE {
        if (current == spot_EMPTY || current == spot_TRAIL || current == spot_LASER) {            
            return spot_PROJECTILE
        } else {
            // projectiles don't overwrite more useful information
            return current
        }
    }

    switch (current) {        
        default: return additional            
        case spot_SHIP, spot_SHIPS:
            if (additional == spot_SHIP) {
                return spot_SHIPS
            } else if (additional == spot_PLANET) {
                return spot_PLANET_WITH_SHIPS
            } else {
                return current
            }
        case spot_PLANET, spot_PLANETS:
            if (additional == spot_SHIP) {
                return spot_PLANET_WITH_SHIPS
            } else if (additional == spot_PLANET) {
                return spot_PLANETS
            } else {
                return current
            }
        case spot_PLANET_WITH_SHIPS: return spot_PLANET_WITH_SHIPS               
    }
}

func displayGrid(r *ConsoleRange, grid [][]spot) {
    w, h := r.GetSize()
    var centerX int = (w - 2) / 2
    var centerY int = (h - 2) / 2
    
    // Display map dots itself
    for rowIndex := 0; rowIndex < len(grid); rowIndex++ {        
        for colIndex := 0; colIndex < len(grid[rowIndex]); colIndex++ {            
            var spot string = "?"
            isTrail := false       
            isLaser := false
            isDestroyed := false
            switch grid[rowIndex][colIndex].spotType {
                case spot_EMPTY: spot=map_BLANK                
                case spot_TRAIL:
                    spot=map_BLANK
                    isTrail = true
                case spot_LASER:
                    spot = Compass_GetLineHeadingIcon(grid[rowIndex][colIndex].laser.Heading)
                    isLaser = true
                case spot_PROJECTILE:
                    spot=map_PROJECTILE
                    isTrail = true
                case spot_SHIP: 
                    isDestroyed = grid[rowIndex][colIndex].ship.IsDestroyed()
                    
                    if (isDestroyed) {
                        spot = map_DESTROYED_SHIP
                    } else {
                        spot = Compass_GetShipHeadingIcon(grid[rowIndex][colIndex].ship.ShipHeadingInDegrees)
                    }                    
                    
                case spot_SHIPS: spot=map_SHIPS 
                case spot_PLANET: spot=map_PLANET       
                case spot_PLANET_WITH_SHIPS: spot=map_PLANET_WITH_SHIPS             
                case spot_PLANETS: spot=map_PLANETS
            }
            
            // Flip the rows vertically. In the math, the rows go from 0 to x downards, but visually
            // they go from positive to negative.
            // ---- This may be a sign that I wrote the algorithm wrong, and it should be consistantly interpreting
            //      it the other way around
            displayY := h - 2 - rowIndex
            
            bg := termbox.ColorBlack
            if (grid[rowIndex][colIndex].wasHit) { 
                bg = termbox.ColorYellow | termbox.AttrBold
            } else if isDestroyed {
                bg = termbox.ColorBlack | termbox.AttrBold
            }
            
            fg := termbox.ColorRed
            if (colIndex == centerX && rowIndex == centerY) {
            } else if isTrail {
                fg = termbox.ColorWhite | termbox.AttrBold
            } else if isLaser {
                fg = termbox.ColorGreen | termbox.AttrBold
            }
            
            r.DisplayTextWithColor(spot, colIndex + 1, displayY, fg, bg)
        }            
    }
}

