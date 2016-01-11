package main

import "fmt"
import "bufio"
import "os"
import "strconv"
import "strings"

const (
	UNKNOWN = iota
	INTEGER = iota
	PLUS    = iota
	MINUS   = iota
	EOF     = iota
)

type ErrSyntax string

func (e ErrSyntax) Error() string {
	return "ErrSyntax: " + string(e)
}

type Token struct {
	Type  int
	Value string
}

type Interpreter struct {
	Text         string
	Pos          int
	CurrentToken Token
}

func (i *Interpreter) Expr() (int, error) {
	var err error
	i.CurrentToken, _ = i.GetNextToken()
	if err != nil {
		return 0, err
	}

	left := i.CurrentToken
	err = i.Eat(INTEGER)
	if err != nil {
		return 0, err
	}

	op := i.CurrentToken
	switch opType := op.Type; opType {
	case PLUS:
		err = i.Eat(PLUS)
	case MINUS:
		err = i.Eat(MINUS)
	default:
		err = i.Eat(UNKNOWN)
	}

	if err != nil {
		return 0, err
	}

	right := i.CurrentToken
	err = i.Eat(INTEGER)
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

func (i *Interpreter) GetNextToken() (token Token, err error) {
	if i.Pos > len(i.Text)-1 {
		return Token{Type: EOF}, nil
	}

	currentChar := string((i.Text[i.Pos]))

	if _, err := strconv.Atoi(currentChar); err == nil {
		i.Pos += 1
		token = Token{INTEGER, currentChar}
	} else if currentChar == "+" {
		i.Pos += 1
		token = Token{PLUS, currentChar}
	} else if currentChar == "-" {
		i.Pos += 1
		token = Token{MINUS, currentChar}
	} else {
		err = ErrSyntax("Unkown token " + currentChar)
	}
	return
}

func (i *Interpreter) Eat(TokenType int) (err error) {
	if i.CurrentToken.Type == TokenType {
		i.CurrentToken, err = i.GetNextToken()
	} else {
		err = ErrSyntax("Expecting different token")
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
