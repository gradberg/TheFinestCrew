/*
    Compass functions comprise showing ships, arrows or other information directionally from 0 to 360.0 degrees
*/

package main
import "math"

func Compass_GetShipHeadingIcon(angle float64) string {
    const segmentSize = 22.5
    // adjust the ship's heading to get the correct
    adjustedHeading := math.Mod(angle + (segmentSize / 2.0), 360.0)
    segment := int16(adjustedHeading / segmentSize)
    
    switch segment {
        case 0: return "^"
        case 1: return "╕"
        case 2: return "┐"
        case 3: return "╖"
        case 4: return ">"
        case 5: return "╜"
        case 6: return "┘"
        case 7: return "╛"
        case 8: return "v"
        case 9: return "╘"
        case 10: return "└"
        case 11: return "╙"
        case 12: return "<"
        case 13: return "╓"
        case 14: return "┌"
        case 15: return "╒"
    }
    return "?"   
}

func Compass_GetNearDirectionArrow(angle float64) (string, int, int) {
    const segmentSize = 22.5
    // adjust the ship's heading to get the correct
    adjustedHeading := math.Mod(angle + (segmentSize / 2.0), 360.0)
    segment := int16(adjustedHeading / segmentSize)
    
    switch segment {
        case 0: return "|", 0, -1
        case 1: return "/", 0, -1
        case 2: return "/", 1, -1
        case 3: return "/", 1, 0
        case 4: return "-", 1, 0
        case 5: return "\\", 1, 0
        case 6: return "\\", 1, 1
        case 7: return "\\", 0, 1
        case 8: return "|", 0, 1
        case 9: return "/", 0, 1
        case 10: return "/", -1, 1
        case 11: return "/", -1, 0
        case 12: return "-", -1, 0
        case 13: return "\\", -1, 0
        case 14: return "\\", -1, -1
        case 15: return "\\", 0, -1
    }
    return "?", 0, -1   
}

func Compass_GetFarDirectionArrow(angle float64) (string, int, int) {
    const segmentSize = 22.5
    // adjust the ship's heading to get the correct
    adjustedHeading := math.Mod(angle + (segmentSize / 2.0), 360.0)
    segment := int16(adjustedHeading / segmentSize)
    
    switch segment {
        case 0: return "|", 0, -2
        case 1: return "/", 1, -2
        case 2: return "/", 2, -2
        case 3: return "/", 2, -1
        case 4: return "-", 2, 0
        case 5: return "\\", 2, 1
        case 6: return "\\", 2, 2
        case 7: return "\\", 1, 2
        case 8: return "|", 0, 2
        case 9: return "/", -1, 2
        case 10: return "/", -2, 2
        case 11: return "/", -2, 1
        case 12: return "-", -2, 0
        case 13: return "\\", -2, -1
        case 14: return "\\", -2, -2
        case 15: return "\\", -1, -2
    }
    return "?", 0, -1   
}

type CompassPoint struct { X, Y int }
func Compass_GetArcPoints(start, end float64) []CompassPoint {
    points := [32]CompassPoint { 
        CompassPoint {  0, -1 },
        CompassPoint {  0, -1 },
        CompassPoint {  1, -1 },
        CompassPoint {  1,  0 },
        CompassPoint {  1,  0 },
        CompassPoint {  1,  0 },
        CompassPoint {  1,  1 },
        CompassPoint {  0,  1 }, 
        CompassPoint {  0,  1 },
        CompassPoint {  0,  1 },
        CompassPoint { -1,  1 },
        CompassPoint { -1,  0 },
        CompassPoint { -1,  0 },
        CompassPoint { -1,  0 },
        CompassPoint { -1, -1 },
        CompassPoint {  0, -1 }, 
        CompassPoint {  0, -1 },
        CompassPoint {  0, -1 },
        CompassPoint {  1, -1 },
        CompassPoint {  1,  0 },
        CompassPoint {  1,  0 },
        CompassPoint {  1,  0 },
        CompassPoint {  1,  1 },
        CompassPoint {  0,  1 }, 
        CompassPoint {  0,  1 },
        CompassPoint {  0,  1 },
        CompassPoint { -1,  1 },
        CompassPoint { -1,  0 },
        CompassPoint { -1,  0 },
        CompassPoint { -1,  0 },
        CompassPoint { -1, -1 },
        CompassPoint {  0, -1 }, 
    } 
    
    // If the start angle is LARGER than the end angle, add 360 to the end angle and it
    // will cause the slice to use the second half of the points (which is what we want)    
    const segmentSize = 22.5
    if start > end { end += 360.0 }
    
    adjustedStart := int16((start + (segmentSize / 2.0)) / segmentSize)
    adjustedEnd := int16((end + (segmentSize / 2.0)) / segmentSize)
    
    return points[adjustedStart:adjustedEnd+1]
}