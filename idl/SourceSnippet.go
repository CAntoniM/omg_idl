package idl

import "fmt"

type SourceSnippet struct {
	file_name   string
	line        string
	line_number int
	start       int
	end         int
}

func generate_spacer(base string, n int) string {
	var buffer string = ""

	for i := 0; i < n; i++ {
		buffer = buffer + base
	}

	return buffer
}

func (s *SourceSnippet) ToString() string {
	var spacer = generate_spacer(" ", s.start)
	var underline = generate_spacer("~", s.end-s.start)
	return fmt.Sprintf("%s:%d | %s\n%s%s^", s.file_name, s.line_number, s.line, spacer, underline)
}
