 
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
        errorNowf(t, "Expected %f, Actual %f", expected, test)
    }
}

func Test_BresenhamLine_Normal(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1  ,2,1  ,3,2  ,4,2  ,5,3  ,6,3  ,7,3  ,8,4  ,9,4  ,10,5   ,11,5),
        BresenhamLine2(1,1,11,5), 
    )
}
func Test_BresenhamLine_NormalReverse(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1  ,2,1  ,3,2  ,4,2  ,5,3  ,6,3  ,7,3  ,8,4  ,9,4  ,10,5   ,11,5),
        //coordinatesToLine(11,5  ,10,5  ,9,4  ,8,4  ,7,3  ,6,3  ,5,3  ,4,2  ,3,2  ,2,1  ,1,1),
        BresenhamLine2(11,5,1,1), 
    )
}
func Test_BresenhamLine_Steep(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1  ,1,2  ,2,3  ,2,4  ,3,5  ,3,6  ,3,7  ,4,8  ,4,9  ,05,10  ,05,11),
        BresenhamLine2(1,1,5,11), 
    )
}
func Test_BresenhamLine_SteepReverse(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1  ,1,2  ,2,3  ,2,4  ,3,5  ,3,6  ,3,7  ,4,8  ,4,9  ,05,10  ,05,11),
        BresenhamLine2(5, 11, 1, 1), 
    )
}
func Test_BresenhamLine_ZeroLength(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1),
        BresenhamLine2(1,1, 1, 1), 
    )
}
func Test_BresenhamLine_VerticalLine(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1  ,2,1  ,3,1  ,4,1  ,5,1),
        BresenhamLine2(1,1, 5, 1), 
    )
}
func Test_BresenhamLine_HorizontalLine(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(1,1  ,1,2  ,1,3  ,1,4  ,1,5),
        BresenhamLine2(1,1, 1, 5), 
    )
}
func Test_BresenhamLine_AllNegativeLine(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(-11,-5  ,-10,-5  ,-9,-4  ,-8,-4  ,-7,-3  ,-6,-3  ,-5,-3  ,-4,-2  ,-3,-2  ,-2,-1  ,-1,-1),
        BresenhamLine2(-1,-1, -11, -5), 
    )
}
func Test_BresenhamLine_AllNegativeLineReverse(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(-11,-5  ,-10,-5  ,-9,-4  ,-8,-4  ,-7,-3  ,-6,-3  ,-5,-3  ,-4,-2  ,-3,-2  ,-2,-1  ,-1,-1),
        BresenhamLine2(-11,-5, -1, -1), 
    )
}
func Test_BresenhamLine_CrossingAxis(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(-4,-4  ,-3,-3  ,-2,-2  ,-1,-1  ,0,0  ,1,1  ,2,2),
        BresenhamLine2(-4,-4, 2, 2), 
    )
}
func Test_BresenhamLine_CrossingReverse(t *testing.T) {
    assertPointsEqual(t,
        coordinatesToLine(-4,-4  ,-3,-3  ,-2,-2  ,-1,-1  ,0,0  ,1,1  ,2,2),
        BresenhamLine2(2,2,  -4,-4), 
    )
}


func coordinatesToLine(coordinates ...int) []LinePoint {
    points := make([]LinePoint, 0, 10)
    
    pairCount := len(coordinates) / 2
    for i := 0; i < pairCount; i++ {
        points = append(points, LinePoint{X: coordinates[i * 2], Y: coordinates[i * 2 + 1]})
    }
    
    return points
}
func assertPointsEqual(t *testing.T, expected, actual []LinePoint) {
    if (len(expected) != len(actual)) {
        errorNowf(t, "Expected %d points, Actual %d points", len(expected), len(actual))
    }
    for i := 0; i < len(expected); i++ {
        if (expected[i].X != actual[i].X || expected[i].Y != actual[i].Y) {
            errorNowf(t, "Expected %d,%d ; Actual %d,%d", expected[i].X, expected[i].Y, actual[i].X, actual[i].Y)
        } 
    }
}

func errorNowf(t *testing.T, format string, args ...interface{}) {
    t.Errorf(format, args...)
    t.FailNow()
}