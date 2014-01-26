package main

import "container/list"

type PlanetList struct {
    list *list.List
}

func NewPlanetList() *PlanetList {
    return & PlanetList {
        list: list.New(),
    }
}

func (pl *PlanetList) Add(p *Planet) {

}

func (pl *PlanetList) Remove(p *Planet) {

}





