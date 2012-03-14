package lang

import (
	"fmt"
	"strings"
)

////////// Tokenization

func tokenize(s string) *TokenList {
	replacedOpen := strings.Replace(s, "(", " ( ", -1)
	replacedAll := strings.Replace(replacedOpen, ")", " ) ", -1)
	tokens := strings.Split(replacedAll, " ")

	filteredTokens := make([]string, 0)
	for _, value := range tokens {
		if value != "" {
			filteredTokens = append(filteredTokens, value)
		}
	}

	return &TokenList{filteredTokens, 0}
}

type TokenList struct {
	tokens []string
	i      int
}

func (self *TokenList) String() string {
	return fmt.Sprintf("%#v", *self)
}

func (self *TokenList) pop() string {
	if self.empty() {
		return ""
	}

	token := self.tokens[self.i]
	self.i++
	return token
}

func (self *TokenList) top() string {
	if self.empty() {
		return ""
	}

	return self.tokens[self.i]
}

func (self *TokenList) empty() bool {
	return self.i == len(self.tokens)
}
