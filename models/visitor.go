package models

import (
	"github.com/hashicorp/go-multierror"

	"procinspect/pkg/semantic"
)

type (
	transVisitor struct {
		semantic.StubNodeVisitor

		source string

		Statements []Statement
	}
)

func newTransVisitor(source string) *transVisitor {
	v := &transVisitor{}
	v.NodeVisitor = v
	v.source = source
	return v
}

func (v *transVisitor) VisitChildren(node semantic.AstNode) error {
	var result *multierror.Error
	for _, child := range semantic.GetChildren(node) {
		e := child.Accept(v)
		if e != nil {
			result = multierror.Append(result, e)
		}
	}

	return result.ErrorOrNil()
}

func (v *transVisitor) VisitScript(node *semantic.Script) error {
	for _, stmt := range node.Statements {
		switch s := stmt.(type) {
		case *semantic.CreateCompoundDmlTriggerStatement:
			err := s.Accept(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *transVisitor) VisitCreateCompoundDmlTriggerStatement(node *semantic.CreateCompoundDmlTriggerStatement) error {
	trigger := CompoundTrigger{
		Replace:     true,
		Name:        node.Name,
		Events:      node.Events,
		OnTableView: node.TableView,
	}

	body := node.TriggerBody.(*semantic.CompoundTriggerBlock)
	for _, tp := range body.TimingPoints {
		trigger.TimingPoints = append(trigger.TimingPoints, &TimingPoint{
			Before:     tp.IsBefore,
			ForEachRow: tp.ForEachRow,
			Body:       v.source[tp.Body.Span().Start : tp.Body.Span().End+1],
		})
	}
	v.Statements = append(v.Statements, &trigger)
	return nil
}
