package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

var (
	tmpl     *template.Template
	tmplOnce sync.Once
)

func getTemplate() *template.Template {
	tmplOnce.Do(func() {
		var err error

		// Create root + register helpers
		tmpl = template.New("").Funcs(template.FuncMap{
			"default": func(def string, val any) any {
				if val == nil || val == "" {
					return def
				}
				return val
			},
			"map": func(values ...any) map[string]any {
				m := make(map[string]any)
				for i := 0; i+1 < len(values); i += 2 {
					if key, ok := values[i].(string); ok {
						if key == "Action" {
							if strVal, ok := values[i+1].(string); ok {
								m[key] = template.JS(strVal)
								continue
							}
						}
						m[key] = values[i+1]
					}
				}
				return m
			},
		})

		// Merge page templates
		pages, _ := filepath.Glob("views/*.html")
		for _, page := range pages {
			_, err = tmpl.ParseFiles(page)
			if err != nil {
				panic(err)
			}
		}

		// Merge components from same folder
		components, _ := filepath.Glob("views/templates/*.html")
		for _, comp := range components {
			_, err = tmpl.ParseFiles(comp)
			if err != nil {
				panic(err)
			}
		}

		fmt.Println("Loaded all templates + components into singleton")
	})
	return tmpl
}

func RenderTemplate(w http.ResponseWriter, name string, data any) {
	t := getTemplate()

	if err := t.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
	}
}
