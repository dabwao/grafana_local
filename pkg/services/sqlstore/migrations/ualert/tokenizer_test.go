package ualert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenString(t *testing.T) {
	t1 := Token{Space: ' '}
	assert.Equal(t, " ", t1.String())
	t2 := Token{Literal: "this is a literal"}
	assert.Equal(t, "this is a literal", t2.String())
	t3 := Token{Variable: "this is a variable"}
	assert.Equal(t, "this is a variable", t3.String())
}

func TestTokenizeLiteral(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		token Token
		pos   int
		err   error
	}{{
		name:  "no characters",
		text:  "",
		token: Token{},
		pos:   0,
	}, {
		name:  "string is valid literal",
		text:  "Instance is down",
		token: Token{Literal: "Instance is down"},
		pos:   16,
	}, {
		name:  "string with numbers is a valid literal",
		text:  "Instance 1 is down",
		token: Token{Literal: "Instance 1 is down"},
		pos:   18,
	}, {
		name:  "all spaces",
		text:  "    ",
		token: Token{Literal: ""},
		pos:   4,
	}, {
		name:  "leading space is removed",
		text:  " Instance 1 is down",
		token: Token{Literal: "Instance 1 is down"},
		pos:   19,
	}, {
		name:  "leading spaces is removed",
		text:  "  Instance 1 is down",
		token: Token{Literal: "Instance 1 is down"},
		pos:   20,
	}, {
		name:  "trailing space is removed",
		text:  "Instance 1 is down ",
		token: Token{Literal: "Instance 1 is down"},
		pos:   18,
	}, {
		name:  "trailing spaces are removed",
		text:  "Instance 1 is down  ",
		token: Token{Literal: "Instance 1 is down"},
		pos:   18,
	}, {
		name:  "string is terminated at $",
		text:  "Instance ${instance} is down",
		token: Token{Literal: "Instance"},
		pos:   8,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, pos, err := tokenizeLiteral([]rune(test.text))
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.pos, pos)
			assert.Equal(t, test.token, token)
		})
	}
}

func TestTokenizeSpace(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		token Token
		pos   int
		err   error
	}{{
		name:  "no spaces",
		text:  "",
		token: Token{Space: 0},
		pos:   0,
	}, {
		name:  "a single space",
		text:  " ",
		token: Token{Space: ' '},
		pos:   1,
	}, {
		name:  "two spaces",
		text:  "  ",
		token: Token{Space: ' '},
		pos:   1,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, pos, err := tokenizeSpace([]rune(test.text))
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.pos, pos)
			assert.Equal(t, test.token, token)
		})
	}
}

func TestTokenizeVariable(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		token Token
		pos   int
		err   error
	}{{
		name:  "variable with no trailing text",
		text:  "${instance}",
		token: Token{Variable: "instance"},
		pos:   11,
	}, {
		name:  "variable with trailing text",
		text:  "${instance} is down",
		token: Token{Variable: "instance"},
		pos:   11,
	}, {
		name:  "varaiable with numbers",
		text:  "${instance1} is down",
		token: Token{Variable: "instance1"},
		pos:   12,
	}, {
		name:  "variable with underscores",
		text:  "${instance_with_underscores} is down",
		token: Token{Variable: "instance_with_underscores"},
		pos:   28,
	}, {
		name:  "two variables without spaces",
		text:  "${variable1}${variable2}",
		token: Token{Variable: "variable1"},
		pos:   12,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, pos, err := tokenizeVariable([]rune(test.text))
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.pos, pos)
			assert.Equal(t, test.token, token)
		})
	}
}

func TestTokenizeTmpl(t *testing.T) {
	tests := []struct {
		name   string
		tmpl   string
		tokens Tokens
		err    error
	}{{
		name:   "simple template can be tokenized",
		tmpl:   "${instance} is down",
		tokens: Tokens{{Variable: "instance"}, {Space: ' '}, {Literal: "is down"}},
	}, {
		name: "complex template can be tokenized",
		tmpl: "More than ${value} ${status_code} in the last 5 minutes",
		tokens: Tokens{
			{Literal: "More than"},
			{Space: ' '},
			{Variable: "value"},
			{Space: ' '},
			{Variable: "status_code"},
			{Space: ' '},
			{Literal: "in the last 5 minutes"},
		},
	}, {
		name: "trailing spaces are removed",
		tmpl: "More than ${value} ${status_code} in the last 5 minutes ",
		tokens: Tokens{
			{Literal: "More than"},
			{Space: ' '},
			{Variable: "value"},
			{Space: ' '},
			{Variable: "status_code"},
			{Space: ' '},
			{Literal: "in the last 5 minutes"},
		},
	}, {
		name:   "variables without spaces can be tokenized",
		tmpl:   "${value}${status_code}",
		tokens: Tokens{{Variable: "value"}, {Variable: "status_code"}},
	}, {
		name:   "variables without spaces then literal can be tokenized",
		tmpl:   "${value}${status_code}in the last 5 minutes",
		tokens: Tokens{{Variable: "value"}, {Variable: "status_code"}, {Literal: "in the last 5 minutes"}},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := tokenizeTmpl(test.tmpl)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.tokens, tokens)
		})
	}
}

func TestTokensStringer(t *testing.T) {
	tokens := Tokens{{Variable: "instance"}, {Space: ' '}, {Literal: "is down"}}
	assert.Equal(t, "{{instance}} is down", tokens.String())
}
