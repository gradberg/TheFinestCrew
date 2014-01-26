/*
    All objects in space have a common set of information, which thus can be relied
    upon regardless of what they are (ship versus planet versus ?)
*/

package main


type ISpaceObject interface {
    GetPoint() Point
    GetCourse() float64
    GetSpeed() float64
    GetHeading() float64
    GetName() string
}