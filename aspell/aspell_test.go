package aspell

import "testing"

func Test_checkWithAspell(t *testing.T) {
	aspell := Aspell{
		Mode:         modeSubject,
		MinLength:    3,
		AllowedWords: []string{"config"},
	}
	tests := []struct {
		name    string
		subject string
		wantErr bool
	}{
		{"OK 1", "BUG/MEDIUM: config: add default location of path to the configuration file", false},
		{"OK 2", "BUG/MEDIUM: config: add default location of path to the configuration file xtra", false},
		{"error - flie", "BUG/MEDIUM: config: add default location of path to the configuration flie", true},
		{"error - locatoin", "CLEANUP/MEDIUM: config: add default locatoin of path to the configuration file", true},
		{"error - locatoin+flie", "CLEANUP/MEDIUM: config: add default locatoin of path to the configuration flie", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := aspell.checkSingle(tt.subject, []string{"xtra"})
			if tt.wantErr && err == nil {
				t.Errorf("checkWithAspell() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("checkWithAspell() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
