package idl

type LexerError struct {
	cause   string
	snippet SourceSnippet
}

func (e LexerError) Error() string {
	var buffer string = e.cause
	if (SourceSnippet{}) != e.snippet {
		buffer = buffer + "\n" + e.snippet.ToString()
	}
	return buffer
}

func makeLexerError(cause string, snippet SourceSnippet) LexerError {
	return LexerError{cause: cause, snippet: snippet}
}
