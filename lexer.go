package equilex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Token represents a lexical token
type Token int

const (
	// Illegal token
	Illegal Token = iota
	// EOF means "end of file"
	EOF

	// WS is whitespace
	WS

	// NewLine is one or more consecutive "\n" or "\r".
	NewLine
	// Comment is any kind of comment token
	Comment

	// Identifier is variable identifier
	Identifier

	// String is a string type declaration
	String
	// Logical is a logical (boolean) type declaration
	Logical
	// Number is a numerical type declaration
	Number
	// Date ia a date type declaration
	Date

	// Constants

	// StringConstant is a " delimited string constant
	StringConstant
	// StringMultilineConstant is a $ delimited string constant which may span multiple lines
	StringMultilineConstant
	// IntegerConstant is an integer constant
	IntegerConstant
	// DecimalConstant is a decimal constant
	DecimalConstant
	// DateOrTimeConstant is either a date or time constant (including empty date or time)
	DateOrTimeConstant

	// Comma is ','
	Comma
	// Equals is '='
	Equals
	// LeftParen is '('
	LeftParen
	// RightParen is ')'
	RightParen
	// LeftSquare is '['
	LeftSquare
	// RightSquare is ']'
	RightSquare
	// LeftAngle is '<'
	LeftAngle
	// RightAngle is '>'
	RightAngle

	// Keywords

	// Subtable is the 'subtable' keyword
	Subtable
	// FindRecord is the 'findrecord' keyword
	FindRecord
	// FileOpen is the 'fileopen' keyword
	FileOpen
	// FilePrint is the 'fileprint' keyword
	FilePrint
	// FirstRecord is the 'firstrecord' keyword
	FirstRecord
	// Method is the 'method' keyword
	Method
	// Text is the 'text' keyword
	Text
	// Lookup is the 'lookup' keyword
	Lookup
	// Alert is the 'alert' keyword
	Alert
	// SetIndex is the 'setindex' keyword
	SetIndex
	// Execute is the 'execute' keyword
	Execute
	// MethodSwap is the 'methodswap' keyword
	MethodSwap
	// MethodSetup is the 'methodsetup' keyword
	MethodSetup
	// Process is the 'process' keyword
	Process
	// FormSwap is the 'formswap' keyword
	FormSwap
	// Form is the 'form' keyword
	Form
	// OptimiseTable is the 'optimisetable' keyword
	OptimiseTable
	// OptimiseTableIndexes is the 'optimisetableindexes' keyword
	OptimiseTableIndexes
	// OptimiseDatabase is the 'optimisedatabase' keyword
	OptimiseDatabase
	// OptimiseDatabaseIndexes is the 'optimisedatabaseindexes' keyword
	OptimiseDatabaseIndexes
	// OptimiseAllDatabases is the 'optimisealldatabases' keyword
	OptimiseAllDatabases
	// OptimiseAllDatabasesIndexes is the 'optimisealldatabasesindexes' keyword
	OptimiseAllDatabasesIndexes
	// OptimiseDatabaseHelper is the 'optimisedatabasehelper' keyword
	OptimiseDatabaseHelper
	// ConvertAllDatabases is the 'convertalldatabases' keyword.
	ConvertAllDatabases
	// Command is the 'command' keyword
	Command
	// Task is the 'task' keyword
	Task
	// Shell is the 'shell' keyword
	Shell
	// Export is the 'export' keyword
	Export
	// Import is the 'import' keyword
	Import
	// EmptyDatabase is the 'emptydatabase' keyword
	EmptyDatabase
	// Query is the 'query' keyword
	Query
	// ReportPreview is the 'reportpreview' keyword
	ReportPreview
	// Report is the 'report' keyword
	Report
	// System is the 'system' keyword
	System

	// Public is the 'public' keyword
	Public
	// Procedure is the 'procedure' keyword
	Procedure
	// External is the 'procedure' keyword
	External

	// If is the 'if' keyword
	If
	// Else is the 'else' keyword
	Else
	// ElseIf is the 'elseif' keyword
	ElseIf
	// EndIf is the 'endif' keyword
	EndIf
	// While is the 'while' keyword
	While
	// End is the 'end' keyword
	End
	// Repeat is the 'repeat' keyword
	Repeat
	// Until is the 'until' keyword
	Until
	// For is the 'for' keyword
	For
	// Next is the 'next' keyword
	Next
	// Step is the 'step' keyword
	Step
	// Then is the 'then' keyword
	Then

	// Block is the 'block' keyword
	Block
	// Switch is the 'switch' keyword
	Switch
	// Case is the 'case' keyword
	Case

	// Not is the 'not' keyword
	Not
	// And is the 'and' keyword
	And
	// Or is the 'or' keyword
	Or
	// Xor is the 'xor' keyword
	Xor
	// True is the 'true' keyword
	True
	// False is the 'false' keyword
	False

	// Plus is '+'
	Plus
	// Minus is '-'
	Minus
	// Multiply is '*'
	Multiply
	// Divide is '/'
	Divide
	// Power is '^'
	Power

	// Ampersand is '&'
	Ampersand

	// Today is the 'today' keyword
	Today

	// Backslash is '\'
	Backslash

	// Dot is '.'
	Dot

	// Semicolon is ';'
	Semicolon

	// SysError is the 'syserror' keyword
	SysError

	// TODO: lots are missing and will currently be picked up as VARIABLE probably.
	// Doesn't matter until/unless we try to parse
)

// Lexer is a lexical analyser of Equinox source code.
type Lexer struct {
	r *bufio.Reader
}

// NewLexer returns a new lexical analyser, given a reader that provides equinox source in UTF8.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

// Scan returns the next Token and corresponding string literal
func (s *Lexer) Scan() (tok Token, lit string, err error) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if ch == '|' {
		s.unread()
		return s.scanComment()
	} else if ch == '"' {
		s.unread()
		return s.scanDoubleQuotedLiteral()
	} else if ch == '$' {
		s.unread()
		return s.scanDollarQuotedLiteral()
	} else if ch == '\'' {
		s.unread()
		return s.scanSingleQuotedLiteral()
	} else if ch == '\n' || ch == '\r' {
		s.unread()
		return s.scanNewline()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	} else if isLetter(ch) || ch == '_' {
		s.unread()
		return s.scanIdentifier()
	}

	switch ch {
	case eof:
		return EOF, "", nil
	case ',':
		return Comma, string(ch), nil
	case '=':
		return Equals, string(ch), nil
	case '(':
		return LeftParen, string(ch), nil
	case ')':
		return RightParen, string(ch), nil
	case '[':
		return LeftSquare, string(ch), nil
	case ']':
		return RightSquare, string(ch), nil
	case '<':
		return LeftAngle, string(ch), nil
	case '>':
		return RightAngle, string(ch), nil
	case '+':
		return Plus, string(ch), nil
	case '-':
		return Minus, string(ch), nil
	case '*':
		return Multiply, string(ch), nil
	case '/':
		return Divide, string(ch), nil
	case '^':
		return Power, string(ch), nil

	case '&':
		return Ampersand, string(ch), nil

	case '.':
		return Dot, string(ch), nil

	case ';':
		return Semicolon, string(ch), nil

	case '\\':
		return Backslash, string(ch), nil
	}

	return Illegal, string(ch), nil
}

func (s *Lexer) scanWhitespace() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String(), nil
}

func (s *Lexer) scanSingleQuotedLiteral() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	isDate := false
	isTime := false

	for {
		ch := s.read()
		switch ch {
		case ':':
			isTime = true
		case '-':
			isDate = true
		}

		if isDate && isTime {
			return Illegal, "", fmt.Errorf("malformed date or time '%v' next char is '%v'\n", buf.String(), ch)
		}

		switch ch {
		case '\'':
			buf.WriteRune(ch)
			return DateOrTimeConstant, buf.String(), nil
		case '\n':
			return Illegal, "", fmt.Errorf("unclosed single quote. (TODO: deal with this better)\nbuffer is '%v' and next char is `%v`\n", buf.String(), ch)
		default:
			buf.WriteRune(ch)
		}
	}
}

func (s *Lexer) scanDoubleQuotedLiteral() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		switch ch {
		case '"':
			buf.WriteRune(ch)
			return StringConstant, buf.String(), nil
		case '\n':
			return Illegal, "", fmt.Errorf("unclosed double quote. (TODO: deal with this better)\nbuffer is '%v' and next char is `%v`\n", buf.String(), ch)
		default:
			buf.WriteRune(ch)
		}
	}
}

func (s *Lexer) scanDollarQuotedLiteral() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		switch ch {
		case '$':
			buf.WriteRune(ch)
			return StringMultilineConstant, buf.String(), nil
		case eof:
			return Illegal, "", fmt.Errorf("unclosed double quote. (TODO: deal with this better)\nbuffer is '%v' and next char is `%v`\n", buf.String(), ch)
		default:
			buf.WriteRune(ch)
		}
	}
}

func (s *Lexer) scanComment() (tok Token, lit string, err error) {
	peeked, err := s.r.Peek(2)
	if err != nil {
		return Illegal, "", err
	}

	if bytes.Equal([]byte("|*"), peeked) {
		return s.scanStandardComment()
	}
	return s.scanSingleLineComment()
}

func (s *Lexer) scanSingleLineComment() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// single line comment. scan to end of line (or EOF if first)
	for {
		ch := s.read()
		switch ch {
		case '\n':
			s.unread()
			return Comment, buf.String(), nil
		case eof:
			return Comment, buf.String(), nil
		default:
			buf.WriteRune(ch)
		}

	}
}

func (s *Lexer) scanStandardComment() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	nest := 0
	for {
		if ch := s.read(); ch == eof {
			return Illegal, "", fmt.Errorf("truncated file?:\n\n%v\n\n", string(buf.Bytes()))
		} else if ch == '|' && peek1(s.r) == '*' {
			buf.WriteRune(ch)
			nest++
		} else if ch == '*' && peek1(s.r) == '|' && buf.Len() > 1 {
			buf.WriteRune(ch)
			if nest == 0 {
				s.read()
				buf.WriteRune('|')
				break
			} else {
				nest--
			}
		} else {
			buf.WriteRune(ch)
		}
	}

	return Comment, buf.String(), nil
}

func (s *Lexer) scanNewline() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if ch != '\r' && ch != '\n' {
			if ch != eof {
				s.unread()
			}
			break
		}
		buf.WriteRune(ch)
	}

	return NewLine, buf.String(), nil
}

func (s *Lexer) scanNumber() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	token := IntegerConstant

	for {
		ch := s.read()
		switch {
		case ch == '.':
			buf.WriteRune(ch)
			if token == IntegerConstant {
				token = DecimalConstant
			} else {
				return Illegal, "", fmt.Errorf("malformed number? : '%s' with next char '%v'", buf.String(), string(ch))
			}
		case isDigit(ch):
			buf.WriteRune(ch)
		default:
			s.unread()
			return token, buf.String(), nil
		}
	}
}

func (s *Lexer) scanIdentifier() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Keywords
	switch strings.ToUpper(buf.String()) {
	case "SUBTABLE":
		return Subtable, buf.String(), nil
	case "FINDRECORD":
		return FindRecord, buf.String(), nil
	case "FILEOPEN":
		return FileOpen, buf.String(), nil
	case "FILEPRINT":
		return FilePrint, buf.String(), nil
	case "FIRSTRECORD":
		return FirstRecord, buf.String(), nil
	case "METHOD":
		return Method, buf.String(), nil
	case "TEXT":
		return Text, buf.String(), nil
	case "LOOKUP":
		return Lookup, buf.String(), nil
	case "ALERT":
		return Alert, buf.String(), nil
	case "SETINDEX":
		return SetIndex, buf.String(), nil
	case "EXECUTE":
		return Execute, buf.String(), nil
	case "METHODSWAP":
		return MethodSwap, buf.String(), nil
	case "METHODSETUP":
		return MethodSetup, buf.String(), nil
	case "PROCESS":
		return Process, buf.String(), nil
	case "FORMSWAP":
		return FormSwap, buf.String(), nil
	case "FORM":
		return Form, buf.String(), nil
	case "OPTIMISETABLE":
		return OptimiseTable, buf.String(), nil
	case "OPTIMISETABLEINDEXES":
		return OptimiseTableIndexes, buf.String(), nil
	case "OPTIMISEDATABASE":
		return OptimiseDatabase, buf.String(), nil
	case "OPTIMISEDATABASEINDEXES":
		return OptimiseDatabaseIndexes, buf.String(), nil
	case "OPTIMISEALLDATABASES":
		return OptimiseAllDatabases, buf.String(), nil
	case "OPTIMISEALLDATABASESINDEXES":
		return OptimiseAllDatabasesIndexes, buf.String(), nil
	case "OPTIMISEDATABASEHELPER":
		return OptimiseDatabase, buf.String(), nil
	case "CONVERTALLDATABASES":
		return ConvertAllDatabases, buf.String(), nil
	case "COMMAND":
		return Command, buf.String(), nil
	case "TASK":
		return Task, buf.String(), nil
	case "SHELL":
		return Shell, buf.String(), nil
	case "EXPORT":
		return Export, buf.String(), nil
	case "IMPORT":
		return Import, buf.String(), nil
	case "EMPTYDATABASE":
		return EmptyDatabase, buf.String(), nil
	case "QUERY":
		return Query, buf.String(), nil
	case "REPORTPREVIEW":
		return ReportPreview, buf.String(), nil
	case "REPORT":
		return Report, buf.String(), nil
	case "SYSTEM":
		return System, buf.String(), nil

	case "PUBLIC":
		return Public, buf.String(), nil
	case "PROCEDURE":
		return Procedure, buf.String(), nil
	case "EXTERNAL":
		return External, buf.String(), nil

	case "NOT":
		return Not, buf.String(), nil

	case "IF":
		return If, buf.String(), nil
	case "ELSE":
		return Else, buf.String(), nil
	case "ELSEIF":
		return ElseIf, buf.String(), nil
	case "ENDIF":
		return EndIf, buf.String(), nil
	case "WHILE":
		return While, buf.String(), nil
	case "END":
		return End, buf.String(), nil
	case "REPEAT":
		return Repeat, buf.String(), nil
	case "UNTIL":
		return Until, buf.String(), nil
	case "FOR":
		return For, buf.String(), nil
	case "NEXT":
		return Next, buf.String(), nil
	case "STEP":
		return Step, buf.String(), nil
	case "THEN":
		return Then, buf.String(), nil

	case "BLOCK":
		return Block, buf.String(), nil
	case "SWITCH":
		return Switch, buf.String(), nil
	case "CASE":
		return Case, buf.String(), nil

	case "AND":
		return And, buf.String(), nil
	case "OR":
		return Or, buf.String(), nil
	case "XOR":
		return Xor, buf.String(), nil

	case "STRING":
		return String, buf.String(), nil
	case "LOGICAL":
		return Logical, buf.String(), nil
	case "DATE":
		return Date, buf.String(), nil
	case "NUMBER":
		return Number, buf.String(), nil

	case "TRUE":
		return True, buf.String(), nil
	case "FALSE":
		return False, buf.String(), nil

	case "TODAY":
		return Today, buf.String(), nil

	case "SYSERROR":
		return SysError, buf.String(), nil
	}

	return Identifier, buf.String(), nil
}

func (s *Lexer) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func peek1(r *bufio.Reader) byte {
	bytes, err := r.Peek(1)
	if err != nil {
		panic(err)
	}
	return bytes[0]
}

func (s *Lexer) unread() {
	err := s.r.UnreadRune()
	if err != nil {
		panic(err)
	}
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' }

func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

var eof = rune(0)
