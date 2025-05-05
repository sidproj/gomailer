package utils

import "strings"

func ReplaceTemplateVariables(content string, templateVariables map[string]string) string{
	for key, value := range templateVariables {
		content = strings.ReplaceAll(content,"{{ "+key+" }}",value)
		content = strings.ReplaceAll(content,"{{"+key+"}}",value)
	}
	return content;
}