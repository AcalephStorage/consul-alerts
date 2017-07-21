package notifier

import (
	"bytes"
	"html/template"
)

type TemplateData struct {
	ClusterName  string
	SystemStatus string
	FailCount    int
	WarnCount    int
	PassCount    int
	Nodes        map[string]Messages
}

func (t TemplateData) IsCritical() bool {
	return t.SystemStatus == SYSTEM_CRITICAL
}

func (t TemplateData) IsWarning() bool {
	return t.SystemStatus == SYSTEM_UNSTABLE
}

func (t TemplateData) IsPassing() bool {
	return t.SystemStatus == SYSTEM_HEALTHY
}

func renderTemplate(t TemplateData, templateFile string, defaultTemplate string) (string, error) {
	var tmpl *template.Template
	var err error
	if templateFile == "" {
		tmpl, err = template.New("base").Parse(defaultTemplate)
	} else {
		tmpl, err = template.ParseFiles(templateFile)
	}

	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, t); err != nil {
		return "", err
	}

	return body.String(), nil
}
