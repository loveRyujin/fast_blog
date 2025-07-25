package store

import (
	"context"
	"errors"

	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"

	"github.com/loveRyujin/fast_blog/internal/apiserver/model"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
	"github.com/loveRyujin/fast_blog/internal/pkg/log"
)

// UserStore 定义了 user 模块在 store 层所实现的方法.
type UserStore interface {
	Create(ctx context.Context, obj *model.User) error
	Update(ctx context.Context, obj *model.User) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.User, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.User, error)

	UserExpansion
}

// UserExpansion 定义了用户操作的附加方法.
type UserExpansion interface{}

// userStore 是 UserStore 接口的实现.
type userStore struct {
	store *dataStore
}

// 确保 userStore 实现了 UserStore 接口.
var _ UserStore = (*userStore)(nil)

// newUserStore 创建 userStore 的实例.
func newUserStore(store *dataStore) *userStore {
	return &userStore{store: store}
}

// Create 插入一条用户记录.
func (s *userStore) Create(ctx context.Context, obj *model.User) error {
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		log.With(ctx).Errorw("Failed to insert user into database", "err", err, "user", obj)
		return errorx.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

// Update 更新用户数据库记录.
func (s *userStore) Update(ctx context.Context, obj *model.User) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.With(ctx).Errorw("Failed to update user in database", "err", err, "user", obj)
		return errorx.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

// Delete 根据条件删除用户记录.
func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.User)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.With(ctx).Errorw("Failed to delete user from database", "err", err, "conditions", opts)
		return errorx.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

// Get 根据条件查询用户记录.
func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.User, error) {
	var obj model.User
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		log.With(ctx).Errorw("Failed to retrieve user from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrUserNotFound
		}
		return nil, errorx.ErrDBRead.WithMessage(err.Error())
	}

	return &obj, nil
}

// List 返回用户列表和总数.
func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.User, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.With(ctx).Errorw("Failed to list users from database", "err", err, "conditions", opts)
		err = errorx.ErrDBRead.WithMessage(err.Error())
	}
	return
}
