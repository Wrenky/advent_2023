namespace import ::tcl::mathop::*
namespace import ::tcl::mathfunc::*
source "../helpers/utils.tcl"


#####################################################################
# Parse the input into a id - sets dictionary
#####################################################################
set data [readFile [lindex $argv 0]]
set linePattern {Game\s(\d+):\s(.*?)$}
set setPattern {(\d+)\s(blue|red|green)}
foreach line [split $data "\n"] {
    lassign [regexp -all -inline -- $linePattern $line] -> id gameInfo
    dict set games $id [lmap x [split $gameInfo ";"] {
        lmap y [split $x ","] {split [string trim $y]}
    }]
}

#####################################################################
# Part 1: Figure out what games are valid, sum the ids
#####################################################################
# These are my filter procs
proc isValid {s} {
    set colors [dict create "red" 12 "green" 13 "blue" 14]
    lassign [split $s] amount color
    return [expr {$amount < [dict get $colors $color]}]
}
proc processSet {s} {
    return [expr ([llength $s] == [llength [filter isValid $s]])]
}

# filter for only valid games 
set validGames [dict filter $games script {id sets} {
   expr [llength $sets] == [llength [filter processSet $sets]]
}]

# validGames/games are keyed by id, so just sum the keys
puts "Part 1: [fold 0 + [dict keys $validGames]]"



#####################################################################
# Part 2: Get the minimum color needed for each game, multiply them
#         then sum all games
#####################################################################
proc minimumRGB {sets} {
    set colors [dict create "red" 0 "green" 0 "blue" 0]
    foreach s [join $sets] {
        lassign [split $s] amount color 
        if {[dict get $colors $color] < $amount} {
            dict set colors $color $amount
        }
    }
    return [dict values $colors]
}

# Multiply the minimum RGB values for each game
set gamePower [lmap {id sets} $games {
    fold 1 * [minimumRGB $sets]
}]
#####################################################################

# Sum the powers
puts "Part 2: [fold 0 + $gamePower]"
