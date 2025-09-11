package main

import (
	"errors"
	"math"
	"os"
	"strconv"
)

/*------------------------------------------------------------------------
--------------------------Funciones auxiliares----------------------------
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
