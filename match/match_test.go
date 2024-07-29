package match

import "testing"

func TestMatchFilter(t *testing.T) {
	type args struct {
		file   string
		filter string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test file test/*",
			args: args{
				file:   "/test/file.txt",
				filter: "test/*",
			},
			want: true,
		},
		{
			name: "test file test",
			args: args{
				file:   "/test/",
				filter: "test/",
			},
			want: true,
		},
		{
			name: "test file api/data*",
			args: args{
				file:   "/api/data.json",
				filter: "api/data*",
			},
			want: true,
		},
		{
			name: "test file api/data*.json",
			args: args{
				file:   "/api/data.json",
				filter: "api/data*.json",
			},
			want: true,
		},
		{
			name: "test file api/data/*.txt",
			args: args{
				file:   "/api/data.txt",
				filter: "api/data/*.txt",
			},
			want: false,
		},
		{
			name: "test file something/*.txt",
			args: args{
				file:   "/something/file.txt",
				filter: "something/*.txt",
			},
			want: true,
		},
		{
			name: "test file something/*",
			args: args{
				file:   "/something/file.dat",
				filter: "something/*",
			},
			want: true,
		},
		{
			name: "*_test.go",
			args: args{
				file:   "/something/file_test.go",
				filter: "*_test.go",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchFilter(tt.args.file, tt.args.filter); got != tt.want {
				t.Errorf("MatchFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
