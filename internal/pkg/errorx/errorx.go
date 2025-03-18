package errorx

import (
	"errors"
	"fmt"
)

type Errorx struct {
	Code    int    `json:"code,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(code int, reason, message string) *Errorx {
	return &Errorx{
		Code:    code,
		Reason:  reason,
		Message: message,
	}
}

func (e *Errorx) Error() string {
	return fmt.Sprintf("error: code: %d, reason: %s, message: %s", e.Code, e.Reason, e.Message)
}

func (e *Errorx) WithMessage(message string) *Errorx {
	e.Message = message
	return e
}

func FromError(err error) *Errorx {
	if err == nil {
		return nil
	}
	if errx := new(Errorx); errors.As(err, &errx) {
		return errx
	}
	return New(ErrInternal.Code, ErrInternal.Reason, ErrInternal.Message)
}
