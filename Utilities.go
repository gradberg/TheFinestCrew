package main

import "math"
import "strings"


func TrimLength(s string, length int) string {
    if len(s) > length {
        return s[0:length]
    } else {
        return s
    } 
}

func TrimOrPadToLength(s string, length int) string {
    if (len(s) > length) {
        return s[0:length]
    } else if (len(s) < length) {
        gap := length - len(s)
        return s + strings.Repeat(" ", gap)
    } else {
        return s
    }
}

// Round return rounded version of x with prec precision.
func Round(x float64, prec int) float64 {
    var rounder float64
    pow := math.Pow(10, float64(prec))
    intermed := math.Abs(x) * pow
    _, frac := math.Modf(intermed)

    if frac >= 0.5 {
        rounder = math.Ceil(intermed)
    } else {
        rounder = math.Floor(intermed)
    }
    
    sign := 1.0
    if (x < 0.0) { sign = -1.0 }

    return rounder * sign / pow
}

func Round64(f float64) int64 {
    //return int64(math.Floor(f+0.5))
    return int64(Round(f, 0))
}



func VectorToXy(angle, magnitude float64) (float64, float64) {
    radians := angle * math.Pi / 180.0
    x := math.Sin(radians) * magnitude
    y := math.Cos(radians) * magnitude
    return x, y
}

func XyToVector(x, y float64) (float64, float64) {
    radians := math.Atan2(x, y)
    degrees := radians * 180.0 / math.Pi    
    if degrees < 0.0 { degrees += 360.0 }
    if degrees >= 360.0 { degrees -= 360.0 }
    if math.IsNaN(degrees) { degrees = 0.0 }
    
    magnitude := math.Sqrt((x * x) + (y * y))
    return degrees, magnitude
}

// Given a number of degreee 0 (inclusive) to 360 (exclusive), it returns the
// opposite direction
func GetOppositeDegrees(angle float64) float64 {
    angle += 180.0
    if angle >= 360.0 { angle -= 360.0 }
    return angle
}
func AddAngles(angle1, angle2 float64) float64 {
    a := angle1 + angle2
    if (a >= 360.0) { a -= 360.0 }
    if (a < 0.0) { a += 360.0 }
    return a
}

type Point struct {
    x, y float64
}
func NewPoint(x, y float64) Point {
    return Point { x:x, y:y, }
}
func (p Point) X() float64 { return p.x }
func (p Point) Y() float64 { return p.y }
func (p Point) ToVector() (float64, float64) {
    return XyToVector(p.x, p.y)
}
func (p Point) AddVector(angle, magnitude float64) Point {
    x, y := VectorToXy(angle, magnitude)
    return NewPoint(p.x + x, p.y + y)
}
func (p Point) Add(p2 Point) Point {
    return NewPoint(p.x + p2.x, p.y + p2.y)
}
func (p Point) Subtract(p2 Point) Point {
    return NewPoint(p.x - p2.x, p.y - p2.y)
}
func (p Point) Round() Point {
    return NewPoint(Round(p.x, 1), Round(p.y, 1))
}
func (p Point) DistanceFrom(p2 Point) float64 {
    a := (p.x - p2.x)*(p.x - p2.x) + (p.y - p2.y)*(p.y-p2.y)
    return math.Sqrt(a)
}

//
//  Copied from http://forums.codeguru.com/showthread.php?194400-Distance-between-point-and-line-segment
//
func DistanceFromLineSegment(test Point, lineA Point, lineB Point) float64 {
    cx := test.X()
    cy := test.Y()
    ax := lineA.X()
    ay := lineA.Y()
    bx := lineB.X()
    by := lineB.Y()

	r_numerator := (cx-ax)*(bx-ax) + (cy-ay)*(by-ay)
	r_denomenator := (bx-ax)*(bx-ax) + (by-ay)*(by-ay)
	r := r_numerator / r_denomenator

	if ( (r >= 0) && (r <= 1) ) {        
        s := ((ay-cy)*(bx-ax)-(ax-cx)*(by-ay) ) / r_denomenator
        return math.Abs(s)*math.Sqrt(r_denomenator)
	} else {
		dist1 := (cx-ax)*(cx-ax) + (cy-ay)*(cy-ay)
		dist2 := (cx-bx)*(cx-bx) + (cy-by)*(cy-by)
		if (dist1 < dist2) {
            return math.Sqrt(dist1)
		} else {
            return math.Sqrt(dist2)
		}
	}
}
 

// returns the degrees to turn and true if it should be clockwise, or false if not
func GetShortestTurn(currentAngle, targetAngle float64) (float64, bool) {
    clockwiseAngle := targetAngle - currentAngle
    if clockwiseAngle < 0.0 { clockwiseAngle += 360.0 }
    
    counterAngle := currentAngle - targetAngle
    if counterAngle < 0.0 { counterAngle += 360.0 }
    
    if clockwiseAngle < counterAngle {
        return clockwiseAngle, true
    } else {
        return counterAngle, false
    }    
}
