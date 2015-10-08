package pascomp

//////////////////////////TOKEN TYPES//////////////////////////

//Enums are just a list of constants in a parenthesis giving it a faux type
type TokenType int

//Golang has one keyword iota that make the consts into enum.
const (
	Tok_Word   TokenType = iota // 1
	Tok_Number                  // 2
	Tok_Op                      // 3
	Tok_Eof                     // 4
)

// [...] is to tell the Go intrepreter/compiler to figure out the array size
var tokenTypes = [...]string{"Tok_Word", "Tok_Number", "Tok_Op", "Tok_Eof"}

const NumTokens = len(tokenTypes)

func (t TokenType) String() string {
	return tokenTypes[t]
}
