 
package main
import "testing"
 
func Test_DistanceFromLineSegment(t *testing.T) {
    //result := DistanceFromLineSegment(NewPoint(0.707, 0.707), NewPoint(0,0), NewPoint(1,0) )
    assertEqualF64(t, 0.707, DistanceFromLineSegment(NewPoint(0.707, 0.707), NewPoint(0,0), NewPoint(1,0)))
    assertEqualF64(t, 1.0,   Round(DistanceFromLineSegment(NewPoint(-0.707, 0.707), NewPoint(0,0), NewPoint(1,0)), 2))
    assertEqualF64(t, 216.37, Round(DistanceFromLineSegment(NewPoint(200, -106), NewPoint(0,0), NewPoint(100,100)),2))
}


func assertEqualF64(t *testing.T, expected, test float64) {
    if (test != expected) {
        t.Errorf("Expected %f, Actual %f", expected, test)
    }
}