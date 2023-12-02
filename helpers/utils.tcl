namespace import ::tcl::mathop::*
package require textutil
proc fold {res op l} {foreach e $l {set res [$op $res $e]}; set res}
proc filter {cond l} {set res {}; foreach e $l {if [{*}$cond $e] {lappend res $e}}; set res}

proc lpop {listVar {count 1}} {
	upvar 1 $listVar l
	set r [lrange $l 0 [incr count -1]]
	set l [lreplace $l 0 $count]
	return $r
}



#############################################
# Amazing snippets from
# http://chiselapp.com/user/aspect/repository/tcl-hacks/file?name=modules/fun-0.tm
#############################################
proc range {a {b ""}} {
    if {$b eq ""} {
        set b $a
        set a 0
    }
    for {set r {}} {$a<$b} {incr a} {
        lappend r $a
    }
    return $r
}
proc lolcat args {join [uplevel lmap $args]}
proc dictify {cmdPrefix ls} {lolcat x $ls {list $x [uplevel 1 $cmdPrefix [list $x]]}}
proc pdict {d} {array set {} $d; parray {}}
proc lconcat args {concat {*}[uplevel 1 lmap $args]}
proc map {cmdPrefix args} {
    set names [range [llength $args]]
    set forArgs [lconcat n $names a $args {list $n $a}]
    set cmdArgs [lconcat name $names {string cat \$ $name}]
    set body "uplevel 1 [list $cmdPrefix] \[list $cmdArgs\]"
    lmap {*}$forArgs $body
}
proc indexed {list} {
    set i -1
    concat {*}[lmap {x} $list {
        list [incr i] $x
    }]
}
proc lgroup {listIn lengthOfSublist} {
    set i 0
    foreach it $listIn {
        lappend tmp $it
        if {[llength $tmp] == $lengthOfSublist} {
            lappend result $tmp
            set tmp {}
        }
        incr i
    }
    if {[llength $listIn] % $lengthOfSublist} {
        lappend result $tmp
    }
    return $result
}
proc repeat {n script args} {
    set script [concat $script {*}$args]
    set res {}
    loop i 0 $n {
        lappend res [uplevel 1 $script]
    }
    set res
}
proc cmdpipe args {
    set anonvar ~
    set args [lassign $args body]
    foreach cmd $args {
        if {[string first $anonvar $cmd] >= 0} {
            set body [string map [list $anonvar "\[$body\]"] $cmd]
        } else {
            set body "$cmd \[$body\]"
        }
    }
    set body
}
## Make min/max respect lists
proc max {args} {tailcall ::tcl::mathfunc::max {*}[concat {*}$args]}
proc min {args} {tailcall ::tcl::mathfunc::min {*}[concat {*}$args]}
proc maxlen {ss} {return [max {*}[map {string length} $ss]]}
#############################################

#############################################
## Helpers
#############################################
proc readFile {filename} {
    if {$filename eq ""} {return ""}
	set fh [open $filename]
	set data [read -nonewline $fh]
	close $fh
	return $data
}

#############################################
# Matrices! Operate on a list of lists like a 
#    matrix. Each destructive action/update
#    return a new matrix.
#############################################
namespace eval ::matrix {
    # Setters/getters
    proc setCell {m x y val} {return [lreplace $m $y $y [lreplace [lindex $m $y] $x $x $val]]}
    proc getCell {m x y} {return [lindex [lindex $m $y] $x]}
    proc transpose {m} {
        set m1 [list]
        for {set i 0} {$i < [llength [lindex $m 0]]} {incr i} {
            lappend m1 [lmap r $m {lindex $r $i}]
        }
        return $m1
    }
    proc print {m} {foreach r $m {puts $r}}
    proc size {m} {
        return [list [llength [lindex $m 0]] [llength $m]]
    }
}
#############################################

