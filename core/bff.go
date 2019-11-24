package core

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Bff struct {
	Routine	*Routine
}

func NewBffExecution(routine string) (*Bff, error) {

	file, err := os.Open(routine)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer file.Close()

	newBFF := Bff{
		Routine: ParseNewRoutine(file),
	}
	return &newBFF, nil
}

func (bff *Bff) GetVariable(search string) string {
	search = strings.Replace(search, "{", "", -1)
	search = strings.Replace(search, "}", "", -1)
	for name, variable := range bff.Routine.Parameters {
		if strings.ToLower(name) == strings.ToLower(search) {
			return variable
		}
	}
	return ""
}

func (bff *Bff) Execute(parameters string) {
	err := bff.Routine.Prepare()
	if err != nil {
		log.Println(err.Error())
		return
	}

	if parameters != "" {
		if strings.Contains(parameters, "=") {
			if strings.Contains(parameters, ";") {
				params := strings.Split(parameters, ";")
				for _, param := range params {
					if strings.Contains(param, "=") {
						info := strings.Split(param, "=")
						bff.Routine.Parameters[info[0]] = info[1]
					}
				}
			} else {
				if strings.Contains(parameters, "=") {
					info := strings.Split(parameters, "=")
					bff.Routine.Parameters[info[0]] = info[1]
				}
			}
		} else {
			log.Fatal("Can't parse parameters please use format type: param1=value1;param2=value2")
			return
		}
	}

	bff.Routine.ExecutedAt = time.Now()
	for _, step := range bff.Routine.Steps {
		if strings.Contains(step.Line, "{") && strings.Contains(step.Line, "}") {
			r, _ := regexp.Compile("\\{(.*?)\\}")
			step.Line = strings.TrimSpace(r.ReplaceAllString(step.Line, bff.GetVariable(r.FindString(step.Line))))
		}

		if bff.Routine.Output == "screen" {
			var outb bytes.Buffer
			cmd := exec.Command("bash", "-c", step.Line)
			cmd.Stdout = &outb
			cmd.Stderr = os.Stderr
			fmt.Printf("%s Execution of routine '%s'\n", color.HiGreenString("====="), step.Name)
			// fmt.Printf("%s Execution of line '%s'\n", color.GreenString("====="), step.Line)
			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println(strings.TrimSpace(outb.String()))
			if step.Store != "" {
				bff.Routine.Parameters[step.Store] = strings.TrimSpace(outb.String())
			}
		}
	}
	bff.Routine.FinishedAt = time.Now()
}
