package pascomp

//Constant declarations
const (
	// 8 characters per tabstop
	tabStop int = 8
	// The size of the name table, hash table,
	// string table and attribute table
	nameTableSize   int = 200
	hashTableSize   int = 100
	stringTableSize int = 1200
	attribTableSize int = 200

	// No more than 120 characters per line + null
	maxLine int = 121

	// The Pascal Subset for this project currently contains 21 keywords
	// and 13 other tokens with entries in the symbol table
	// their is also 1 additional to handle the special case float
	numKeywords int = 21
	numOthers   int = 13
	numTokens   int = numKeywords + numOthers
	labelSize   int = 10
)

// The structure of the attribute table entry, which includes:
//	semantic, token and data types
//	index of the procedure in which this symbol appears
//	index of the name table entry
//	index of the symbol's attribute table entry outside this scope
//	index of the next sentry in this scope so we can close them all
//	value of the constant
//	a label usually indicating address in the object code.
type attribTabType struct {
	smtype                SemanticType
	tok_class             TokenType
	dataclass             DataType
	owningprocedure       int
	thisname              int
	outerscope, scopenext int
	value                 valRec
	label                 [labelSize]rune
}

// The structure for name table entries, i.e, a starting point in
// a long array, a pointer to the entry in the attribute table and
//	the next lexeme with the same hash value.
type nameTabType struct {
	strstart  int
	strlength int
	symtabptr int
	nextname  int
}

//////////////////////////COMPLIMENTING DATA OBJECTS//////////////////////////

// The value can be either integer or real.  We save the tag which
// tells is which it is and store the result in a union.
type tagType int

const (
	tint tagType = iota
	treal
)

//Go has no support for unions and current conceivable
//work arounds seem convoluted and horribly contrived
//type valType struct { //union {
//	ival int
//	rval float
//}

// The structure that stores the tag and the value
type valRec struct {
	tag tagType
	//Since no union support keep generic and type decided on declaration
	val interface{} //valType
}

// The structure for each item that is pushed on the procedure stack
// This includes:
//	index of the procedure in the attribute table
//	index of the first attribute table entry for
//		this scope
//	index of the next attribute table entry for
//		this scope
type procstackitem struct {
	proc          int
	sstart, snext int
}

//////////////////////////TOKEN TYPES//////////////////////////

//Enums are just a list of constants in a parenthesis giving it a faux type
//by using a capital letter for the first letter the type is exported and usable outside of the package
type TokenType int

//Golang has one keyword iota that make the consts into enum.

const (
	Tokbegin   TokenType = iota //1
	Tokcall                     //2
	Tokdeclare                  //3
	Tokdo                       //4
	Tokelse
	Tokend
	Tokendif
	Tokenduntil
	Tokendwhile
	Tokif
	Tokinteger
	Tokparameters
	Tokprocedure
	Tokprogram
	Tokread
	Tokreal
	Tokset
	Tokthen
	Tokuntil
	Tokwhile
	Tokwrite
	Tokstar
	Tokplus
	Tokminus
	Tokslash
	Tokequals
	Toksemicolon
	Tokcomma
	Tokperiod
	Tokgreater
	Tokless
	Toknotequal
	Tokopenparen
	Tokcloseparen
	Tokfloat
	Tokidentifier
	Tokconstant
	Tokerror
	Tokeof
	Tokunknown
)

// [...] is to tell the Go intrepreter/compiler to figure out the array size
var tokenTypes = [...]string{"tokbegin",
	"tokcall", "tokdeclare", "tokdo",
	"tokelse", "tokend", "tokendif", "tokenduntil",
	"tokendwhile", "tokif", "tokinteger",
	"tokparameters", "tokprocedure", "tokprogram",
	"tokread", "tokreal", "tokset", "tokthen",
	"tokuntil", "tokwhile", "tokwrite", "tokstar",
	"tokplus", "tokminus", "tokslash", "tokequals",
	"toksemicolon", "tokcomma", "tokperiod",
	"tokgreater", "tokless", "toknotequal",
	"tokopenparen", "tokcloseparen", "tokfloat",
	"tokidentifier", "tokconstant", "tokerror",
	"tokeof", "tokunknown"}

func (tok TokenType) String() string {
	return tokenTypes[tok]
}

//	The key words and operators - used in initializing the symbol
//	table
var keywords = [...]string{"begin", "call", "declare",
	"do", "else", "end", "endif", "enduntil", "endwhile",
	"if", "integer", "parameters", "procedure", "program",
	"read", "real", "set", "then", "until", "while",
	"write", "*", "+", "-", "/", "=", ";",
	",", ".", ">", "<", "!", "(", ")", "_float"}

//////////////////////////SYMANTIC TYPES//////////////////////////
// The semantic types, i.e, keywords, procedures, variables, constants
type SemanticType int

const (
	Stunknown SemanticType = iota
	Stkeyword
	Stprogram
	Stparameter
	Stvariable
	Sttempvar
	Stconstant
	Stenum
	Ststruct
	Stunion
	Stprocedure
	Stfunction
	Stlabel
	Stliteral
	Stoperator
)

var semanticTypes = [...]string{
	"stunknown", "stkeyword", "stprogram",
	"stparameter", "stvariable", "sttempvar",
	"stconstant", "stenum", "ststruct",
	"stunion", "stprocedure", "stfunction",
	"stlabel", "stliteral", "stoperator"}

func (st SemanticType) String() string {
	return semanticTypes[st]
}

//////////////////////////DATA TYPES//////////////////////////
// The data types, i.e, real and integer
type DataType int

const (
	Dtunknown DataType = iota
	Dtnone
	Dtprogram
	Dtprocedure
	Dtinteger
	Dtreal
)

// The data types, i.e, real and integer
var dataTypes = [...]string{
	"dtunknown", "dtnone", "dtprogram",
	"dtprocedure", "dtinteger", "dtreal"}

func (dat DataType) String() string {
	return dataTypes[dat]
}
