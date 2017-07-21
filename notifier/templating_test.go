package notifier

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRenderTemplateWithInvalidFile(t *testing.T) {
	templateData := TemplateData{
		ClusterName:  "some cluster",
		SystemStatus: "some status",
		FailCount:    1,
		WarnCount:    2,
		PassCount:    3,
		Nodes:        make(map[string]Messages),
	}

	renderedTemplate, err := renderTemplate(templateData, "some-file-that-does-not-exist", "")

	if err == nil {
		t.Errorf("Expected error but rendered something: %s", renderedTemplate)
	} else if err.Error() != "open some-file-that-does-not-exist: no such file or directory" {
		t.Errorf("Expected error in opening file but got: %s", err.Error())
	}
}

func TestRenderTemplate(t *testing.T) {
	tmpfile, err := templateFile("{{ .SystemStatus }} - {{ .FailCount }}")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	templateData := TemplateData{
		ClusterName:  "some cluster",
		SystemStatus: "some status",
		FailCount:    1,
		WarnCount:    2,
		PassCount:    3,
		Nodes:        make(map[string]Messages),
	}

	renderedTemplate, err := renderTemplate(templateData, tmpfile.Name(), "")

	if err != nil {
		t.Errorf("Rendering failed: %v", err)
	} else if renderedTemplate != "some status - 1" {
		t.Errorf("Expected 'temporary file' but was '%v'", renderedTemplate)
	}
}

func TestRenderDefaultTemplate(t *testing.T) {
	templateData := TemplateData{
		ClusterName:  "some cluster",
		SystemStatus: "some status",
		FailCount:    1,
		WarnCount:    2,
		PassCount:    3,
		Nodes:        make(map[string]Messages),
	}

	renderedTemplate, err := renderTemplate(templateData, "", "{{ .SystemStatus }} - {{ .FailCount }}")

	if err != nil {
		t.Errorf("Rendering failed: %v", err)
	} else if renderedTemplate != "some status - 1" {
		t.Errorf("Expected 'temporary file' but was '%v'", renderedTemplate)
	}
}

func templateFile(content string) (*os.File, error) {
	tmpfile, err := ioutil.TempFile("", "consulAlertsTest")
	if err != nil {
		return nil, err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}
	return tmpfile, nil
}
