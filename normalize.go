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
			newFields = append(newFields, "USERNAME")
		} else if strings.HasPrefix(word, "http://") ||
			strings.HasPrefix(word, "https://") {
			newFields = append(newFields, "URL")
		} else {
			word = strings.ToLower(word)
			word = removeRepeatedLetters(word)
			newFields = append(newFields, word)
		}
	}
	return strings.Join(newFields, " ")
}

// SeparatePunctuation adds spaces around clusters of
// punctuation.
func SeparatePunctuation(text string) string {
	var res []string
	for _, f := range strings.Fields(text) {
		res = append(res, separatePunctuationWord(f)...)
	}
	return strings.Join(res, " ")
}

// separatePunctuationWord separates punctuation in a
// single word/field.
func separatePunctuationWord(word string) []string {
	punct := map[rune]bool{'!': true, '.': true, ',': true, '?': true}

	var words []string
	var cur string
	var last rune

	for _, ch := range word {
		if len(cur) == 0 {
			last = ch
			cur = string(ch)
			continue
		}
		p := punct[ch]
		if p != punct[last] {
			words = append(words, cur)
			cur = string(ch)
		} else {
			cur += string(ch)
		}
		last = ch
	}
	words = append(words, cur)
	return words
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
