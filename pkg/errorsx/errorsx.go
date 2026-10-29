package errorsx

import "fmt"

// Problem describes a RFC 9457 compatible application error.
type Problem struct {
	Type   string
	Title  string
	Status int
	Detail string
	Code   string
}

func (p Problem) Error() string {
	return fmt.Sprintf("%s: %s", p.Title, p.Detail)
}

func NewProblem(status int, title, detail, code string) Problem {
	return Problem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
		Code:   code,
	}
}
