package pascomp

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/kdgwill/golang_dev/JAPC_WIG/pascomp/datastructures"
)

type SymbolTable struct {
	// The string table is a long string in which all lexemes are stored
	stringtable [stringTableSize]rune
	// A series of indices pointing to the lexeme within the string table
	//	 as well as to the relevant attribute table entry
	nametable [nameTableSize]nameTabType
	// The attribute table entries
	attribTable [attribTableSize]attribTabType
	hashTable   [hashTableSize]int
	// The lengths of the string table, name table and attribute table
	strTabLen, namTabLen, attribTabLen, auxTabLen int

	thisproc  procstackitem        // A stack entry for the current procedure
	procStack datastructures.Stack //<procstackitem> // The procedure stack
}

func NewSymbolTable() *SymbolTable {
	var nameindex int

	st := new(SymbolTable)
	//initialize the first entry for the procedure stack
	st.thisproc = st.initprocentry(-1)

	//Initialize the hash table, the name table's next
	//field and the attribute table's fields as -1
	for i := 0; i < hashTableSize; i++ {
		st.hashTable[i] = -1
	}

	//	The attribute table's fields are all initially
	//	unknown, 0 or -1 (if they're indices).
	for i := 0; i < attribTableSize; i++ {
		st.attribTable[i] = attribTabType{
			smtype:          Stunknown,
			tok_class:       Tokunknown,
			dataclass:       Dtunknown,
			owningprocedure: -1,
			thisname:        -1,
			outerscope:      -1,
			scopenext:       -1,
			value: valRec{
				tag: tint,
				val: 0}} //, //ival
		//label: 0}
		st.attribTable[i].label[0] = 0
	}

	//Install the keywords and operators in the name table and
	//Set their attribute to keyword
	var i int
	for i = 0; i < numTokens; i++ {
		st.Installname(keywords[i], &nameindex)

		var sym SemanticType = Stkeyword
		if i >= numKeywords {
			//now load operators
			sym = Stoperator
		}
		st.Setattrib(nameindex, sym, TokenType(i))
	}

	// Initialize the entry for float, the routine
	// that converts values from integer to real
	st.Installname(keywords[i], &nameindex)
	st.Setattrib(nameindex, Stfunction, TokenType(i))
	st.Installdatatype(nameindex, Stfunction, Dtreal)

	return st
}

// InstallName() - Check if the name is already in the table.
// If not add it to the name table and create
// an attribute table entry.
func (this *SymbolTable) Installname(name string, tabindex *int) bool {
	//TODO:Quick Fix pertaining to proper token identification
	name = strings.ToUpper(name)

	var code, nameindex int

	// Use the function ispresent to see if the token string
	// is in the table.  If so, return a pointer to its
	// attribute table entry.

	length := len(name)
	if this.ispresent(name, &code, &nameindex) {
		if this.nametable[nameindex].symtabptr == -1 {
			*tabindex = this.Installattrib(nameindex)
			return false
		} else {
			*tabindex = this.nametable[nameindex].symtabptr
			return true
		}
	}

	// If not create entries in the name table, copy the name
	// into the string table and create a hash table entry
	// (linking it to its previous entry if necessary) and
	// create an entry in the attribute table with the
	// bare essentials.
	nameindex = this.namTabLen
	this.namTabLen++
	this.nametable[nameindex].strstart = this.strTabLen

	for i := 0; i < length; i++ {
		this.stringtable[this.strTabLen] = rune(name[i])
		this.strTabLen++
	}

	this.stringtable[this.strTabLen] = 0
	this.strTabLen++
	this.nametable[nameindex].nextname = this.hashTable[code]
	this.hashTable[code] = nameindex
	*tabindex = this.Installattrib(nameindex)
	return false

}

// ispresent() - after finding the hash value, it traces
//		 through the hash list, link by link looking to see
//		 if the current token string is there.
//		 this version is private and is intended for use
//		 by Installname
func (this *SymbolTable) ispresent(name string, code *int, nameIndex *int) (found bool) {

	found = false
	var oldnameindex, k int

	/* initialize the old name's index to -1;
	   it may not be there */
	oldnameindex = -1

	// find the hash value
	*code = this.hashcode(name)

	j, length := 0, len(name)
	// starting with the entry in the hash table, trace through
	// the name table's link list for that hash value.
	for *nameIndex = this.hashTable[*code]; !found && *nameIndex != -1; oldnameindex, *nameIndex = *nameIndex, this.nametable[*nameIndex].nextname {
		//k = this.nametable[*nameIndex].strstart
		//found = name != (this.stringtable + string(k))
		for j, k = 0, this.nametable[*nameIndex].strstart; j < length &&
			rune(name[j]) == this.stringtable[k]; j, k = j+1, k+1 {
		}
		if j == length {
			found = true
		}
	}

	// if it's there, we actually went right past it.
	if found {
		*nameIndex = oldnameindex
	}
	return //using found directly do not need to return;
}

// IsPresent() -	After finding the hash value, it traces
//					through the hash list, link by link looking to see
//					if the current token string is there.
//					This version is public and can be used to
//					determine if the lexeme is in the symbol table
//					without Installing it and also returns the index
//					within the attribute table.
func (this *SymbolTable) IsPresent(name string, tabIndex *int) (found bool) {

	found = false
	var nameindex int

	// Initialize the old name's index to -1;
	//	it may not be there
	oldnameindex := -1

	// Find the hash value
	code := this.hashcode(name)

	j, k, length := 0, 0, len(name)
	// Starting with the entry in the hash table, trace through
	// the name table's link list for that hash value.
	for nameindex = this.hashTable[code]; !found && nameindex != -1; oldnameindex, nameindex =
		nameindex, this.nametable[nameindex].nextname {
		//k := this.nametable[nameindex].strstart
		//found = name != (this.stringtable + string(k))
		for j, k = 0, this.nametable[nameindex].strstart; j < length &&
			rune(name[j]) == this.stringtable[k]; j, k = j+1, k+1 {
		}
		if j == length {
			found = true
		}
	}

	// If it's there, we actually went right past it.
	if found {
		nameindex = oldnameindex
	}
	*tabIndex = this.nametable[nameindex].symtabptr
	return
}

// HashCode() -	A hashing function which uses the characters
//				from the end of the token string.  The algorithm comes
//				from Matthew Smosna of NYU.
func (this *SymbolTable) hashcode(name string) int {

	length := len(name)
	//	The number of shifts cannot exceed the runes in an integers
	//	less 8; any more and bits within a given character will be
	//	lost.
	var numshifts int = length
	var temp int = 8*int(unsafe.Sizeof(numshifts)) - 8 //cannot check type so use numshifts in place of int
	if temp < numshifts {
		numshifts = temp
	}

	startchar := ((length - numshifts) % 2)
	var code uint = 0

	//	Left shift one place and add the current character's ASCII
	//	value to the total.
	for i := startchar; i <= startchar+numshifts-1; i++ {
		code = (code << 1) + uint(name[i])
	}

	//	Divide by the table size and use the remainder as the hash
	//	value.
	return int(code) % hashTableSize
}

// InstallAttrib() -	Create a new entry in the attribute
//						table and have this name table
//						entry point to it.
func (this *SymbolTable) Installattrib(nameindex int) int {

	var tabindex int = this.attribTabLen
	this.nametable[nameindex].symtabptr = tabindex
	this.attribTabLen++
	this.attribTable[tabindex].thisname = nameindex
	this.attribTable[tabindex].smtype = Stunknown
	this.attribTable[tabindex].dataclass = Dtunknown

	// Return the index of the attribute table entry
	return tabindex
}

// SetAttrib() -	Set attribute table information, given
//					a pointer to the correct entry in the table.
func (this *SymbolTable) Setattrib(tabindex int, symbol SemanticType, token TokenType) {

	//	Install semantic type and token class
	this.attribTable[tabindex].smtype = symbol
	this.attribTable[tabindex].tok_class = token

	//	Reserved words and operators do not need data types
	if this.attribTable[tabindex].smtype == Stkeyword || this.attribTable[tabindex].smtype == Stoperator {
		this.attribTable[tabindex].dataclass = Dtnone
	} else {
		//	Other symbols have not yet had their
		//	data types determined
		this.attribTable[tabindex].dataclass = Dtunknown
	}

	//	If it's an identifier and it isn't global
	if this.gettok_class(tabindex) == Tokidentifier && this.thisproc.proc != -1 {
		//	If no other scope has a variable with this name
		//	connect its listing to other identifiers in
		//	this scope
		if this.thisproc.sstart == -1 {
			this.thisproc.sstart = tabindex
			this.thisproc.snext = tabindex
		} else {
			//	Otherwise, connect it to the attribute table
			//	entries for this name in outer scopes
			this.attribTable[this.thisproc.snext].scopenext = tabindex
			this.thisproc.snext = tabindex
		}
	}
}

// InstallDataType() -	Install basic data type information,
//						i.e., the data type and semantic type
func (this *SymbolTable) Installdatatype(tabindex int, stype SemanticType, dclass DataType) {
	this.attribTable[tabindex].smtype = stype
	this.attribTable[tabindex].dataclass = dclass

}

// OpenScope() -	Open a new scope for this identifier
func (this *SymbolTable) openscope(tabindex int) int {

	var newtabindex, nameindex int

	// Get the index in the name table
	nameindex = this.attribTable[tabindex].thisname
	//	Create a new attribute table entry and
	//	initialize its information
	newtabindex = this.Installattrib(nameindex)
	this.Setattrib(newtabindex, Stunknown, Tokidentifier)
	// Have this entry point to the outer scope's entry
	this.attribTable[newtabindex].outerscope = tabindex
	return newtabindex
}

// CloseScope() -	Close the scope for ALL the
//					identifiers for the scope
func (this *SymbolTable) closescope() {

	var nmptr, symptr int

	//	Start at the first identifier that belongs to the
	//	procedure and for each identifier
	for symptr = this.thisproc.sstart; symptr != -1; symptr = this.attribTable[symptr].scopenext {
		// Have it point to the outer scope's
		// attribute table entry
		nmptr = this.attribTable[symptr].thisname
		this.nametable[nmptr].symtabptr = this.attribTable[symptr].outerscope
	}

}

// SetProc() -	Set the identifier's owning procedure
func (this *SymbolTable) Setproc(thisproc int, tabindex int) {
	this.attribTable[tabindex].owningprocedure = thisproc
}

// GetProc() -	Returns the identifier's owning procedure
func (this *SymbolTable) Getproc(tabindex int) int {
	return this.attribTable[tabindex].owningprocedure
}

// SetValue() -	Set the value for a real identifier
func (this *SymbolTable) SetFvalue(tabindex int, val float32) {
	this.attribTable[tabindex].value.tag = treal
	this.attribTable[tabindex].value.val = val
}

// SetValue() -	Set the value for an integer identifier
func (this *SymbolTable) SetIvalue(tabindex, val int) {
	this.attribTable[tabindex].value.tag = tint
	this.attribTable[tabindex].value.val = val
}

//	getlabel() -		Gets a label which is used by the final code
//						generator.  If the label is not Installed in the
//						symbol table, it creates one and returns it.
func (this *SymbolTable) Getlabel(tabindex int, varlabel []rune) {
	if this.attribTable[tabindex].label[0] != 0 {
		copy(varlabel, this.attribTable[tabindex].label[:])
	} else {
		this.makelabel(tabindex, &varlabel)
	}
}

// makelabel() -	Makes a label which is used by the final code
//					generator and Installs it in the symbol table.
func (this *SymbolTable) makelabel(tabindex int, label *[]rune) {

	var ivalue int
	var indexstr string // [5]rune

	*label = make([]rune, labelSize)

	switch this.Getsmclass(tabindex) {
	case Stliteral:
		if this.Getdatatype(tabindex) == Dtinteger {
			ivalue = this.attribTable[tabindex].thisname
			*label = append(this.stringtable[:], []rune{rune(this.nametable[ivalue].strstart), rune(0)}...)
			break
		}
		fallthrough

	case Sttempvar:
		copy(*label, []rune("_t")[:])
		//strconv.Itoa(123)
		indexstr = strconv.FormatInt(int64(tabindex), 10)
		//label = fmt.Sprint(label, indexstr)
		*label = append((*label), []rune(indexstr)...)
		break
	case Stlabel:
		copy(*label, []rune("_loop")[:])
		indexstr = strconv.FormatInt(int64(tabindex), 10)
		*label = append((*label), []rune(indexstr)...)
		break
	case Stprogram:
		fallthrough
	case Stvariable:
		fallthrough
	case Stparameter:
		fallthrough
	case Stprocedure:
		ivalue = this.attribTable[tabindex].thisname
		*label = append(this.stringtable[:], []rune{rune(this.nametable[ivalue].strstart), rune(0)}...)

		if len(*label) >= 5 {
			indexstr = strconv.FormatInt(int64(tabindex), 10)
			*label = append(*label, []rune(indexstr)...)
		}
	}
	copy(this.attribTable[tabindex].label[:], *label)
}

func (this *SymbolTable) labelscope(procindex int) int {

	var oldsymptr, symptr, numbytes int //totalbytes int
	var label = make([]rune, labelSize)

	for symptr = this.Getivalue(procindex); symptr != 0; symptr = this.Getivalue(symptr) {
		//numrunes += (this.getdatatype(symptr) == dtinteger)? 2 : 4;
		numbytes += 2
		if this.Getdatatype(symptr) != Dtinteger {
			numbytes += 2
		}

	}

	//totalbytes = numbytes
	for oldsymptr, symptr = symptr, this.Getivalue(procindex); symptr != 0; oldsymptr, symptr = symptr, this.Getivalue(symptr) {
		this.paramlabel(symptr, &label, &numbytes)
	}

	numbytes -= 2
	if this.Getivalue(procindex) == 0 {
		symptr = procindex + 1
	} else {
		symptr = this.attribTable[oldsymptr].scopenext
	}
	for ; symptr != -1 && this.Getsmclass(symptr) != Stprocedure; symptr = this.attribTable[symptr].scopenext {
		this.paramlabel(symptr, &label, &numbytes)
	}
	return (-numbytes - 2)
}

func (this *SymbolTable) paramlabel(tabindex int, label *[]rune, bytecount *int) {
	var indexstr string
	//	enum symboltype thissymbol;

	if *bytecount < 0 {
		if this.Getdatatype(tabindex) == Dtinteger {
			copy(*label, []rune("[bp")[:])

		} else {
			copy(*label, []rune("[bp")[:])
		}
	} else {
		copy(*label, []rune("[bp"))
	}
	if *bytecount > 0 {
		*label = append((*label), []rune("+")...)
	}

	indexstr = strconv.FormatInt(int64(*bytecount), 10)
	*label = append((*label), []rune(indexstr)...)
	//*runecount -= this.Ggetdatatype(tabindex) == dtinteger? 2: 4;
	*bytecount -= 2
	if this.Getdatatype(tabindex) != Dtinteger {
		*bytecount -= 2 //- 4
	}

	*label = append((*label), []rune("]")...)
	copy(this.attribTable[tabindex].label[:], *label)
}

// PrintLexeme() - Print the lexeme for a given token
func (this *SymbolTable) Printlexeme(tabindex int) {
	i := this.attribTable[tabindex].thisname
	j := this.nametable[i].strstart
	s := append(this.stringtable[:], rune(j))
	fmt.Print(string(s))
}

// PrintToken() -	Print the token class's name given the token
//                  class.
func (this *SymbolTable) Printtoken(i int) {
	fmt.Println(tokenTypes[this.gettok_class(i)])
}

// GetTok_Class()	- Returns the token class for the symbol
func (this *SymbolTable) gettok_class(tabindex int) TokenType {
	return this.attribTable[tabindex].tok_class
}

// InitProcEntry() -	Initialize the entry on the
//						Procedure Stack with the current symbol
func (this *SymbolTable) initprocentry(tabindex int) procstackitem {
	var thisproc procstackitem

	thisproc.proc = tabindex
	thisproc.sstart = -1
	thisproc.snext = -1
	return thisproc
}

// LexemeInCaps() -	Print the lexeme in capital letters
//					This makes it more distinctive
func (this *SymbolTable) LexemeInCaps(tabindex int) {
	//	Get the index within the string table
	//	where the lexeme starts
	i := this.attribTable[tabindex].thisname
	//	Until you encounter the ending null rune,
	//	Print the character in upper case.
	for j := this.nametable[i].strstart; this.stringtable[j] != 0; j++ {
		fmt.Print(strings.ToUpper(string(this.stringtable[j])))
	}

	fmt.Println()
}

func (this *SymbolTable) Getrvalue(tabindex int) float32 {
	return this.attribTable[tabindex].value.val.(float32)
}

func (this *SymbolTable) Getivalue(tabindex int) int {
	return this.attribTable[tabindex].value.val.(int)
}

func (this *SymbolTable) Getdatatype(tabindex int) DataType {
	return this.attribTable[tabindex].dataclass
}

func (this *SymbolTable) Getsmclass(tabindex int) SemanticType {
	return this.attribTable[tabindex].smtype
}

//	True if a valid type (real or integer); False if not
func (this *SymbolTable) Isvalidtype(tabindex int) bool {
	dataclass := this.attribTable[tabindex].dataclass
	return dataclass == Dtinteger || dataclass == Dtreal
}

// Returns the size of the attribute table
func (this *SymbolTable) Tablesize(tabindex int) int {
	return this.attribTabLen
}
