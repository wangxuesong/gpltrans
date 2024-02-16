package models

import (
	"bytes"
	"strings"
	"text/template"
)

type (
	CompoundTrigger struct {
		Replace      bool
		Name         string
		Events       []string
		OnTableView  string
		TimingPoints []*TimingPoint
	}

	TimingPoint struct {
		Before     bool
		ForEachRow bool
		Body       string
	}

	TriggerTrans struct {
	}
)

func (c *CompoundTrigger) statement() {}

func NewTriggerTrans() *TriggerTrans {
	return &TriggerTrans{}
}

func (t *TriggerTrans) TransSql(trigger Statement) (string, error) {
	const createTmpl = `CREATE {{if .Replace}}OR REPLACE {{end}}TRIGGER {{.Name}}`
	const beforeTmpl = `{{if .IsBefore}}BEFORE {{else}}AFTER {{end}}`
	const forTmpl = `{{template "events" .Events}} ON {{.OnTableView}}`
	const forEachTmpl = `{{if .ForEachRow}}FOR EACH ROW {{end}}`
	const bodyTmpl = "BEGIN\n{{.Body}}\nEND;"
	const eventsTmpl = `{{define "events"}}{{range $index, $value := .}}{{if $index}} OR {{end}}{{$value}}{{end}}{{end}}`

	tmpl, err := template.New("create").Parse(
		strings.Join([]string{createTmpl,
			beforeTmpl,
			forTmpl,
			forEachTmpl,
			bodyTmpl,
		},
			"\n"))
	if err != nil {
		return "", err
	}
	_, err = tmpl.Parse(eventsTmpl)
	if err != nil {
		return "", err
	}

	trig := trigger.(*CompoundTrigger)
	tp := trig.TimingPoints[0]
	xxx := struct {
		Replace     bool
		Name        string
		Events      []string
		OnTableView string
		IsBefore    bool
		ForEachRow  bool
		Body        string
	}{
		Replace:     trig.Replace,
		Name:        trig.Name,
		Events:      trig.Events,
		OnTableView: trig.OnTableView,
		IsBefore:    tp.Before,
		ForEachRow:  tp.ForEachRow,
		Body:        tp.Body,
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, xxx)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
