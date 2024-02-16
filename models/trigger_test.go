package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	testCase struct {
		name   string
		input  string
		expect string
		Func   func(t *testing.T, root any)
	}
)

func TestTransCompoundTrigger(t *testing.T) {
	tests := []testCase{
		{
			name: "simple coumpound trigger",
			input: `create or replace trigger comp_test
for insert or update or delete
on emp_test
compound trigger
      before statement is
      begin
      DBMS_OUTPUT.PUT_LINE('1');
      DBMS_OUTPUT.PUT_LINE('2');
      end  before statement;

      before each row is
      begin
      DBMS_OUTPUT.PUT_LINE('2');
      end before each row;

      after each row is
      begin
      DBMS_OUTPUT.PUT_LINE('3');
      end after each row;

      after statement is
      begin
      DBMS_OUTPUT.PUT_LINE('4');
      end after statement;
end;`,
			expect: `CREATE OR REPLACE TRIGGER comp_test
BEFORE 
insert OR update OR delete ON emp_test

BEGIN
DBMS_OUTPUT.PUT_LINE('1');
      DBMS_OUTPUT.PUT_LINE('2');
END;`,
			Func: func(t *testing.T, root any) {
				testCase, ok := root.(testCase)
				require.True(t, ok)
				require.NotNil(t, testCase)

				plsql := NewPlsqlTrans(testCase.input)
				sql, err := plsql.TransSql()
				require.Nil(t, err)
				require.Equal(t, testCase.expect, sql)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.Func(t, test)
		})
	}
}
