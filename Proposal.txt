# Summary #
In this game, you play one of the important crew members responsible for controlling a star-trek-style starship in its travels and battles in the galaxy. You control only a single crew member who is responsible for their aspect of the ship, and thus must rely on the other AI crew members doing their job? 

# Why #
First off, the one of the very few games that you actually get to control a major ship is Faster Than Light. That game however focuses entirely on manipulating the crew. Rather, I would like the game to focus on the ship, its interactions with other ships and planets, and how you control that as your particular crew member.

Also, this game is not an attempt to have any sort of realistic physics, or feasible space faring societies and military organizations (discussions such as big ships vs fighters, drones, etc). This is meant to be a fun game where you are just one part of a big ship, which may be part of a big fleet, etc. 

# Ideas #
Each 'class' needs ways to optimize their job. Helmsman cannot just be about turnning and moving the ship. You need to be able to spend cycles investigating ways to do 'better'. Perhaps you can program manuevers prior to executing them, so that while something is happening, you can queue up the next set of moves if it works out just fine.

All actions take time, just like a normal roguelike. The difference is that the ship is still moving and other crew members are still doing stuff too.

The HUD displays orders the captain has given. Perhaps he keeps track of who is doing a good job and not, and replaces poor crew members.

Only status messages referring to YOU yourself (like the ship makes funny noises, you are hurt) will you need to be informed outside of the personnel panel. Anything else, like what people say, is shown on the personnel panel. If there are too many things said (likely), then it automatically alternates between them so the player doesn't have to scroll through them manually (because they would have heard everything automatically), which differs from the communication officer, who would have to manually receive and read each message individually.

For whatever reason, there COULD be fighter craft. If the player controls it, then his command console has all of the different panels necessary for piloting such a craft.

Commands need to clearly state whether they consume time (and thus you lose turns) or whether they represent information you simply have and probably heard while doing your own tasks.

Different roles and ships could have different panel arrangements on the screen, not all just a simple 4x4 grid.

Commands that take time will be a caution color (like yellow). Commands that do not take time will be a safe color, like Green.
Commands that take alot of time could be like orange?

Movement. I could take a hybrid approach whereby after a certain maximum speed attempts to accelerate provide dimenshing returns, but still use the same amount of fuel. Thus if you want to escape from an enemy, you could have the helmsman Hit the afterburners to try to escape? Probably would be a slow decay of speed as well

Players should have the ability to play time real-time, perhaps 1 tick per second (and can be adjusted by the player). Thus if they know they will need to wait for a while, they can just let it go. If something happens that they need to be concerned about, the game can stop playback for them.

There could be slightly differently gameplay modes. One where you work your way up from a lower crew member like helmsman on a large warship. Another where you start out as a captain on a dumpy little spaceship and you work you way up in ships. Crew members would have some skills and abilities which could level-up during this time. Furthermore, ships could collect better weapons and equipment (maybe) during their voyages. Like a traditional roguelike, you would have missions in random maps which are equivalent to a level in a dungeon.

if the "Intertial Dampers" go out, then it is up to the helmsman to not produce so bad of G-Forces that the crew start throwing up (which loses them ticks, and lower's the Helmsman's status rating)

Ships could have varying numbers of shield generators. Thus small ships may have as few a 1, which means the single set of shields covvers the whole ship, whereas large ones may have numerous (though capped, because I don't know how it would display well in the GUI) generators. The orientation of the shields could be set at any angle, not just 0, 90, 180, 360

Physical weapons (like torpedoes and missles) could have some sort of pattern or programability. Given shields are directional, some missiles may be able to strike from the side (but not flying directly at the ship, and then turning to hit a side). If an enemy ship's shields are damaged and they're leaving a side open open, a certain type of missile could take it out.

One aspect of mission progression is that you could accumulate additional ships that follow your ship's orders. Thus you can tackle more and worse enemies because you have several ships, until a point where you are just too outmatched and get blown to bits (woot Roguelike)

What if in addition to ship subsystems, ships could actually have pieces that break off. Specifically, if you think of the Star Trek federation ship, if one of their warp nacelles where broken off, it would be both unrepairable as its missing, and it would form debris in space that could be an issue in combat. This fits nicely into potential asteroid fields on the map : D

Planets should be very large (so that you can run into them) and their size should be represented on the map. In fact, if ships started having significant size, then zooming in close enough would show the ship taking up multiple locations. For this to not suck though, ships would then almost need a size or shapethat would get displayed.

# Battles and Combat? #
How do weapons work? How does damage work? Is it a simple HP/Shield number (shields are oriented around the ship?)?
Ships need to be composed of subsections, that way everything must be repaired in order to work.

I think at least half of weapons should be kinetic impact, because that allows weapons to divert the direction that ships go. Perhaps all weapons hit shields first, but any physical ones, if impacting the hull, and really send it spinning and flying.

# Missions #
So similiar to how a traditional roguelike has dungeon levels, this has missions. You arrive on a map, have some objective to accomplish, and when you are done, you jump out. You can jump out at any time, but if you do, you (1) lose bonuses given to you after the mission, (2) risk losing the game by getting fired or demoted (THERE YOU GO, you get DEMOTED to where you cannot control the jump-out), and (3) lose any bonuses you may gain by fighting (by discovering new technologies, new weaponary, etc?)

Maps should have a decent amount on them, of which only 25% is relevannt to your given mission. Missions could include Destroy a specific fleet, escort a friendly ship to avoid destruction, attack a planet, retrieve a disabled friendly ship from an asteroid field.

# Game Modes #

Career Mode: You start as a lower-level crewmember and advance (helmsman -> communications -> defense officer, tactical officer, captain, admiral). Each mission or so that the ship survives and for which you do a 'good' job, you would advance. 

>>> Career Mode needs a GOAL. Is there something you can acheive all the way at the end? What is it? Does it end with you as the admiral of a big fleet (in your flagship, of course), ordering fleet movements across a very large map, to conquer the enemy capital planet?

Survival: 

Preset Battles: 

Custom Battles:

# Command Panels #

## Helm Control
* Has the ability to switch between direct mode  (where pressing certain buttons thrusts the ship a little each direction) and pilot mode where you specify the heading and speed you want, and the ship adjusts to match it. 
* More advanced ships may have a 'jump drive' which rapidly acccelerates and then deccelerates a ship for a long distance (slightly randomized, based on ship quality and crew skill). It would need to be 'charged' for some amount of time before being available, at which point it can be used (and recharged again). If and when power levels are implemented, charging the jump drive would be expensive in terms of power
* You must first 'plan' a maneuver (which takes a tick), then set the Auto Pilot to execute it. In the planning window, it will show preset maneuver plus many of them you have used so far. This saves the player time, but they still must 'enter' the manuever in (even if they don't have to plan out each step) so they are still forced to spend a tick getting it ready to be run by the Auto Pilot.

## Helm Status
* Its possible this could be eliminated in favor of the Helm Map, but no guarantee
* Display any G-Force information?

## Helm Map
* Displays locations of objects, including their velocity and heading
* Differs from regular tactical map in that it displays information about each objects' 

## Tactical Map
* Displays locations of objects, their velocity, their knowable sheild, weapon, and armor statuses

## Captain
* Provides a quick system for choosing a given crewmember, and based on their role, giving them orders. (Sort of a tree-like system)

## Personnel
* Look around the room. Looking around the room takes time (as you're not looking at your console to control or see things), so it first tells you what you saw LAST time and when that was. (This represents your memory that you would have as a person in that situation) You can then select to look around the room again. This is a way to ascertain the health, status, and even presence of the other crew members, including visible damage to the bridge of the ship. 

## Shield Control
* Orient the shields
* Distribute power between strength and recharge
* Change sheild frequency to reduce weapons lock relability (Weapons can lock on better by identifying the target ship's shield frequency. By changing it, it eliminates that lock and forces the enemy tactical officer)

## Tactical Scan
* Shows your ship versus another target ship. It shows the ship orientation of each with respect to each other, AND (if sensors are not damaged) the shield ratings and their orientation. Thus if involved in multi-ship combat, one of the ships may wait for the opportune moment to fire when the enemy's ships have a gap and they're trying to orient them to protect from a different enemy.
* Shows distance, heading, bearing, speed of each, potentially engine 'status'?

# Roles #
## Captain
* Directs rest of crew, and thus controls the whole ship's actions.
* Can view mission (maybe a mission panel?) to see what objectives remain

## Helmsman
* Steers ship and plots courses

## Tactical Officer
* Targets enemy ships, fires weapons, tries to improve hit chances
* Reports upcoming threats to captain
 
## Defense Officer
* Manages shields, analyzes damage status and directs repair crews
* Launches countermeasures and the like to confuse missiles

## Communications Officer
* Receives messages from other ships and conveys them to the captain. 
* Sends messages to other ships

## Science Officer?
## Chief Engineer
* Directs power to different areas of ship based on need. 
* Increases power output
* Increases speed output?
* Handles power recharhing/reloading

#
#
# Ship Systems #
These are systems that can be damaged by combat, and thus negatively affect the ship in various ways until repaired. It is possible that they can be permanently destroyed for the mission. After each mission, there is a chance that a permanently damaged system is repaired (parts are found from friendly sources?). But that also means there is a chance that you go into the next mission with a crippled ship.

All sub-systems can take a range of damage that accumulates, the otherwise simply either function or not. Thinking about Rules of Engagement 2, when systems got damaged and just worked WORSE (not on/off), it becomes confusing and irritating not knowing what their status is. I can avoid that by allowing them to be damaged, but at some point they are either functioning correctly or not. This makes it straight-forward for a defense officer to prioritize damage repair teams.

Additionally, its possible that long-term each one of these systems could also have a power-level, which allows the chief engineer to turn off power to non-essential systems if generators are knocked-out. I do not know if I like this idea or not.

* Inertial Dampers: The ship changing its heading (rotational G-Forces) and speeding-up / slowing-down causes G-Forces on the crew. If too much G-Forces are exerted in a short time then crew members throw-up, causing them to lose ticks, or they can black-out entirely.
* Primary Life-support: With primary life support out, crew members are forced to rely on respirators and such rooms that have their owm auxliary life support. This means crew members randomly lose a tick gasping for air or adjusting their breathing mask
* Engines: Ships will have one or more engine units, each which can have a different power output (so ships may have big engines and small engines, likely focusing repairs on the biggest ones first). When a given engine is nonfunctional, the ship loses that thrust capability. In face, these are likely directional, so rotational thrusters can be damaged, as well as forwards engines and backwards engines. (In bad situations, you may be reduced to navigating BACKWARDS with manual control. The Automatic Pilot will be unable to handle a capability beng gone.)
* Weapons: Obviously weapons can be damaged
* Long-range Sensors: Maps past a certain distance will be unviewable??
* Short-range Sensors: Maps still viewable without Long-range ensors will lose all detail about enemy ships and planets
* Shields: As ships have multi directional shields, shield generators can be knockced out which created gaps in those relevant ship sections.
* Main Computer: While the ship is properply designed to be functional in the event of main computer failure, it does cause things like Auto Piloting, weapons locking, and any other sort of analytical functionality to be unavailable.
* Communications Antenna: Ship loses the ability to send or receive communications from other ships
* Power Generators: Like engines, ships have multiple generators which each feed power. As they are knocked-out, less power is available which much be focused on different subsystems.
* Jump Drive: The ship may have a jump drive, giving it the ability to quickly travel a slightly random, far distance. Damamged, it is unavailable.


