package main

/*
------------------------------------------------------------------------
-----------------------------instruction list---------------------------
------------------------------------------------------------------------
*/

type instruction struct {
	name string
	item any
}

type instructionList []instruction

var instructions instructionList

//var varIndex map[any]int
