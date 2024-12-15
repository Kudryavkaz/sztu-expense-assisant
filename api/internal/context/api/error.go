package api

import "github.com/gofiber/fiber/v2"

type Error struct {
	Status int    `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
}

var (
	ErrMissToken      = NewError(fiber.StatusBadRequest, 100001, "Missing password or account")
	ErrToken          = NewError(fiber.StatusBadRequest, 100002, "Token error")
	ErrParseRequest   = NewError(fiber.StatusBadRequest, 200001, "Parse request error")
	ErrRepeatPassword = NewError(fiber.StatusBadRequest, 200002, "Password and repeat password must be same")
	ErrQueryUser      = NewError(fiber.StatusBadRequest, 200003, "Query user error")
	ErrCreateUser     = NewError(fiber.StatusBadRequest, 200004, "Create user error")
	ErrUpdateUser     = NewError(fiber.StatusBadRequest, 200005, "Update user error")
	ErrRepeatAccount  = NewError(fiber.StatusBadRequest, 200006, "Account already exists")
	ErrPassword       = NewError(fiber.StatusBadRequest, 200007, "Password error")
	ErrGenerateToken  = NewError(fiber.StatusBadRequest, 200008, "Generate token error")
	ErrGetLock        = NewError(fiber.StatusBadRequest, 200009, "Get lock error")
	ErrQueryExpense   = NewError(fiber.StatusBadRequest, 200010, "Query expense error")
	ErrCreateExpense  = NewError(fiber.StatusBadRequest, 200011, "Create expense error")
)

func (e Error) Error() string {
	return e.Msg
}

func NewError(status int, code int, msg string) (err *Error) {
	err = &Error{
		Status: status,
		Code:   code,
		Msg:    msg,
	}
	return
}
