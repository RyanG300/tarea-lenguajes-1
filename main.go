package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
			if index+1 == len(text) {
				item = convertTextToVariable(itemRune)
				whichExecute(instruccionString, item)
				return
			}
			if text[index+1] == '\r' {
				item = convertTextToVariable(itemRune)
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

func convertTextToVariable(text []rune) any {
	numers := []rune{48, 49, 50, 51, 52, 53, 54, 55, 56, 57}
	whatVariable := 1 // 1=int, 2=float, 3=rune/character, 4=string
	noMoreFloat := false
	negative := true
	for _, val := range text {
		if negative && val == '-' {
			continue
		} else {
			negative = false
		}
		if contains(numers, val) {
			continue
		} else if val == '.' && (text[0] != val || text[len(text)-1] != val) && !noMoreFloat {
			whatVariable = 2
			noMoreFloat = true
		} else if len(text) == 1 {
			whatVariable = 3
			break
		} else {
			whatVariable = 4
			break
		}
	}
	var n any
	switch whatVariable {
	case 1:
		{
			n, _ := strconv.Atoi((string(text)))
			fmt.Println(n)
			return n
		}
	case 2:
		{
			n, _ := strconv.ParseFloat(string(text), 32)
			fmt.Println(n)
			return n
		}
	case 3:
		{
			n = text
			return n
		}
	case 4:
		{
			n = string(text)
			return n
		}
	}
	return n
}

func contains(slice []rune, val rune) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
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
	var txt string = "Pruebas_de_interprete\\example1.txt"
	fileError, errorName := fileExist(txt)
	if fileError {
		fmt.Println("Si va")
	} else {
		fmt.Println("no va xdxdxd: ", errorName)
	}
	dataTemp, _ := readFile(txt)
	data := string(dataTemp)
	fmt.Println(data)
	lecturaByteCode(data)

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
