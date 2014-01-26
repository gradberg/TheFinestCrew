/*
    This interface defines how the program can run the actions for an AI crew member.
    *NOTE* These objects are used in a STATEFUL manner, so a given crew member is initialized,
    assigned to a ship, and the object stays alive so that the crew member can have a 'memory'
    
    If I want these to be test-able, it will need to be on an individual basis.
*/

package main

type IAiCrew interface {       
    
    // The crew member performs their actions and returns a number indicating how many ticks the action
    // takes up.
    DoAction(s *Ship, g *Game, cm *CrewMember) int
    
}
