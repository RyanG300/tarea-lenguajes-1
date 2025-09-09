package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

/*
------------------------------------------------------------------------
------------------Pila y memoria de las variables (struct)----------------
------------------------------------------------------------------------
*/

// La pila, tiene los metodos push, pop, top, isEmpty
type Stack struct {
	items []any
}

type variablesMemory struct { //Es lo mismo al anterior struct, estas separadas
	// solo para evitar utilizar los atributos de la pila en la memoria de variables
	nombre   any
	variable any
}

var stack Stack

type varMemory []variablesMemory

var varMem varMemory

//var varIndex map[any]int

/*------------------------------------------------------------------------
---------------------------Atributos de la pila---------------------------
------------------------------------------------------------------------*/

func (s *Stack) push(v any) {
	s.items = append(s.items, v)
}

func (s *Stack) pop() bool {
	if len(s.items) == 0 {
		return false
	}
	i := len(s.items) - 1
	//v := s.items[i]
	s.items = s.items[:i]
	return true
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

// Lee cada una de las instrucciones del archivo txt
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
				itemRune = []rune{}
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

// Convierte la variable del archivo txt en string, rune, int, float32 (La lista es de otra forma)
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
			//fmt.Println(n)
			return n
		}
	case 2:
		{
			n, _ := strconv.ParseFloat(string(text), 32)
			//fmt.Println(n)
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

// Comprueba si existe val en un slice cualquiera
func contains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// Que instrucción ejecuta
func whichExecute(instruccion string, item any) {
	switch instruccion {
	case "LOAD_CONST":
		EXECUTE_LOAD_CONST(item)
	case "STORE_FAST":
		varMem.EXECUTE_STORE_FAST(item)
	case "BINARY_MULTIPLY":
		EXECUTE_BINARY_MULTIPLY()

	}
}

// Coloca el valor de la constante en el tope de la pila
func EXECUTE_LOAD_CONST(item any) {
	stack.push(item)
	//fmt.Println(stack.items[0])
}

// Escribe el contenido del tope de la pila en la variable
func (varMe *varMemory) EXECUTE_STORE_FAST(varname any) {
	//var index = len(*varMe)-1
	//varIndex[varname] = index
	variableItem, top := stack.top()
	if !top {
		panic("(Error) No se puede almacenar variables sin un tope en la pila")
	}
	*varMe = append(*varMe, variablesMemory{nombre: varname, variable: variableItem})
	stack.pop()
}

// Realiza la multiplicación de los dos valores en el tope de la pila
func EXECUTE_BINARY_MULTIPLY() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar multiplicación")
	}

	toFloat := func(val any) float64 {
		switch v := val.(type) {
		case int:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return v
		default:
			panic("(Error) Tipo no soportado en multiplicación")
		}
	}

	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()

	n1 := toFloat(val1)
	n2 := toFloat(val2)
	stack.push(n1 * n2)
}

// Realiza la división de los dos valores en el tope de la pila
func EXECUTE_BINARY_DIVIDE() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar división")
	}

	toFloat := func(val any) float64 {
		switch v := val.(type) {
		case int:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return v
		default:
			panic("(Error) Tipo no soportado en división")
		}
	}
	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()

	n1 := toFloat(val1)
	n2 := toFloat(val2)

	if n2 == 0 {
		panic("(Error) No se puede dividir entre 0")
	}

	stack.push(n1 / n2)
}

// Realiza la operación AND entre los dos valores en el tope de la pila
func EXECUTE_BINARY_AND() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar AND")
	}

	toBool := func(val interface{}) bool {
		b, ok := val.(bool)
		if !ok {
			panic("(Error) Tipo no soportado en AND")
		}
		return b
	}

	val1, _ := stack.top()
	stack.pop()
	val2, _ := stack.top()
	stack.pop()

	stack.push(toBool(val1) && toBool(val2))
}

// Realiza la operación OR entre los dos valores en el tope de la pila
func EXECUTE_BINARY_OR() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar OR")
	}

	toBool := func(val interface{}) bool {
		b, ok := val.(bool)
		if !ok {
			panic("(Error) Tipo no soportado en OR")
		}
		return b
	}

	val1, _ := stack.top()
	stack.pop()
	val2, _ := stack.top()
	stack.pop()

	stack.push(toBool(val1) || toBool(val2))
}

// Realiza la operación MODULO entre los dos valores en el tope de la pila
func EXECUTE_BINARY_MODULO() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar modulo")
	}

	toFloat := func(val any) float64 {
		switch v := val.(type) {
		case int:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return v
		default:
			panic("(Error) Tipo no soportado en modulo")
		}
	}
	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()

	n1 := toFloat(val1)
	n2 := toFloat(val2)

	if n2 == 0 {
		panic("(Error) No se puede dividir entre 0")
	}

	stack.push(math.Mod(n1, n2))
}

// Construye una lista con los n elementos del tope de la pila
func EXECUTE_BUILD_LIST(elements int) {
	if len(stack.items) < elements {
		panic("(Error) No hay suficientes elementos en la pila para construir la lista")
	}
	list := make([]any, elements)
	for i := 1; i <= elements; i++ {
		val, _ := stack.top()
		stack.pop()
		list[elements-i] = val
	}
	stack.push(list)
}

// Realiza la subscripción de una lista
func EXECUTE_BINARY_SUBSCR() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar subscripción")
	}
	index, _ := stack.top()
	stack.pop()
	list, _ := stack.top()
	stack.pop()

	listSlice, ok := list.([]any)
	if !ok {
		panic("(Error) El elemento no es una lista")
	}

	idx, ok := index.(int)
	if !ok {
		panic("(Error) El índice no es un entero")
	}

	if idx < 0 || idx >= len(listSlice) {
		panic("(Error) Índice fuera de rango")
	}

	stack.push(listSlice[idx])
}

// Realiza la asignación por subscripción de una lista
func EXECUTE_STORE_SUBSCR() {
	if len(stack.items) < 3 {
		panic("(Error) No hay suficientes elementos en la pila para realizar asignación por subscripción")
	}
	index, _ := stack.top()
	stack.pop()
	list, _ := stack.top()
	stack.pop()
	value, _ := stack.top()
	stack.pop()

	listSlice, ok := list.([]any)
	if !ok {
		panic("(Error) El elemento no es una lista")
	}

	idx, ok := index.(int)
	if !ok {
		panic("(Error) El índice no es un entero")
	}

	if idx < 0 || idx >= len(listSlice) {
		panic("(Error) Índice fuera de rango")
	}
	listSlice[idx] = value
}

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
