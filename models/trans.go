package models

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

func (t *PlsqlTrans) TransSql() string {
	return ""
}
