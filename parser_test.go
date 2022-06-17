package main

import "testing"

func TestFindCharset(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{name: "empty", text: "", want: defaultCharset},
		{name: "normal", text: "<meta charset=\"cp1251\">", want: "cp1251"},
		{name: "case", text: "<MeTA ChaRseT=\"CP1251\">", want: "CP1251"},
		{
			name: "multiline",
			text: `
<html>
<meta charset="cp1251">
<meta charset="cp2020">
</html>
`,
			want: "cp1251",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findCharset(tt.text); got != tt.want {
				t.Errorf("findCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}
