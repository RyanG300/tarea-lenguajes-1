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

var stack Stack

type variablesMemory struct { //Es lo mismo al anterior struct, estas separadas
	// solo para evitar utilizar los atributos de la pila en la memoria de variables
	nombre   any
	variable any
}

type varMemory []variablesMemory
var varMem varMemory

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
------------------------- Lector ByteCode y más --------------------------
------------------------------------------------------------------------*/

// Lee cada una de las instrucciones del archivo txt
func (inst *instructionList)lecturaByteCode(text string) {
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
				//whichExecute(instruccionString, item)
				*inst = append(*inst, instruction{name: instruccionString, item: item})
				return
			}
			if text[index+1] == '\r' {
				item = convertTextToVariable(itemRune)
				//whichExecute(instruccionString, item)
				*inst = append(*inst, instruction{name: instruccionString, item: item})
				instruccionString = ""
				instruccionRune = []rune{}
				intsRead = false
				itemRune = []rune{}
			}
		} else if chr == '\t' && !readingInts && !intsRead {
			readingInts = true
		} else if chr == '\r' && readingInts {
			instruccionString = string(instruccionRune)
			//whichExecute(instruccionString, item)
			*inst = append(*inst, instruction{name: instruccionString, item: item})
			instruccionString = ""
			instruccionRune = []rune{}
			readingInts = false

		} else if chr == '\t' && readingInts {
			instruccionString = string(instruccionRune)
			readingInts = false
			intsRead = true
		} else if readingInts {
			instruccionRune = append(instruccionRune, chr)
			if index+1 == len(text) {
				instruccionString = string(instruccionRune)
				//whichExecute(instruccionString, item)
				*inst = append(*inst, instruction{name: instruccionString, item: item})
				return
			}
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
			n = string(text)
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

// Compara dos variables de cualquier tipo (Solo string, []int32,float32 e int)
func EqualAny(a, b any) bool {
	// nil handling
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch va := a.(type) {
	case string:
		vb, ok := b.(string)
		if !ok {
			return false
		}
		return va == vb

	case []int32:
		vb, ok := b.([]int32)
		if !ok {
			return false
		}
		// tratar nil y slice vacío como iguales
		if len(va) == 0 && len(vb) == 0 {
			return true
		}
		if len(va) != len(vb) {
			return false
		}
		for i := range va {
			if va[i] != vb[i] {
				return false
			}
		}
		return true

	case float32:
		// b puede ser float32 o int
		switch vb := b.(type) {
		case float32:
			af := float64(va)
			bf := float64(vb)
			if math.IsNaN(af) || math.IsNaN(bf) {
				return false
			}
			return af == bf
		case int:
			af := float64(va)
			bf := float64(vb)
			if math.IsNaN(af) || math.IsNaN(bf) {
				return false
			}
			return af == bf
		default:
			return false
		}

	case int:
		// b puede ser int o float32
		switch vb := b.(type) {
		case int:
			return va == vb
		case float32:
			af := float64(va)
			bf := float64(vb)
			if math.IsNaN(af) || math.IsNaN(bf) {
				return false
			}
			return af == bf
		default:
			return false
		}

	default:
		// no deberían llegar otros tipos
		return false
	}
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
func whichExecute( inst *instructionList) {
	for _, ins := range *inst {
		switch ins.name {
		case "LOAD_CONST":
			EXECUTE_LOAD_CONST(ins.item)
		case "STORE_FAST":
			varMem.EXECUTE_STORE_FAST(ins.item)
		case "BINARY_MULTIPLY":
			EXECUTE_BINARY_MULTIPLY()
		case "LOAD_FAST":
			EXECUTE_LOAD_FAST(ins.item)
		case "BINARY_DIVIDE":
			EXECUTE_BINARY_DIVIDE()
		case "BINARY_AND":
			EXECUTE_BINARY_AND()
		case "BINARY_OR":
			EXECUTE_BINARY_OR()
		case "BINARY_MODULO":
			EXECUTE_BINARY_MODULO()
		case "BUILD_LIST":
			EXECUTE_BUILD_LIST(ins.item.(int))
		case "BINARY_SUBSCR":
			EXECUTE_BINARY_SUBSCR()
		case "STORE_SUBSCR":
			EXECUTE_STORE_SUBSCR()
		case "LOAD_GLOBAL":
			EXECUTE_LOAD_GLOBAL(ins.item)
		case "CALL_FUNCTION":
			EXECUTE_CALL_FUNCTION(ins.item.(int))
		case "COMPARE_OP":{
			EXECUTE_COMPARE_OP(ins.item.(string))
			/*switch ins.item.(type) {
				
				case []int32:
					var itemStr string
					for _, val := range ins.item.([]int32) {
						itemStr += strconv.Itoa(int(val))
					}
					EXECUTE_COMPARE_OP(itemStr)
				case string:
					EXECUTE_COMPARE_OP(ins.item.(string))	
				default:
					panic("(Error) Operador de comparación no reconocido")
				}*/
		}
		case "BINARY_SUBSTRACT":
			EXECUTE_BINARY_SUBSTRACT()
		case "BINARY_ADD":
			EXECUTE_BINARY_ADD()
		default:
			panic("(Error) Instrucción " + ins.name + " no reconocida")
		}
	}
}

/*------------------------------------------------------------------------
----------------------- Instrucciones de bytecode ------------------------
------------------------------------------------------------------------*/

// Coloca el valor de la constante en el tope de la pila
func EXECUTE_LOAD_CONST(item any) {
	stack.push(item)
	fmt.Println(stack.items[len(stack.items)-1])
}

//Coloca el valor del contenido de la variable en la pila
func EXECUTE_LOAD_FAST(varname any) {
	for i := len(varMem) - 1; i >= 0; i-- {
		if EqualAny(varMem[i].nombre, varname) {
			stack.push(varMem[i].variable)
			fmt.Println("Variable ", varMem[i].nombre, " cargada con el valor de: ", varMem[i].variable, " en el tope de la pila")
			return
		}
	}
	panic("(Error) La variable " + fmt.Sprint(varname) + " no está definida")
}

//Carga en el tope de la pila o el valor de la referencia a la función
func EXECUTE_LOAD_GLOBAL(varname any) {
	if(EqualAny(varname, "print")){
		stack.push("print")
		fmt.Println("Función print cargada en el tope de la pila")
		return
	}
	panic("(Error) La función " + fmt.Sprint(varname) + " no está definida")
}

//Realiza un salto a la dirección de código de la función (Solo print jaja)
func EXECUTE_CALL_FUNCTION(numArgs int) {
	if len(stack.items) < numArgs+1 {
		panic("(Error) No hay suficientes elementos en la pila para realizar la llamada a función")
	}
	args := make([]any, numArgs)
	for i := numArgs - 1; i >= 0; i-- {
		arg, _ := stack.top()
		stack.pop()
		args[i] = arg
	}
	fmt.Print("Output de print: ")
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	
}

//Realiza una comparación booleana según el op que reciba
func EXECUTE_COMPARE_OP(op string) {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar la comparación")
	}
	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()
	var result bool
	switch op {
	case "==":
		result = EqualAny(val1, val2)
	case "!=":
		result = !EqualAny(val1, val2)
	case "<":
		result = fmt.Sprint(val1) < fmt.Sprint(val2)
	case "<=":
		result = fmt.Sprint(val1) <= fmt.Sprint(val2)
	case ">":
		result = fmt.Sprint(val1) > fmt.Sprint(val2)
	case ">=":
		result = fmt.Sprint(val1) >= fmt.Sprint(val2)
	default:
		panic("(Error) Operador de comparación no soportado")
	}
	stack.push(result)
	fmt.Println("Resultado de la comparación ", op, ": ", result)
}

//Realiza una suma de operandos
func EXECUTE_BINARY_SUBSTRACT() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar resta")
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
			panic("(Error) Tipo no soportado en resta")
		}
	}
	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()
	n1 := toFloat(val1)
	n2 := toFloat(val2)
	stack.push(n1 - n2)
	a,_ := stack.top()
	fmt.Println("Resultado de la resta: ", a)
}

//ealiza una suma de operandos
func EXECUTE_BINARY_ADD() {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar suma")
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
			panic("(Error) Tipo no soportado en suma")
		}
	}
	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()
	n1 := toFloat(val1)
	n2 := toFloat(val2)
	stack.push(n1 + n2)
	a,_ := stack.top()
	fmt.Println("Resultado de la suma: ", a)
}

// Escribe el contenido del tope de la pila en la variable
func (varMe *varMemory) EXECUTE_STORE_FAST(varname any) {
	//var index = len(*varMe)-1
	//varIndex[varname] = index
	variableItem, top := stack.top()
	if !top {
		panic("(Error) No se puede almacenar variables sin un tope en la pila")
	}
	for i := len(*varMe) - 1; i >= 0; i-- {
		if EqualAny((*varMe)[i].nombre, varname) {
			(*varMe)[i].variable = variableItem
			stack.pop()
			fmt.Println("Variable ", (*varMe)[i].nombre, " actualizada con el valor de: ", (*varMe)[i].variable)
			return
		}
	}
	//Si no existe la variable, se crea
	*varMe = append(*varMe, variablesMemory{nombre: varname, variable: variableItem})
	stack.pop()
	fmt.Println("Variable ", (*varMe)[len(*varMe)-1].nombre, " guardada con el valor de: ", (*varMe)[len(*varMe)-1].variable)
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
