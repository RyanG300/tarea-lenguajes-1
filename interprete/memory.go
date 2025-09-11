package main

/*
------------------------------------------------------------------------
---------------------------------Memoria--------------------------------
------------------------------------------------------------------------
*/
type variablesMemory struct { //Es lo mismo al anterior struct, estas separadas
	// solo para evitar utilizar los atributos de la pila en la memoria de variables
	nombre   any
	variable any
}

type varMemory []variablesMemory

var varMem varMemory
