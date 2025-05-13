package validation

import (
	"context"
	"errors"

	"github.com/loveRyujin/fast_blog/internal/pkg/contextx"
	v1 "github.com/loveRyujin/fast_blog/pkg/api/apiserver/v1"
)

func (v *Validator) ValidateLoginRequest(ctx context.Context, rq *v1.LoginRequest) error {
	// Validate username
	if rq.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(rq.Username) < 4 || len(rq.Username) > 32 {
		return errors.New("username must be between 4 and 32 characters")
	}

	// Validate password
	if rq.Password == "" {
		return errors.New("password cannot be empty")
	}
	if len(rq.Password) < 8 || len(rq.Password) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	return nil
}

func (v *Validator) ValidateRefreshTokenRequest(ctx context.Context, rq *v1.RefreshTokenRequest) error {
	userID := contextx.UserID(ctx)
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	return nil
}

func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, rq *v1.ChangePasswordRequest) error {
	userID := contextx.UserID(ctx)
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	if rq.OldPassword == "" {
		return errors.New("old password cannot be empty")
	}
	if len(rq.OldPassword) < 8 || len(rq.OldPassword) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	if rq.NewPassword == "" {
		return errors.New("old password cannot be empty")
	}
	if len(rq.NewPassword) < 8 || len(rq.NewPassword) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	return nil
}

func (v *Validator) ValidateCreateUserRequest(ctx context.Context, rq *v1.CreateUserRequest) error {
	// Validate username
	if rq.Username == "" {
		return errors.New("username cannot be empty")
	}
	if len(rq.Username) < 4 || len(rq.Username) > 32 {
		return errors.New("username must be between 4 and 32 characters")
	}

	// Validate password
	if rq.Password == "" {
		return errors.New("password cannot be empty")
	}
	if len(rq.Password) < 8 || len(rq.Password) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	// Validate nickname (if provided)
	if rq.Nickname != nil && *rq.Nickname != "" {
		if len(*rq.Nickname) > 32 {
			return errors.New("nickname cannot exceed 32 characters")
		}
	}

	// Validate email
	if rq.Email == "" {
		return errors.New("email cannot be empty")
	}

	// Validate phone number
	if rq.Phone == "" {
		return errors.New("phone number cannot be empty")
	}

	return nil
}

func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *v1.UpdateUserRequest) error {
	if rq.Username != nil && (len(*rq.Username) < 4 || len(*rq.Username) > 32) {
		return errors.New("username must be between 4 and 32 characters")
	}

	if rq.Nickname != nil && *rq.Nickname != "" {
		if len(*rq.Nickname) > 32 {
			return errors.New("nickname cannot exceed 32 characters")
		}
	}

	if rq.Email != nil && *rq.Email == "" {
		return errors.New("email cannot be empty")
	}

	if rq.Phone != nil && *rq.Phone == "" {
		return errors.New("phone number cannot be empty")
	}

	return nil
}

func (v *Validator) ValidateDeleteUserRequest(ctx context.Context, rq *v1.DeleteUserRequest) error {
	userID := contextx.UserID(ctx)
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	return nil
}

func (v *Validator) ValidateGetUserRequest(ctx context.Context, rq *v1.GetUserRequest) error {
	userID := contextx.UserID(ctx)
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	return nil
}

func (v *Validator) ValidateListUserRequest(ctx context.Context, rq *v1.ListUserRequest) error {
	if rq.Offset < 0 {
		return errors.New("offset cannot be negative")
	}

	if rq.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	return nil
}
