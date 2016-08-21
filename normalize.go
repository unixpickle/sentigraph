package sentigraph

import "strings"

// Normalize applies some basic rewrite rules to ensure
// that text from social media posts don't contain
// specific digital information and misspellings.
func Normalize(text string) string {
	words := strings.Fields(text)
	var newFields []string
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			continue
		}
		word = strings.ToLower(word)
		word = removeRepeatedLetters(word)
		newFields = append(newFields, word)
	}
	return strings.Join(newFields, " ")
}

// removeRepeatedLetters removes occurrences of letters so
// that no letter is repeated more than twice.
// This was suggested in
// http://cs.stanford.edu/people/alecmgo/papers/TwitterDistantSupervision09.pdf.
func removeRepeatedLetters(s string) string {
	var res string
	var last rune
	var count int
	for _, ch := range s {
		if ch == last {
			count++
		} else {
			last = ch
			count = 1
		}
		if count <= 2 {
			res += string(ch)
		}
	}
	return res
}
