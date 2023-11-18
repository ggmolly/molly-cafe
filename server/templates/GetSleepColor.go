package templates

import "html/template"

func GetSleepColor(secondsSlept int) template.HTML {
	if secondsSlept > SleepGoal {
		return template.HTML("green")
	} else if secondsSlept > SleepGoal*0.925 {
		return template.HTML("yellow")
	}
	return template.HTML("red")
}
