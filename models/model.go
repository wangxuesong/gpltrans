package models

type (
	Statement interface {
		statement()
	}

	Transform interface {
		TransSql(trigger Statement) (string, error)
	}
)
