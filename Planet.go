

package main

type Planet struct {
    Point Point // location
    Name string
}

func NewPlanet() *Planet {
    return & Planet {            
        Name: "NAME THE PLANET",
    }
}


// Satisfy the ISpaceObject interface
func (p *Planet) GetPoint() Point { return p.Point }
func (p *Planet) GetCourse() float64 { return 0.0 } 
func (p *Planet) GetSpeed() float64 { return 0.0 }
func (p *Planet) GetHeading() float64 { return 0.0 } 
func (p *Planet) GetName() string { return p.Name }



