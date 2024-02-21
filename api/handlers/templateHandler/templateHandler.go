package templatehandler

import (
	"html/template"
	"net/http"
)

type TemplateHandler struct {
	templates []string
}

func NewTemplateHandler(templatePath string, baseTemplate string) *TemplateHandler {
	arr := []string{baseTemplate, templatePath}
	return &TemplateHandler{templates: arr}
}

func (th *TemplateHandler) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	tmpl, err := template.ParseFiles(th.templates...)
	if err != nil {
    http.Error(w, "Not able to parse the templates. Erro:" + err.Error(), http.StatusInternalServerError)
		return
	}

  err = tmpl.ExecuteTemplate(w,"base", data)
	if err != nil {
    http.Error(w, "Not able to parse the templates. Erro:" + err.Error(), http.StatusInternalServerError)
		return
	}

}
