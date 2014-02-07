/*
    The Overlay Panel is the one the display any additional inforjmation above and below
    the panels
*/

package main
//import "fmt"
import "termbox-go"

type PanelTitleScreen struct { }

func (p *PanelTitleScreen) ProcessInput(g *Game, ch rune, key termbox.Key) *InputResult {
    result := &InputResult {
        Exit: false,
        TicksToPass: 0,    
    }

    if (g.GameStatus == GameStatusTitleScreen) {    
        switch ch {
            // Quit
            case '1':
                //g.doSetupForDevelopment()
                g.GameStatus = GameStatusPickShip
                
            case 'Q', 'q': 
                result.Exit = true
                
            default:
                return nil
        }
    } else if (g.GameStatus == GameStatusPickShip) {
        switch ch {
            case '1':
                g.GameStatus = GameStatusTitleScreen                
            case '2':
                g.GameStatus = GameStatusPickRole
                g.GameSetup.Ship = PlayableShipMilitiaCorvette
                
            case '3':
                g.GameStatus = GameStatusPickRole
                g.GameSetup.Ship = PlayableShipPirateFighter
                
            default:
                return nil
        }
    } else if (g.GameStatus == GameStatusPickRole) {
        switch ch {
            case '1':
                g.GameStatus = GameStatusPickShip
            
            case '2':
                if (g.GameSetup.Ship == PlayableShipMilitiaCorvette) {                
                    g.GameSetup.CrewRole = CrewRoleHelmsman
                    g.GameStatus = GameStatusInstructions        
                } else {                    
                    g.GameSetup.CrewRole = CrewRolePilot
                    g.GameStatus = GameStatusInstructions
                }
                
            case '3':                 
                if (g.GameSetup.Ship == PlayableShipMilitiaCorvette) {                
                    g.GameSetup.CrewRole = CrewRoleCommander
                    g.GameStatus = GameStatusInstructions       
                } else {
                    return nil
                }            
                
            default:
                return nil
        }
    } else if (g.GameStatus == GameStatusInstructions) {
        switch ch {
            case '1':
                g.GameStatus = GameStatusPickRole
            
            case '2':                
                //g.doSetupForDevelopment()
                g.setupForStart()
                g.setupForNextRound() // I do not like this here...
                g.GameStatus = GameStatusPlaying            
                
            default:
                return nil
        }    
    }
    
    return result
}

func (p *PanelTitleScreen) Display(g *Game, r *ConsoleRange) {
    if (g.GameStatus == GameStatusInstructions) {
        p.displayInstructions(g,r)
    } else {
        p.displayNormal(g,r)
    }
}

func (p *PanelTitleScreen) displayNormal(g *Game, r *ConsoleRange) {
    title := []string {
        // 0
        "                                                                                ",
        "    Here's to                                                                   ",
        "                                                                                ",
        "    TTTTT  H   H  EEEEE          FFFFF  IIIII  N   N  EEEEE   SSSS  TTTTT       ",        
        "      T    H   H  E              F        I    NN  N  E      S        T         ",
        "      T    HHHHH  EEE            FFF      I    N N N  EEE     SSS     T         ",
        "      T    H   H  E              F        I    N  NN  E          S    T         ",
        "      T    H   H  EEEEE          F      IIIII  N   N  EEEEE  SSSS     T         ",
        "                                                                                ",
        "                                                                                ",
        "            CCCC  RRRR   EEEEE  W   W                                           ",
        "           C      R   R  E      W W W                                           ",
        "           C      RRRR   EEE    W W W                                           ",
        "           C      R  R   E       W W                                            ",
        "            CCCC  R   R  EEEEE   W W                                            ",
        "                                                                                ",
        "                        in the fleet.                                           ",
        "                                                                                ",
        "                                                                                ",
        "                                                                                ",
        // 20
    }        
    for i, s := range title {
        r.DisplayText(s, 0, i)
    }   
    
    // ---- make this randomly display a different ship portrait on the title screen    
    p.displayShipCorvette(g,r,45,10)
    //p.displayShipFighter(g,r,45,10)
    
    black := termbox.ColorBlack
    white := termbox.ColorWhite | termbox.AttrBold
    
    switch g.GameStatus {
        case GameStatusTitleScreen:            
            r.Com("[1]"," Start Game", 5, 19, black, white)
            r.Com("[Q]"," Quit", 5, 20, black, white)    
            
        case GameStatusPickShip:
            r.Com("[1]", " Go Back", 5, 19, black, white)
            r.Com("[2]", " Militia Corvette", 5, 20, black, white)
            r.Com("[3]", " Pirate Fighter", 5, 21, black, white)
        
        case GameStatusPickRole:
            r.DisplayText(" > " + g.GameSetup.Ship.ToString(), 5, 18)
            r.Com("[1]"," Go Back", 5, 19, black, white)
            
            if (g.GameSetup.Ship == PlayableShipMilitiaCorvette) {
                r.Com("[2]"," Play as Helmsman", 5, 20, black, white)
                r.Com("[3]"," Play as Commander", 5, 21, black, white)                
            } else {
                r.Com("[2]"," Play as Pilot", 5, 20, black, white)
            }
    }
}

func (p *PanelTitleScreen) displayInstructions(g *Game, r *ConsoleRange) {
    black := termbox.ColorBlack
    white := termbox.ColorWhite | termbox.AttrBold
    grey := termbox.ColorWhite

    tl := NewFlowDocument(38, 11)
    tl.AddParagraph("Game Mode: Survival", white, black)
    tl.AddParagraph("", grey, black)
    tl.AddParagraph("Each round increasing numbers of enemy ships will be present, which you must all destroy to proceed to the next round. The only end to this mode is when your ship is destroyed.", grey, black)
    tl.Write(r, 1,1)
    
    bl := NewFlowDocument(38, 11)
    bl.AddParagraph("Player Role: " + g.GameSetup.CrewRole.ToString(), white, black)
    bl.AddParagraph("", grey, black)
    switch (g.GameSetup.CrewRole) {
        case CrewRoleCommander:
            bl.AddParagraph("As commander, your job is to give orders to the other supporting crew members while you directly control the weapons systems.", grey, black)    
        case CrewRoleHelmsman:
            bl.AddParagraph("You fly the ship per orders from your commanding officer.", grey, black)    
        case CrewRolePilot:
            bl.AddParagraph("As the sole crew member, you both fly the ship and cnotrol the weapons.", grey, black)    
    }
    bl.Write(r, 1, 13)
        
    switch (g.GameSetup.Ship) {
        case PlayableShipMilitiaCorvette:
            rd := NewFlowDocument(38, 23)                
            rd.AddParagraph(g.GameSetup.Ship.ToString(), white, black)            
            description := []string {
                "Bridge Crew - 2",
                "Weapons - Forward Facing Cannon, 360° Laser Coverage",
                "Shields - None",
                "Rotation - Very Good (18°/t)",
                "Thrust - Decent (1/t)",
                " ",
                "Often used by poorly equiped, regional police forces, the Militia Corvette is ill suited for serious combat. Lacking shield generators, it cannot last long in a fight, and instead must appear in numbers to provide any real threat.",
            }
            
            for _, s := range description {
                rd.AddParagraph(s, grey, black)
            }                   
            rd.Write(r, 41, 6)
                
            p.displayShipCorvette(g, r, 48, 1)
            
        case PlayableShipPirateFighter:            
            rd := NewFlowDocument(38, 23)
            rd.AddParagraph(g.GameSetup.Ship.ToString(), white, black)
            description := []string {
                "Bridge Crew - 1",
                "Weapons - Forward Facing Laser",
                "Shields - None",
                "Rotation - Very Good (18°/t)",
                "Thrust - Terrible (0.4/t)",
                " ",
                "The worst of the worst, fighter-craft used by pirates are a joke. Pilots are usually coerced into flying them, as such an assignment is almost assuredly a death warrant.",
            }
            
            for _, s := range description {
                rd.AddParagraph(s, grey, black)
            }                   
            rd.Write(r, 41, 4)
    
            p.displayShipFighter(g,r,52,1)
    }    
    
        
    
    r.Com("[1]"," Go Back", 1, 24, black, white)
    r.Com("[2]"," Start", 50, 24, black, white)   

    

    
}


func (p *PanelTitleScreen) displayShipCorvette(g *Game, r *ConsoleRange, x, y int) {
    title := []string {
        "           --▄▄          ",
        "    \\▄    ▄██████     █  ",
        " ▄█████████████████████  ",
        "██▀   /▀     ▀\\    ████  ",
        // 20
    }        
    for i, s := range title {
        r.DisplayText(s, x, y + i)
    }   

    r.DisplayTextWithColor("▄█", x + 1, y + 2, termbox.ColorBlue | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayVerticalTextWithColor("███", x + 23, y + 1, termbox.ColorRed, termbox.ColorBlack)
    r.DisplayTextWithColor("██", x + 24, y + 1, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayTextWithColor("███", x + 24, y + 2, termbox.ColorYellow | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayTextWithColor("██", x + 24, y + 3, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack)
}

func (p *PanelTitleScreen) displayShipFighter(g *Game, r *ConsoleRange, x, y int) {
    title := []string {
        "        ▄█  ",
        "_▄▄████████ ",
        "       ▀▀▀  ",
        // 20
    }        
    for i, s := range title {
        r.DisplayText(s, x, y + i)
    }   

    r.DisplayTextWithColor("█", x + 3, y + 1, termbox.ColorBlue | termbox.AttrBold, termbox.ColorBlack)
    r.DisplayTextWithColor("█", x + 11, y + 1, termbox.ColorRed, termbox.ColorBlack)
    r.DisplayTextWithColor("██", x + 12, y + 1, termbox.ColorRed | termbox.AttrBold, termbox.ColorBlack)
}



 