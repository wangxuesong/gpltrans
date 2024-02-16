package models

import "gpltrans/utils"

type (
	TransRequest struct {
		Id     string `json:"id" swaggertype:"string" description:"trans id"`
		Source string `json:"source" swaggertype:"string" description:"source sql"`
	}

	TransResponse struct {
		Id     string `json:"id" swaggertype:"string" description:"trans id"`
		Source string `json:"source" swaggertype:"string" description:"source sql"`
		Target string `json:"target" swaggertype:"string" description:"target sql"`
	}
)

type PlsqlTrans struct {
	Source string
	Target string
}

func NewPlsqlTrans(source string) *PlsqlTrans {
	return &PlsqlTrans{
		Source: source,
	}
}

func (t *PlsqlTrans) TransSql() (string, error) {
	script, err := utils.ParseSql(t.Source)
	if err != nil {
		return "", err
	}

	visitor := newTransVisitor(t.Source)
	err = script.Accept(visitor)
	if err != nil {
		return "", err
	}

	for _, stmt := range visitor.Statements {
		switch s := stmt.(type) {
		case *CompoundTrigger:
			trans := NewTriggerTrans()
			target, err := trans.TransSql(s)
			if err != nil {
				return "", err
			}
			t.Target = target
		}
	}

	return t.Target, nil
}
