package main

import (
	"strings"
)

type ScopeTokenizer struct {
	parse_index int
	tokens      []Token
	code        string
	scopeType   ScopeType
}

type ScopeType int

const (
	Parens ScopeType = iota
	Brackets
	Braces
	Global
)

type TokenType int

const (
	Identifier TokenType = iota
	Number
	String
	Operator
	Keyword
)

var keywords = []string{
	"int",
	"float",
	"bool",
	"string",
}

type Token struct {
	Type  TokenType
	Value string
}

type OF_TOKEN interface {
	NextToken() Token
}

func (s *ScopeTokenizer) NextToken(allowed_chars string) string {
	for s.parse_index < len(s.code) {
		if strings.ContainsRune(allowed_chars, rune(s.code[s.parse_index])) {
			s.parse_index++
			continue
		}
		break
	}
	return s.code[:s.parse_index]
}

var scope_stack []ScopeType

func (s *ScopeTokenizer) parseToken(allowed_chars string) {
	s.consumeSpaces()
	if s.code[s.parse_index] == '\'' {
		s.tokens = append(s.tokens, s.NextString())
	}
	if s.code[s.parse_index] == '"' {
		s.tokens = append(s.tokens, s.NextDbString())
	}
	if strings.ContainsRune(valid_in_identifier_start, rune(s.code[s.parse_index])) {
		s.tokens = append(s.tokens, s.NextIdentifierOrKeyword())
	}
	if strings.ContainsRune(operators, rune(s.code[s.parse_index])) {
		s.tokens = append(s.tokens, s.NextOperator())
	}
	if strings.ContainsRune(punctuation, rune(s.code[s.parse_index])) {
		s.tokens = append(s.tokens, s.NextPunctuation())
	}
	if strings.ContainsRune(valid_numbers_start, rune(s.code[s.parse_index])) {
		s.tokens = append(s.tokens, s.NextNumber())
	}
	if scopeType, ok := scope_openers[rune(s.code[s.parse_index])]; ok {
		s.tokens = append(s.tokens, s.NextScope(scopeType))
		scope_stack = append(scope_stack, scopeType)
	}
	if scopeType, ok := scope_closers[rune(s.code[s.parse_index])]; ok {
		assert(scopeType == s.scopeType)
		scope_stack = scope_stack[:len(scope_stack)-1]
	}

}

func assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
}

func (s *ScopeTokenizer) consumeSpaces() {
	s.NextToken(" \t\n\r")
}

func (s *ScopeTokenizer) NextTokenNotIn(allowed_chars string) string {
	for s.parse_index < len(s.code) {
		if !strings.ContainsRune(allowed_chars, rune(s.code[s.parse_index])) {
			s.parse_index++
			continue
		}
		break
	}
	return s.code[:s.parse_index]
}

func (s *ScopeTokenizer) NextNumber() Token {
	return Token{
		Type:  Number,
		Value: s.NextToken(valid_numbers),
	}
}

func contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func (s *ScopeTokenizer) NextIdentifierOrKeyword() Token {
	next := s.NextToken(valid_in_identifier)
	if contains(keywords, next) {
		return Token{
			Type:  Keyword,
			Value: next,
		}
	}
	return Token{
		Type:  Identifier,
		Value: next,
	}
}

var operators = "+-*/=<>!&|"

func (s *ScopeTokenizer) NextOperator() Token {
	return Token{
		Type:  Operator,
		Value: s.NextToken(operators),
	}
}

var punctuation = "?|$@,.^"

func (s *ScopeTokenizer) NextPunctuation() Token {
	return Token{
		Type:  Operator,
		Value: s.NextToken(punctuation),
	}
}

func (s *ScopeTokenizer) NextString() Token {
	return Token{
		Type:  String,
		Value: s.NextTokenNotIn("'"),
	}
}

func (s *ScopeTokenizer) NextDbString() Token {
	return Token{
		Type:  String,
		Value: s.NextTokenNotIn("\""),
	}
}

func (s *ScopeTokenizer) tillEndOfComment() Token {
	return Token{
		Type:  String,
		Value: s.NextTokenNotIn("\n"),
	}
}

var numbers = "0123456789"
var valid_numbers = numbers + "."
var valid_numbers_start = numbers
var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var valid_in_identifier = numbers + letters + "_"
var valid_in_identifier_start = letters + "_"
var scope_openers = map[rune]ScopeType{
	'(': Parens,
	'[': Brackets,
	'{': Braces,
}
var scope_closers = map[rune]ScopeType{
	')': Parens,
	']': Brackets,
	'}': Braces,
}

func main() {
	code := `
	@add(int a, int b){
		return a + b;
	}



	class Person{
		int age
		string name	
		@have_birthday(){
			.age+=1;
		}
		@Get_status(){
			return "Person is " + age + " years old";
		}
	}
	
	@main(){
		Person person = Person(1, "John");
		person.have_birthday();
		print(person.Get_status());
	}
	`

	split := strings.Split(code, "\n")

	for _, line := range split {

	}

}
