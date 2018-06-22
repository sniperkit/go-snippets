**Project's new home: https://github.com/mattcan/csgodem**

--

## Counter-Strike:Global Offensive Demo Parser

The goals of this project are to:

* port Valve's demoinfo application to Go
* create a faster parser than the one provided by Valve
* provide an increased number of profiles for specific data (all kills, all
  deaths, etc)
* have the option to save directly to a database or CSV

## Demos

1. simple_d2_t_01 - *20.9s and 1331 frames*
  * As a terrorist, spawn facing down mid
	* Buy Tec-9
	* Jump out of spawn and knife run down mid to double doors
	* Fire 4 shots straight ahead
2. simple_d2_t_02 - *16.9s and 1080 frames*
  * From double doors, knife run into CT spawn
	* Turn 180 degrees
	* Fire four shots and B doors
	* Reload
3. simple_d2_t_03 - *72s and 4595 frames*
  * From spawn, B side, knife run through tunnels to B
	* Plant bomb
	* Wait for explosion
