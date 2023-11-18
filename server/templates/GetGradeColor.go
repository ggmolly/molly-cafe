package templates

import "html/template"

func GetGradeColor(grade int8) template.HTML {
	if grade < 100 {
		return template.HTML("grade-fail")
	} else if grade < 125 {
		return template.HTML("grade-ok")
	} else if grade == 125 {
		return template.HTML("grade-max-bonus")
	}
	return template.HTML("grade-wip")
}
