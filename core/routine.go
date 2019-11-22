package core

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type Routine struct {
	Name 		string 				`yaml:"name"`
	Description string 				`yaml:"description"`
	Steps 		[]Step 				`yaml:"steps"`
	Output 		string 				`yaml:"output"`
	Require 	[]string			`yaml:"require"`
	Parameters  map[string]string 	`yaml:"parameters"`
	ExecutedAt 	time.Time
	FinishedAt 	time.Time
}

func ParseNewRoutine(file *os.File) *Routine {

	var routine Routine
	fileReader, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	_ = yaml.Unmarshal([]byte(fileReader), &routine)

	return &routine
}

func (r *Routine) Prepare() error {
	// checking required tools
	for _, require := range r.Require {
		_, err := exec.LookPath(require)
		if err != nil {
			return errors.New("routine '"+r.Name+"' require external tools '"+require+"'")
		}
	}

	for _, step := range r.Steps {
		if step.Store != "" {
			if _, ok := r.Parameters[step.Store]; !ok {
				return errors.New("Parameter '"+step.Store+"' not found in routine configuration")
			}
		}
	}
	return nil
}
