package pascomp

import "fmt"

var tokclstring = [...]string{"begin     ", "call      ",
	"declare   ", "do        ", "else      ", "end       ",
	"endif     ", "enduntil  ", "endwhile  ", "if        ",
	"integer   ", "parameters", "procedure ", "program   ",
	"read      ", "real      ", "set       ", "then      ",
	"until     ", "while     ", "write     ", "star      ",
	"plus      ", "minus     ", "slash     ", "equals    ",
	"semicolon ", "comma     ", "period    ", "greater   ",
	"less      ", "notequal  ", "openparen ", "closeparen",
	"float     ", "identifier", "constant  ", "error     ",
	"eof       ", "unknown   "}

//	The names of the semantic types in a format that can be
//	printed  in a symbol table dump
var symtypestring = [...]string{"unknown  ", "keyword  ", "program  ",
	"parameter", "variable ", "temp. var",
	"constant ", "enum     ", "struct   ",
	"union    ", "procedure", "function ",
	"label    ", "literal  ", "operator "}

//	The names of the data types in a format that can be
//	printed  in a symbol table dump
var datatypestring = [...]string{"unknown", "none   ", "program",
	"proced.", "integer", "real   "}

// DumpSymbolTable() -	Prints out the basic symbol table
//						information, including the name and token
//						class
func DumpSymbolTable(st *SymbolTable) {
	var i, j int
	var printstring string

	//	Print the symbol table's heading
	fmt.Print("SYMBOL TABLE DUMP\n-----------------\n\n")
	fmt.Print("                   Token       Symbol     Data")
	fmt.Print("              Owning\n")
	fmt.Print("Index   Name       Class       Type       Type")
	fmt.Print("          Value   Procedure    Label\n")
	fmt.Print("-----   ----       -----       ------     ----")
	fmt.Print("          -----   ---------    ---------------\n")

	//	Print the data for each entry
	for i = 0; i < st.attribTabLen; i++ {
		//	Pause every tenth line
		//if (i%10 == 9) st.getchar();

		//	Print the entry number and lexeme
		fmt.Printf("%5d\t", i)
		st.Printlexeme(i)

		//
		//	After printing the lexeme, move to column 20.  If
		//	the name is too long to permit, go to the next
		//	line
		q := st.attribTable[i].thisname
		s := st.nametable[q].strstart
		e := s + st.nametable[q].strlength
		printstring = string(st.stringtable[s:e])
		if len(printstring) < 11 {
			for j = 0; j < 11-len(printstring); j++ {
				fmt.Print(" ")
			}
		} else {
			fmt.Print("\n          ")
		}
		// Print the token class, symbol type and data type
		fmt.Print(tokclstring[st.attribTable[i].tok_class], "  ")
		fmt.Print(symtypestring[st.attribTable[i].smtype], "  ")
		fmt.Print(datatypestring[st.attribTable[i].dataclass], "  ")

		//	If the value is real or integer, print the
		//	value in the correct format.
		if st.attribTable[i].value.tag == tint {
			fmt.Printf("%10d", st.attribTable[i].value.val.(int))
		} else {
			fmt.Printf("%1.4E", st.attribTable[i].value.val.(float32))
		}
		//	If there is no procedure that owns the symbol
		//	(which is the case for reserved words, operators,
		//	and literals), print "global."
		if st.attribTable[i].owningprocedure == -1 {
			fmt.Print("   global")
			//	Otherwise print the name of the owning
			//	procedure in capital letters to make it
			//	stand out.
		} else {
			fmt.Print("   ")
			st.LexemeInCaps(st.attribTable[i].owningprocedure)
		}

		//	Print the assembly language label.
		fmt.Print("       ", st.attribTable[i].label)
		fmt.Println()
	}

}

func DumpSymbolTable2(st *SymbolTable) {
	var i int

	for i = 0; i < st.namTabLen; i++ {
		//if (i%10 == 9) getchar();
		fmt.Printf("%d\t%d\t%d\t%d\t%d\n", i,
			st.nametable[i].strstart,
			st.nametable[i].strlength,
			st.nametable[i].symtabptr,
			st.nametable[i].nextname)
	}

	for i = 0; i < st.attribTabLen; i++ {
		//if (i%10 == 9) getchar();
		fmt.Printf("%d  %d  %d  %d  %d  %d  %d  %d  %d\t", i,
			st.attribTable[i].smtype,
			st.attribTable[i].tok_class,
			st.attribTable[i].dataclass,
			st.attribTable[i].owningprocedure,
			st.attribTable[i].thisname,
			st.attribTable[i].outerscope,
			st.attribTable[i].scopenext,
			st.attribTable[i].value.tag)
		if st.attribTable[i].value.tag == treal {
			fmt.Printf("%f\t", st.attribTable[i].value.val.(float32))
		} else {
			fmt.Printf("%d\t", st.attribTable[i].value.val.(int))
		}
		fmt.Printf("%s\n", st.attribTable[i].label)
	}

	for i = 0; i < st.strTabLen; i++ {
		//if (i%60 == 59) getchar();
		fmt.Print(st.stringtable[i])
	}
}
