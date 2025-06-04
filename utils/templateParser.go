package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func ReplaceTemplateVariables(content string, templateVariables map[string]string) string{
	for key, value := range templateVariables {
		content = strings.ReplaceAll(content,"{{ "+key+" }}",value)
		content = strings.ReplaceAll(content,"{{"+key+"}}",value)
	}
	return content;
}

func RenderTemplate(w http.ResponseWriter, path string, data any) {
	t, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template parsing error:", err)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
	}
}

func TrimSpaces(s string) string {
    return strings.TrimSpace(s)
}