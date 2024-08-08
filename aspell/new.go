package aspell

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func New(filename string) (Aspell, error) {
	var data []byte
	var err error
	fileExists := true
	if data, err = os.ReadFile(filename); err != nil {
		log.Printf("warning: aspell exceptions file not found (%s)", err)
		fileExists = false
	}

	var aspell Aspell
	err = yaml.Unmarshal(data, &aspell)
	if err != nil {
		return Aspell{}, err
	}

	for i, word := range aspell.AllowedWords {
		aspell.AllowedWords[i] = strings.ToLower(word)
	}

	if aspell.MinLength < 1 {
		aspell.MinLength = 3
	}

	switch aspell.Mode {
	case modeDisabled:
	case modeSubject:
	case modeCommit:
	case modeAll:
	case "":
		aspell.Mode = modeSubject
	default:
		return Aspell{}, fmt.Errorf("invalid mode: %s", aspell.Mode)
	}

	log.Printf("aspell mode set to %s", aspell.Mode)
	if fileExists {
		aspell.HelpText = `aspell can be configured with .aspell.yml file.
Add words to allowed list if its false positive`
	} else {
		aspell.HelpText = `aspell can be configured with .aspell.yml file.
content example:
mode: subject
min_length: 3
ignore:
  - go.mod
  - go.sum
  - '*test.go'
  - 'gen/*'
allowed:
  - aspell
  - config
`
	}

	ignoreFiles := []string{"go.mod", "go.sum"}
	for _, file := range ignoreFiles {
		if _, err := os.Stat(file); err == nil {
			log.Printf("aspell: added %s to ignore list", file)
			aspell.Ignore = append(aspell.Ignore, file)
		}
	}

	return aspell, nil
}
