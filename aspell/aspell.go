package aspell

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"slices"
	"sort"
	"strings"

	"check-commit/match"

	"github.com/fatih/camelcase"
)

type Aspell struct {
	Mode         mode     `yaml:"mode"`
	MinLength    int      `yaml:"min_length"`
	Ignore       []string `yaml:"ignore"`
	AllowedWords []string `yaml:"allowed"`
	HelpText     string   `yaml:"-"`
}

var (
	camelCaseOK = map[string]struct{}{
		"HAProxy": {},
		"golang":  {},
		"ascii":   {},
		"api":     {},
	}
	camelCaseNotOK = map[string]struct{}{}
)

func (a Aspell) checkSingle(data string, allowedWords []string) error {
	var words []string
	var badWords []string

	checkRes, err := checkWithAspellExec(data)
	if checkRes != "" {
		words = strings.Split(checkRes, "\n")
	}
	if err != nil {
		return err
	}

	for _, word := range words {
		wordLower := strings.ToLower(word)
		if len(word) < a.MinLength {
			continue
		}
		if _, ok := camelCaseNotOK[wordLower]; ok {
			badWords = append(badWords, wordLower)
			continue
		}
		if _, ok := camelCaseOK[wordLower]; ok {
			continue
		}
		if slices.Contains(a.AllowedWords, wordLower) || slices.Contains(allowedWords, wordLower) {
			continue
		}
		splitted := camelcase.Split(word)
		if len(splitted) < 2 {
			splitted = strings.FieldsFunc(word, func(r rune) bool {
				return r == '_' || r == '-'
			})
		}
		if len(splitted) > 1 {
			for _, s := range splitted {
				er := a.checkSingle(s, allowedWords)
				if er != nil {
					camelCaseNotOK[wordLower] = struct{}{}
					badWords = append(badWords, word+":"+s)
					break
				}
			}
		} else {
			camelCaseNotOK[wordLower] = struct{}{}
			badWords = append(badWords, word)
		}
	}

	if len(badWords) > 0 {
		m := map[string]struct{}{}
		for _, w := range badWords {
			m[w] = struct{}{}
		}
		badWords = []string{}
		for k := range m {
			badWords = append(badWords, k)
		}
		sort.Strings(badWords)
		return fmt.Errorf("aspell: %s", badWords)
	}
	return nil
}

func (a Aspell) Check(subjects []string, commitsFull []string, content []map[string]string) error {
	var response string
	var checks []string
	switch a.Mode {
	case modeDisabled:
		return nil
	case modeSubject:
		checks = subjects
	case modeCommit:
		checks = commitsFull
	case modeAll:
		for _, file := range content {
			for name, v := range file {
				nextFile := false
				for _, filter := range a.Ignore {
					if match.MatchFilter(name, filter) {
						// log.Println("File", name, "in ignore list")
						nextFile = true
						continue
					}
				}
				if nextFile {
					continue
				}
				var imports []string
				if strings.HasSuffix(name, ".go") {
					imports = match.GetImportWordsFromGoFile(name)
				}
				if err := a.checkSingle(v, imports); err != nil {
					log.Println(name, err.Error())
					response += fmt.Sprintf("%s\n", err)
				}
			}
		}
		checks = commitsFull
	default:
		checks = subjects
	}

	for _, subject := range checks {
		if err := a.checkSingle(subject, []string{}); err != nil {
			log.Println("commit message", err.Error())
			response += fmt.Sprintf("%s\n", err)
		}
	}

	if len(response) > 0 {
		return fmt.Errorf("%s", response)
	}
	return nil
}

func checkWithAspellExec(subject string) (string, error) {
	cmd := exec.Command("aspell", "--list")
	cmd.Stdin = strings.NewReader(subject)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("aspell error: %s, stderr: %s", err, stderr.String())
		return "", err
	}

	return stdout.String() + stderr.String(), nil
}
