package filter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type NumericFilter []string

func (filter NumericFilter) Where(column string) qm.QueryMod {

	qSlice := []string{}

	for _, q := range filter {
		qSlice = append(qSlice, q)
	}

	return qm.Where(strings.Join(qSlice, " AND "))
}

func (filter NumericFilter) And(column string) qm.QueryMod {

	qSlice := []string{}

	for _, q := range filter {
		qSlice = append(qSlice, column+q)
	}

	return qm.And(strings.Join(qSlice, " AND "))
}

func (filter NumericFilter) Or(column string) qm.QueryMod {

	qSlice := []string{}

	for _, q := range filter {
		qSlice = append(qSlice, column+q)
	}

	return qm.Or(strings.Join(qSlice, " AND "))
}

func (fltr NumericFilter) Validate() error {

	for _, expression := range fltr {

		//strip spaces
		expression = spaceStrip(expression)

		//verify that string prefixed with valid operator
		if expression[0] != '>' && expression[0] != '<' && expression[0] != '=' {
			return errors.New(fmt.Sprint("numeric filter support the following operators < > = recieved ", expression))
		}

		expression = expression[1:]

		//check that this is indeed numeric value
		if !isNumeric(expression) {
			return errors.New(fmt.Sprint("numeric filter support floating point numbers as value ", expression))
		}

	}
	return nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

//strip strings from spaces
func spaceStrip(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
