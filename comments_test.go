package equilex

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStandardComment(t *testing.T) {
	assert := assert.New(t)

	c := `foo bar	|* comment 1 *|   foo bar`

	assert.Equal([]string{"|* comment 1 *|"}, findComments(c))
}

func TestStandardCommentSplit(t *testing.T) {
	assert := assert.New(t)

	c := `foo bar	|* comment with 
line break *|   foo bar`

	assert.Equal([]string{"|* comment with \nline break *|"}, findComments(c))
}

func TestEOLComment(t *testing.T) {
	assert := assert.New(t)

	c := `	foo bar	| comment 2  ` // no EOL before EOF
	assert.Equal([]string{"| comment 2  "}, findComments(c))

}

func TestEOLCommentEOF(t *testing.T) {
	assert := assert.New(t)

	c := `	foo bar	| comment 2  
`
	assert.Equal([]string{"| comment 2  "}, findComments(c))

}

func TestCommentInsideString(t *testing.T) {
	assert := assert.New(t)

	c := `
	foo bar	" | not a comment"  
`
	assert.Equal([]string{}, findComments(c))

}

func TestNoComments(t *testing.T) {
	assert := assert.New(t)

	c := `
	foo bar	  
`
	assert.Equal([]string{}, findComments(c))

}

func TestCommentDoesntEndImmediately(t *testing.T) {
	assert := assert.New(t)

	c := `
	foo bar |*| baz *| banana
`
	assert.Equal([]string{"|*| baz *|"}, findComments(c))

}

func findComments(input string) (comments []string) {
	l := NewLexer(strings.NewReader(input))

	comments = []string{}

	for {
		token, literal, err := l.Scan()
		if err != nil {
			panic(err)
		}
		switch token {
		case EOF:
			return
		case Comment:
			comments = append(comments, literal)
		default:
		}
	}
}
