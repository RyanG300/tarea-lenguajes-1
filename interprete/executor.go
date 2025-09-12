package main

import (
	"fmt"
	"math"
)

// Que instrucción ejecuta
func whichExecute(inst *instructionList) {
	pc := 0
	for pc < len(*inst) {
		ins := (*inst)[pc]
		switch ins.name {
		case "LOAD_CONST":
			EXECUTE_LOAD_CONST(ins.item)
			pc++
		case "STORE_FAST":
			varMem.EXECUTE_STORE_FAST(ins.item)
			pc++
		case "BINARY_MULTIPLY":
			EXECUTE_BINARY_MULTIPLY()
			pc++
		case "LOAD_FAST":
			EXECUTE_LOAD_FAST(ins.item)
			pc++
		case "BINARY_DIVIDE":
			EXECUTE_BINARY_DIVIDE()
			pc++
		case "BINARY_AND":
			EXECUTE_BINARY_AND()
			pc++
		case "BINARY_OR":
			EXECUTE_BINARY_OR()
			pc++
		case "BINARY_MODULO":
			EXECUTE_BINARY_MODULO()
			pc++
		case "BUILD_LIST":
			EXECUTE_BUILD_LIST(ins.item.(int))
			pc++
		case "BINARY_SUBSCR":
			EXECUTE_BINARY_SUBSCR()
			pc++
		case "STORE_SUBSCR":
			EXECUTE_STORE_SUBSCR()
			pc++
		case "LOAD_GLOBAL":
			EXECUTE_LOAD_GLOBAL(ins.item)
			pc++
		case "CALL_FUNCTION":
			EXECUTE_CALL_FUNCTION(ins.item.(int))
			pc++
		case "COMPARE_OP":
			{
				EXECUTE_COMPARE_OP(ins.item.(string))
				pc++
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
			pc++
		case "BINARY_ADD":
			EXECUTE_BINARY_ADD()
			pc++
		case "JUMP_ABSOLUTE":
			pc = EXECUTE_JUMP_ABSOLUTE(ins.item.(int))
		case "JUMP_IF_TRUE":
			pc = EXECUTE_JUMP_IF_TRUE(ins.item.(int), pc)
		case "JUMP_IF_FALSE":
			pc = EXECUTE_JUMP_IF_FALSE(ins.item.(int), pc)
		case "END":
			EXECUTE_END()
			return
		default:
			panic("(Error) Instrucción " + ins.name + " no reconocida")
		}
	}
}

/*------------------------------------------------------------------------
--------------------------Ejecuciones de bytecode-------------------------
------------------------------------------------------------------------*/

// Coloca el valor de la constante en el tope de la pila
func EXECUTE_LOAD_CONST(item any) {
	stack.push(item)
	fmt.Println(stack.items[len(stack.items)-1])
}

// Coloca el valor del contenido de la variable en la pila
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

// Carga en el tope de la pila o el valor de la referencia a la función
func EXECUTE_LOAD_GLOBAL(varname any) {
	if EqualAny(varname, "print") {
		stack.push("print")
		fmt.Println("Función print cargada en el tope de la pila")
		return
	}
	panic("(Error) La función " + fmt.Sprint(varname) + " no está definida")
}

// Realiza un salto a la dirección de código de la función (Solo print jaja)
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
	stack.pop()
	fmt.Print("Output de print: ")
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	fmt.Println("")

}

// Realiza una comparación booleana según el op que reciba
func EXECUTE_COMPARE_OP(op string) {
	if len(stack.items) < 2 {
		panic("(Error) No hay suficientes elementos en la pila para realizar la comparación")
	}
	val2, _ := stack.top()
	stack.pop()
	val1, _ := stack.top()
	stack.pop()
	cmp, err := CompareAny(val1, val2, false) // false -> no fallback arbitrario
	if err != nil {
		panic(err) // o manejar error
	}

	var result bool
	switch op {
	case "<":
		result = cmp < 0
	case "<=":
		result = cmp <= 0
	case ">":
		result = cmp > 0
	case ">=":
		result = cmp >= 0
	case "==":
		result = cmp == 0
	case "!=":
		result = cmp != 0
	default:
		panic("operador no soportado")
	}
	stack.push(result)
	fmt.Println("Resultado de la comparación ", op, ": ", result)
}

// Realiza una suma de operandos
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
	a, _ := stack.top()
	fmt.Println("Resultado de la resta: ", a)
}

// ealiza una suma de operandos
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

	if a, ok := val1.(int); ok {
		if b, ok2 := val2.(int); ok2 {
			stack.push(a + b)
			top, _ := stack.top()
			fmt.Println("Resultado de la suma (int):", top)
			return
		}
	}
	
	n1 := toFloat(val1)
	n2 := toFloat(val2)
	stack.push(n1 + n2)
	a, _ := stack.top()
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

// Realiza un salto incondicional a la instrucción en la posición target
func EXECUTE_JUMP_ABSOLUTE(target int) int {
	if target < 0 || target >= len(instructions) {
		panic("(Error) Target de salto fuera de rango")
	}
	return target
}

// Realiza un salto condicional si el valor en el tope de la pila es verdadero
func EXECUTE_JUMP_IF_TRUE(target, pc int) int {
	val, ok := stack.top()
	if !ok {
		panic("(Error) No hay suficientes elementos en la pila para realizar salto")
	}
	stack.pop()

	boolVal, ok := val.(bool)
	if !ok {
		panic("(Error) El valor en el tope de la pila no es booleano")
	}
	if boolVal {
		if target < 0 || target >= len(instructions) {
			panic("(Error) Target de salto fuera de rango")
		}
		return target
	}
	return pc + 1

}

// Realiza un salto condicional si el valor en el tope de la pila es falso
func EXECUTE_JUMP_IF_FALSE(target, pc int) int {
	val, ok := stack.top()
	if !ok {
		panic("(Error) No hay suficientes elementos en la pila para realizar salto")
	}
	stack.pop()

	boolVal, ok := val.(bool)
	if !ok {
		panic("(Error) El valor en el tope de la pila no es booleano")
	}
	if !boolVal {
		if target < 0 || target >= len(instructions) {
			panic("(Error) Target de salto fuera de rango")
		}
		return target
	}
	return pc + 1
}

// Finaliza la ejecución del programa
func EXECUTE_END() {
	fmt.Println("Programa finalizado")
}
