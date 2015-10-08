package pascomp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unicode"
)

//GO's scope visability is determined by the letter case Upper for public and lower for packaged private
const endOfFile rune = 0 //careful with this look into another time

//////////////////////////Scanner Implementation//////////////////////////
//Go does not support classes per say but what they do support is a struct with specialized functions that act as methods when initialized
type Scanner struct {
	St         *SymbolTable
	originFile *os.File
	reader     *bufio.Reader
	lineNum    int
	lookahead  rune
}

////////////////////////////////////////////////////////////////////
//Mark: Public Scanner Functions
////////////////////////////////////////////////////////////////////

func (this *Scanner) NewScanner(arg ...string) {
	this.St = NewSymbolTable()

	var filename string

	if argc := len(arg); argc == 2 {
		filename = arg[1]
	} else {
		fmt.Print("Enter the filename of the Pascal Program?\t")
		fmt.Scanf("%s", &filename)
	}

	if filepath.Ext(filename) != ".pas" {
		//log.Fatal("Usage: ", filepath.Base(arg[0]), " <FileName>.pas")
		log.Fatal("Only pascal files with extension `.pas` allowed")
	}

	file, err := os.Open(filename) // For read access.
	if err == io.ErrUnexpectedEOF {
		log.Fatal("Cannout Open ", filename, "\n", err)
		//Fatal calls os.exit()
	}
	this.originFile = file
	this.reader = bufio.NewReader(this.originFile)

	this.lineNum = 1
	if err != io.EOF {
		this.lookahead = this.firstChar()
	}
}

//Fake Destructor please remember to call defer in main method
func (this *Scanner) DeinitScanner() {

	if err := this.originFile.Close(); err != nil {
		panic(err)
		//Calls panic which is akin to a throws poping up the call
		//stack until something handles it or the program fails
	}
}

//GO's scope visability is determined by the letter case Upper for public and lower for packaged private
func (this *Scanner) GetToken(tabIndex *int) (token TokenType, lexeme string) {

	var char rune = this.lookahead
	if char == endOfFile {
		//If end of file return it
		return Tokeof, lexeme
	}

	lexeme = string(char)
	this.lookahead = this.getc()

	switch {
	case unicode.IsLetter(char):
		this.scanWord(&token, &lexeme, tabIndex)
	case unicode.IsNumber(char):
		this.scanNum(&token, &lexeme, tabIndex)
	default:
		this.scanOp(&token, &lexeme, tabIndex)
	}

	//Do not need to return since using return values directly
	return
}

////////////////////////////////////////////////////////////////////
//Mark: Private Scanner Functions
////////////////////////////////////////////////////////////////////

func (this *Scanner) scanWord(token *TokenType, lexeme *string, tabIndex *int) {

	var char = this.lookahead

	for char != endOfFile && (unicode.IsLetter(char) || unicode.IsNumber(char)) {
		*lexeme += string(char)
		char = this.getc()
	}
	//Put back last invalid character
	this.ungetc(char)
	//Set lookahead for next lexeme
	this.lookahead = this.firstChar()

	//Finally check the symbol table for correctness
	/*
	 *	If the lexeme is already in the symbol table,
	 *	return its tokenclass.  If it isn't,            it must
	 *	be an identifier whose type we do not know yet.
	 */
	if this.St.Installname(*lexeme, tabIndex) {
		*token = this.St.gettok_class(*tabIndex)
	} else {
		this.St.Setattrib(*tabIndex, Stunknown, Tokidentifier)
		*token = Tokidentifier
	}

}

func (this *Scanner) scanNum(token *TokenType, lexeme *string, tabIndex *int) {

	var isFloat bool
	var char = this.lookahead

	for char != endOfFile && unicode.IsNumber(char) {
		*lexeme += string(char)
		char = this.getc()
	}

	if char == '.' || char == 'E' {
		//this is a floating point number
		isFloat = true
		//append to string
		*lexeme += string(char) //Append 'E' or '.' to lexeme
		//Check for special exponential condition
		if char == 'E' {
			char = this.getc() //get next character to see sign
			if char == '-' {
				//Consume to prevent error
				*lexeme += string(char)
			} else if char != '+' {
				//It will be assumed + and nothing before a number is the same
				this.ungetc(char)
			}
		}
		char = this.getc()
		for char != endOfFile && unicode.IsNumber(char) {
			*lexeme += string(char)
			char = this.getc()
		}
	}

	//Put back last invalid character
	this.ungetc(char)
	//Set lookahead for next lexeme
	this.lookahead = this.firstChar()

	//Finally check the symbol table for correctness

	//If there is no fractional part, it is an integer literal
	// constant.  Otherwise, it is a real literal constant.
	// Firstly, is it already in the symbol table?
	if this.St.Installname(*lexeme, tabIndex) {
		*token = this.St.gettok_class(*tabIndex)
		return
	} else {

		this.St.Setattrib(*tabIndex, Stunknown, Tokconstant)
		this.St.Installdatatype(*tabIndex, Stliteral, Dtreal)

		if isFloat {
			// If is it real?
			rval, _ := strconv.ParseFloat(*lexeme, 32)
			this.St.SetFvalue(*tabIndex, float32(rval))
		} else {
			// Must be an integer literal
			ival, _ := strconv.Atoi(*lexeme)
			this.St.SetIvalue(*tabIndex, ival)
		}
		*token = this.St.gettok_class(*tabIndex)
		return
	}
	this.ungetc(this.lookahead)
	*token = this.St.gettok_class(*tabIndex)
	return
}

func (this *Scanner) scanOp(token *TokenType, lexeme *string, tabIndex *int) {
	//Put back last invalid character
	this.ungetc(this.lookahead)
	//Set lookahead for next lexeme
	this.lookahead = this.firstChar()
	//Finally check the symbol table for correctness

	//	The only valid operators are already in the symbol
	//	table.  If it's not there already, it must be
	//	invalid

	if !this.St.IsPresent(*lexeme, tabIndex) {
		log.Fatal(lexeme, " is an illegal operator on line #", this.lineNum)
		//exit(3);
	}

	*token = this.St.gettok_class(*tabIndex)
	return
}

func (this *Scanner) getc() rune {
	var (
		char rune
		size int
		err  error
	)

	char, size, err = this.reader.ReadRune()
	//catch unintended errors
	if err != nil && err != io.EOF {
		panic(err)
	}
	//Return adequate identifier if at end of file
	if err == io.EOF {
		char = endOfFile
	} else if char == '\n' {
		this.lineNum++
	} else if size == 0 {
		char = ' ' //uneeded
	}

	return unicode.ToUpper(char) //Force all alphabet to upper
}

func (this *Scanner) ungetc(char rune) {
	if char == '\n' {
		this.lineNum--
	}
	this.reader.UnreadRune()
}

func (this *Scanner) peek() rune {
	char := this.getc()
	this.ungetc(char)
	return char
}

func (this *Scanner) firstChar() (char rune) {
	for {
		//check for spaces
		if char = this.getc(); char == endOfFile {
			return endOfFile
		} else if char == '{' {
			//Ignore Comments
			//Documentation specifies no nested comments
			for char != endOfFile && char != '}' {
				char = this.getc()
			}
		} else if !unicode.IsSpace(char) {
			//Handle spaces
			//we finally found a viable character
			return char
		}
	}
}
