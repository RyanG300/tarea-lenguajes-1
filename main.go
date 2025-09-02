package main

import (
	"errors"
	"fmt"
	"os"
)

/*
------------------------------------------------------------------------
------------------Pila y memoria de las variables (struct)----------------
------------------------------------------------------------------------
*/
type Stack struct {
	items []any
}

type variablesMemory struct { //Es lo mismo al anterior struct, estas separadas solo para evitar utilizar los atributos de la pila en la memoria de variables
	variables []any
}

var stack Stack

/*------------------------------------------------------------------------
---------------------------Atributos de la pila---------------------------
------------------------------------------------------------------------*/

func (s *Stack) push(v any) {
	s.items = append(s.items, v)
}

func (s *Stack) pop() (any, bool) {
	if len(s.items) == 0 {
		return nil, false
	}
	i := len(s.items) - 1
	v := s.items[i]
	s.items = s.items[:i]
	return v, true
}

func (s *Stack) top() (any, bool) {
	var x any
	if len(s.items) == 0 {
		return x, false
	}
	x = s.items[len(s.items)-1]
	return x, true
}

func (s *Stack) isEmpty() bool {
	return len(s.items) == 0
}

/*------------------------------------------------------------------------
-------------------------- Lectura de archivos ---------------------------
------------------------------------------------------------------------*/

func fileExist(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

/*------------------------------------------------------------------------
----------------------- Instrucciones de bytecode ------------------------
------------------------------------------------------------------------*/

func lecturaByteCode(text string) {
	readingInts := false
	intsRead := false 
	var instruccionRune []rune
	var instruccionString string
	var itemRune []rune
	var item any
	for index, chr := range text {
		if chr != '\t' && !readingInts && intsRead {
			itemRune = append(itemRune, chr)
			if index+1 == len(text){
				item = any(itemRune)
				whichExecute(instruccionString, item)
				return
			} 
			if text[index+1] == '\r' {
				item = any(itemRune)
				whichExecute(instruccionString, item)
				instruccionString = ""
				instruccionRune = []rune{}
				intsRead = false
			}
		} else if chr == '\t' && !readingInts && !intsRead {
			readingInts = true
		} else if chr == '\r' && readingInts {
			instruccionString = string(instruccionRune)
			whichExecute(instruccionString, item)
			instruccionString = ""
			instruccionRune = []rune{}
			readingInts = false
		} else if chr == '\t' && readingInts {
			instruccionString = string(instruccionRune)
			readingInts = false
			intsRead = true
		} else if readingInts {
			instruccionRune = append(instruccionRune, chr)
		}
	}
}

func whichExecute(instruccion string, item any) {
	switch instruccion {
	case "LOAD_CONST":
		EXECUTE_LOAD_CONST(item)
	case "STORE_FAST":
		fmt.Println("Something happened 2...")
	}
}

func EXECUTE_LOAD_CONST(item any) {
	stack.push(item)
	fmt.Println(stack.items[0])
}

func main() {
	/*var txt string = "Pruebas_de_interprete\\example1.txt"
	fileError, errorName := fileExist(txt)
	if fileError {
		fmt.Println("Si va")
	} else {
		fmt.Println("no va xdxdxd: ", errorName)
	}
	dataTemp, _ := readFile(txt)
	data := string(dataTemp)
	fmt.Println(data)
	lecturaByteCode(data)*/
	tal:= 'a'
	e := int(tal)
	fmt.Println(e)
	//
	/*var tal string = "dsds"
	fmt.Println(any(tal))*/
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
