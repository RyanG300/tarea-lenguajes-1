package main

/*------------------------------------------------------------------------
------------------------- Lector ByteCode y más --------------------------
------------------------------------------------------------------------*/

// Lee cada una de las instrucciones del archivo txt
func (inst *instructionList) lecturaByteCode(text string) {
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
