package main

import "fmt"
import "bufio"
import "os"
import "strconv"
import "strings"

// Token types
const (
	UNKNOWN = iota
	INTEGER = iota
	PLUS    = iota
	MINUS   = iota
	EOF     = iota
)

type errSyntax string

func (e errSyntax) Error() string {
	return "ErrSyntax: " + string(e)
}

// Token keeps Type and Value of user's input
type Token struct {
	Type  int
	Value string
}

// Interpreter handles user's input
type Interpreter struct {
	Text         string
	Pos          int
	CurrentToken Token
}

// Expr parses and interprets user's input
func (i *Interpreter) Expr() (int, error) {
	var err error
	i.CurrentToken, _ = i.getNextToken()
	if err != nil {
		return 0, err
	}

	left := i.CurrentToken
	err = i.eat(INTEGER)
	if err != nil {
		return 0, err
	}

	op := i.CurrentToken
	switch opType := op.Type; opType {
	case PLUS:
		err = i.eat(PLUS)
	case MINUS:
		err = i.eat(MINUS)
	default:
		err = i.eat(UNKNOWN)
	}

	if err != nil {
		return 0, err
	}

	right := i.CurrentToken
	err = i.eat(INTEGER)
	if err != nil {
		return 0, err
	}

	leftValue, _ := strconv.Atoi(left.Value)
	rightValue, _ := strconv.Atoi(right.Value)

	result := 0
	switch opType := op.Type; opType {
	case PLUS:
		result = leftValue + rightValue
	case MINUS:
		result = leftValue - rightValue
	}

	return result, err
}

func (i *Interpreter) getNextToken() (token Token, err error) {
	if i.Pos > len(i.Text)-1 {
		return Token{Type: EOF}, nil
	}

	currentChar := string((i.Text[i.Pos]))

	if _, err := strconv.Atoi(currentChar); err == nil {
		i.Pos++
		token = Token{INTEGER, currentChar}
	} else if currentChar == "+" {
		i.Pos++
		token = Token{PLUS, currentChar}
	} else if currentChar == "-" {
		i.Pos++
		token = Token{MINUS, currentChar}
	} else {
		err = errSyntax("Unkown token " + currentChar)
	}
	return
}

func (i *Interpreter) eat(TokenType int) (err error) {
	if i.CurrentToken.Type == TokenType {
		i.CurrentToken, err = i.getNextToken()
	} else {
		err = errSyntax("Expecting different token")
	}
	return
}

func main() {
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, err := reader.ReadString('\n')

		if err != nil {
			break
		}
		if strings.TrimSpace(text) == "" {
			continue
		}

		interpreter := &Interpreter{Text: text}
		result, err := interpreter.Expr()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}
}
