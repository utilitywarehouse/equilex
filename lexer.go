package equilex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
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
	// TODO: lots are missing and will currently be picked up as VARIABLE probably.
	// Doesn't matter until/unless we try to parse

	// If is the 'if' keyword
	If
	// Else is the 'else' keyword
	Else
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
func (s *Lexer) Scan() (tok Token, lit string) {

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
		return EOF, ""
	case ',':
		return Comma, string(ch)
	case '=':
		return Equals, string(ch)
	case '(':
		return LeftParen, string(ch)
	case ')':
		return RightParen, string(ch)
	case '[':
		return LeftSquare, string(ch)
	case ']':
		return RightSquare, string(ch)
	case '<':
		return LeftAngle, string(ch)
	case '>':
		return RightAngle, string(ch)
	case '+':
		return Plus, string(ch)
	case '-':
		return Minus, string(ch)
	case '*':
		return Multiply, string(ch)
	case '/':
		return Divide, string(ch)
	case '^':
		return Power, string(ch)

	case '&':
		return Ampersand, string(ch)

	case '.':
		return Dot, string(ch)

	case ';':
		return Semicolon, string(ch)

	case '\\':
		return Backslash, string(ch)
	}

	return Illegal, string(ch)
}

func (s *Lexer) scanWhitespace() (tok Token, lit string) {
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

	return WS, buf.String()
}

func (s *Lexer) scanSingleQuotedLiteral() (tok Token, lit string) {
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
			log.Panicf("malformed date or time '%v' next char is '%v'\n", buf.String(), ch)
		}

		switch ch {
		case '\'':
			buf.WriteRune(ch)
			return DateOrTimeConstant, buf.String()
		case '\n':
			log.Panicf("unclosed single quote. (TODO: deal with this better)\nbuffer is '%v' and next char is `%v`\n", buf.String(), ch)
		default:
			buf.WriteRune(ch)
		}
	}

}

func (s *Lexer) scanDoubleQuotedLiteral() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		switch ch {
		case '"':
			buf.WriteRune(ch)
			return StringConstant, buf.String()
		case '\n':
			log.Panicf("unclosed double quote. (TODO: deal with this better)\nbuffer is '%v' and next char is `%v`\n", buf.String(), ch)
		default:
			buf.WriteRune(ch)
		}
	}

}

func (s *Lexer) scanDollarQuotedLiteral() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		switch ch {
		case '$':
			buf.WriteRune(ch)
			return StringMultilineConstant, buf.String()
		case eof:
			log.Panicf("unclosed double quote. (TODO: deal with this better)\nbuffer is '%v' and next char is `%v`\n", buf.String(), ch)
		default:
			buf.WriteRune(ch)
		}
	}

}

func (s *Lexer) scanComment() (tok Token, lit string) {
	peeked, err := s.r.Peek(2)
	if err != nil {
		panic(err)
	}

	if bytes.Equal([]byte("|*"), peeked) {
		return s.scanStandardComment()
	}
	return s.scanSingleLineComment()
}

func (s *Lexer) scanSingleLineComment() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// single line comment. scan to end of line (or EOF if first)
	for {
		ch := s.read()
		switch ch {
		case '\n':
			s.unread()
			return Comment, buf.String()
		case eof:
			return Comment, buf.String()
		default:
			buf.WriteRune(ch)
		}

	}
}

func (s *Lexer) scanStandardComment() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	nest := 0
	for {
		if ch := s.read(); ch == eof {
			panic(fmt.Sprintf("truncated file?:\n\n%v\n\n", string(buf.Bytes())))
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

	return Comment, buf.String()
}

func (s *Lexer) scanNewline() (tok Token, lit string) {
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

	return NewLine, buf.String()
}

func (s *Lexer) scanNumber() (tok Token, lit string) {
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
				log.Panicf("malformed number? : '%s' with next char '%v'", buf.String(), ch)
			}
		case isDigit(ch):
			buf.WriteRune(ch)
		default:
			s.unread()
			return token, buf.String()
		}
	}
}

func (s *Lexer) scanIdentifier() (tok Token, lit string) {
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
		return Subtable, buf.String()
	case "FINDRECORD":
		return FindRecord, buf.String()
	case "FILEOPEN":
		return FileOpen, buf.String()
	case "FILEPRINT":
		return FilePrint, buf.String()
	case "FIRSTRECORD":
		return FirstRecord, buf.String()
	case "METHOD":
		return Method, buf.String()
	case "TEXT":
		return Text, buf.String()
	case "LOOKUP":
		return Lookup, buf.String()
	case "ALERT":
		return Alert, buf.String()
	case "SETINDEX":
		return SetIndex, buf.String()

	case "NOT":
		return Not, buf.String()

	case "IF":
		return If, buf.String()
	case "ELSE":
		return Else, buf.String()
	case "ENDIF":
		return EndIf, buf.String()
	case "WHILE":
		return While, buf.String()
	case "END":
		return End, buf.String()
	case "REPEAT":
		return Repeat, buf.String()
	case "UNTIL":
		return Until, buf.String()
	case "FOR":
		return For, buf.String()
	case "NEXT":
		return Next, buf.String()
	case "STEP":
		return Step, buf.String()
	case "THEN":
		return Then, buf.String()

	case "BLOCK":
		return Block, buf.String()
	case "SWITCH":
		return Switch, buf.String()
	case "CASE":
		return Case, buf.String()

	case "AND":
		return And, buf.String()
	case "OR":
		return Or, buf.String()
	case "XOR":
		return Xor, buf.String()

	case "STRING":
		return String, buf.String()
	case "LOGICAL":
		return Logical, buf.String()
	case "DATE":
		return Date, buf.String()
	case "NUMBER":
		return Number, buf.String()

	case "TRUE":
		return True, buf.String()
	case "FALSE":
		return False, buf.String()

	case "TODAY":
		return Today, buf.String()

	case "SYSERROR":
		return SysError, buf.String()
	}

	return Identifier, buf.String()
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
