package core

type Step struct {
	Name 		string 		`yaml:"name"`
	Program 	string 		`yaml:"program"`
	Arguments 	[]string 	`yaml:"arguments"`
	Line 		string		`yaml:"line"`
	Store 		string 		`yaml:"store"`
}
