package validation

import (
	"context"

	v1 "github.com/onexstack_practice/fast_blog/pkg/api/apiserver/v1"
)

func (v *Validator) ValidateCreatePostRequest(ctx context.Context, rq *v1.CreatePostRequest) error {
	return nil
}

func (v *Validator) ValidateUpdatePostRequest(ctx context.Context, rq *v1.UpdatePostRequest) error {
	return nil
}
