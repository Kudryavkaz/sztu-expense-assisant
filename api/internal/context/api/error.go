package api

type Error struct {
	Status int    `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(status int, code int, msg string, retried bool) (err *Error) {
	err = &Error{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
	return
}
