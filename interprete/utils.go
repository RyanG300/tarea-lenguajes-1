package main

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
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

// toRat intenta convertir v a *big.Rat. Devuelve (rat, true) si pudo.
func toRat(v interface{}) (*big.Rat, bool) {
	if v == nil {
		return nil, false
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return new(big.Rat).SetInt64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		u := rv.Uint()
		bi := new(big.Int).SetUint64(u)
		return new(big.Rat).SetInt(bi), true
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		r := new(big.Rat)
		r.SetFloat64(f) // convierte el float a rational (puede perder exactitud para floats)
		return r, true
	case reflect.String:
		s := rv.String()
		r := new(big.Rat)
		if _, ok := r.SetString(s); ok { // acepta "10", "3.14", "-5/2"...
			return r, true
		}
		return nil, false
	default:
		// Intenta parsear la representación textual como número (por ejemplo si es tipo personalizado)
		s := fmt.Sprint(v)
		r := new(big.Rat)
		if _, ok := r.SetString(s); ok {
			return r, true
		}
		return nil, false
	}
}

// CompareAny compara a y b y devuelve:
//
//	-1 si a < b
//	 0 si a == b
//	+1 si a > b
//
// y error si no se pueden comparar (a menos que allowDeterministicFallback=true).
func CompareAny(a, b interface{}, allowDeterministicFallback bool) (int, error) {
	// igualdad rápida
	if reflect.DeepEqual(a, b) {
		return 0, nil
	}

	// intento numérico: si ambos son convertibles a big.Rat, comparar numéricamente
	if ra, oka := toRat(a); oka {
		if rb, okb := toRat(b); okb {
			return ra.Cmp(rb), nil // -1,0,1
		}
	}

	// bool
	if ab, oka := a.(bool); oka {
		if bb, okb := b.(bool); okb {
			if ab == bb {
				return 0, nil
			}
			if !ab && bb {
				return -1, nil // false < true
			}
			return 1, nil
		}
	}

	// string
	if as, oka := a.(string); oka {
		if bs, okb := b.(string); okb {
			return strings.Compare(as, bs), nil // -1,0,1
		}
	}

	// no comparable semántico: fallback determinístico o error
	if allowDeterministicFallback {
		sa := fmt.Sprintf("%T|%v", a, a)
		sb := fmt.Sprintf("%T|%v", b, b)
		return strings.Compare(sa, sb), nil
	}

	return 0, fmt.Errorf("no se pueden comparar los tipos %T y %T", a, b)
}
