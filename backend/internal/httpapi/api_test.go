package httpapi

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "space separated words keep every character",
			input: "cats dogs",
			want:  []string{"cats", "dogs"},
		},
		{
			name:  "commas and mixed spacing",
			input: "sleep,  #gainz,ai",
			want:  []string{"sleep", "gainz", "ai"},
		},
		{
			name:  "hash prefixes and case are normalised",
			input: "#AI #Gainz",
			want:  []string{"ai", "gainz"},
		},
		{
			name:  "duplicates collapse",
			input: "ai ai ai",
			want:  []string{"ai"},
		},
		{
			name:  "a word carrying punctuation is dropped whole; hyphens, underscores and digits survive",
			input: "ok! yes? more-ai gainz_plan 5x",
			want:  []string{"more-ai", "gainz_plan", "5x"},
		},
		{
			name:  "nothing valid yields no tags",
			input: "!!! ??? ###",
			want:  []string{},
		},
		{
			name:  "the cap is sixteen tags",
			input: strings.Join(tagWords(40), " "),
			want:  tagWords(16),
		},
		{
			name:  "words longer than forty characters are dropped",
			input: strings.Repeat("a", 40) + " " + strings.Repeat("b", 41),
			want:  []string{strings.Repeat("a", 40)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseTags(test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("parseTags(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func tagWords(count int) []string {
	words := make([]string, 0, count)
	for i := 0; i < count; i++ {
		words = append(words, "tag"+string(rune('a'+i%26))+string(rune('0'+i/26)))
	}
	return words
}
