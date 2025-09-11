package main

import (
	"fmt"
)

/*------------------------------------------------------------------------
------------------------ Main Main Main Main Main ------------------------
------------------------------------------------------------------------*/

func main() {
	var txt string = "Pruebas_de_interprete/example1.txt"
	fileError, errorName := fileExist(txt)
	if fileError {
		fmt.Println("Si va")
	} else {
		fmt.Println("no va xdxdxd: ", errorName)
	}
	//varIndex = make(map[any]int)
	dataTemp, _ := readFile(txt)
	data := string(dataTemp)
	fmt.Println(data)
	instructions.lecturaByteCode(data)
	whichExecute(&instructions)
	//
	/*var tal string = "dsds"
	fmt.Println(any(tal))
	/*tal := '0'
	fmt.Println(tal)*/

}

/*Prueba

var s Stack
	s.Push(10)
	s.Push("AAAAA")
	if s.IsEmpty() {
		fmt.Println("Pingas")
	}
	fmt.Println(s.items[1])
	s.items[1] = 20

	var e any
	var y bool
	e,y = s.Top()
	if !y{
		fmt.Println("Pingas2")
	}else{
		fmt.Println(e)
	}

*/
