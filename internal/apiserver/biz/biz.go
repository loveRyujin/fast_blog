package biz

import (
	postv1 "github.com/onexstack_practice/fast_blog/internal/apiserver/biz/v1/post"
	userv1 "github.com/onexstack_practice/fast_blog/internal/apiserver/biz/v1/user"
	"github.com/onexstack_practice/fast_blog/internal/apiserver/store"
)

type IBiz interface {
	UserV1() userv1.UserBiz
	PostV1() postv1.PostBiz
}

type Biz struct {
	store store.IStore
}

var _ IBiz = (*Biz)(nil)

func NewBiz(store store.IStore) IBiz {
	return &Biz{store: store}
}

func (b *Biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store)
}

func (b *Biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}
