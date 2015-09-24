package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func main() {
	//Get Arguments from command line
	//TODO: Use flags isntead of bare arguments
	args := os.Args
	//if len(arg) != 2 {
	//	log.Fatal("Usage: ", filepath.Base(arg[0]), " FileName.pas")
	//}
	//EXT := "pas"
	//fileExtension := string(filepath.Ext(arg[0]))
	//fmt.Println(fileExtension)
	//if strings.EqualFold(fileExtension, EXT) == false {
	//	log.Fatal("Extension must be ", EXT)
	//}
	var scanner Scanner //Since their is nill arguments can do C++ style init
	//would use scanner:= new(Scanner) to allocate memory for fields set all members to respective zero and return a pointer
	defer scanner.deinitScanner() //Kind of a hack to mimic destructors
	scanner.initScanner(args...)  //no constructor so do this way
	for {
		tok, tokString := scanner.scan()
		if tok == tokeof {
			break
		}
		fmt.Printf("%-9s %s\n", tokString, tok)
	}
}

//TODO: This will be moved to its own class and extended in the future
//////////////////////////TOKEN TYPES//////////////////////////

//Enums are just a list of constants in a parenthesis giving it a faux type
type TokenType int

//Golang has one keyword iota that make the consts into enum.
const (
	tokword   TokenType = iota // 1
	toknumber                  // 2
	tokop                      // 3
	tokeof                     // 4
)

// [...] is to tell the Go intrepreter/compiler to figure out the array size
var tokenTypes = [...]string{"tokword", "toknumber", "tokop", "tokeof"}

func (t TokenType) String() string {
	return tokenTypes[t]
}

//////////////////////////Scanner Implementation//////////////////////////
//Go does not support classes per say but what they do support is a struct with specialized functions that act as methods when initialized
type Scanner struct {
	originFile *os.File
	reader     *bufio.Reader
}

func (this *Scanner) initScanner(arg ...string) {

	var filename string

	if argc := len(arg); argc == 1 {
		fmt.Print("Enter the filename of the Pascal Program?\t")
		fmt.Scanf("%s", &filename)
	} else if argc == 2 {
		filename = arg[1]
	} else {
		log.Fatal("Usage: ", filepath.Base(arg[0]), " FileName.pas")
	}

	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	this.originFile = file
	this.reader = bufio.NewReader(this.originFile)
}

//Fake Destructor please remember to call defer in main method
func (this *Scanner) deinitScanner() {

	if err := this.originFile.Close(); err != nil {
		panic(err)
	}
}

func (this *Scanner) scan() (token TokenType, tokenString string) {
	char, size, err := this.reader.ReadRune()

	//skip white space and comments
	for {
		//catch unintended errors
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			break
		}
		//skip it when their is no data or a space
		if size == 0 {
			continue
		}
		//skip it when their is no data or a space
		if char == '{' {
			//Ignore Comments
			//Documentation specifies no nested comments
			for char != '}' {
				char, size, err = this.reader.ReadRune()
			}
		} else if !unicode.IsSpace(char) {
			break
		}

		//get next
		char, size, err = this.reader.ReadRune()
	}
	//If end of file return it
	if err == io.EOF {
		return tokeof, tokenString
	}

	switch {
	case unicode.IsLetter(char):
		this.scanWord(&token, &tokenString, string(char))
	case unicode.IsNumber(char):
		this.scanNum(&token, &tokenString, string(char))
	default:
		token = tokop
		tokenString = string(char)
	}
	//Do not need to return since using return values directly
	return
}

func (this *Scanner) scanWord(token *TokenType, tokenString *string, prefix ...string) {
	*token = tokword
	*tokenString += strings.Join(prefix, "")
	char, size, err := this.reader.ReadRune()

	for {
		*tokenString += string(unicode.ToUpper(char))
		char, size, err = this.reader.ReadRune()

		if err == io.EOF || size == 0 || !unicode.IsLetter(char) {
			this.reader.UnreadRune()
			break
		}

	}
}

func (this *Scanner) scanNum(token *TokenType, tokenString *string, prefix ...string) {
	*token = tokword
	*tokenString += strings.Join(prefix, "")
	char, size, err := this.reader.ReadRune()

	for {
		*tokenString += string(unicode.ToUpper(char))
		char, size, err = this.reader.ReadRune()

		if err == io.EOF || size == 0 || !unicode.IsNumber(char) {
			this.reader.UnreadRune()
			break
		}

	}
}
