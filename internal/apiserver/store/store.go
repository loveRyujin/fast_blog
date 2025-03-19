package store

import (
	"context"
	"sync"

	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"
)

var (
	once sync.Once
	// 方便其它包调用已初始化好的dataStore实例
	S *dataStore
)

type IStore interface {
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
	TX(ctx context.Context, fn func(tx *gorm.DB) error) error

	User() UserStore
	Post() PostStore
}

type transactionKey struct{}

type dataStore struct {
	db *gorm.DB
}

// _ 确保dataStore实现了IStore接口
var _ IStore = (*dataStore)(nil)

// NewStore 返回一个实现了IStore接口的实例
func NewStore(db *gorm.DB) *dataStore {
	once.Do(func() {
		S = &dataStore{
			db: db,
		}
	})

	return S
}

// DB 返回一个新的数据库实例
func (s *dataStore) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := s.db
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}

	for _, whr := range wheres {
		db = whr.Where(db)
	}
	return db
}

// TX 返回一个新的事务实例
func (s *dataStore) TX(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return s.db.Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(tx)
		},
	)
}

// User 返回一个实现UserStore接口的实例
func (s *dataStore) User() UserStore {
	return newUserStore(s)
}

// Post 返回一个实现PostStore接口的实例
func (s *dataStore) Post() PostStore {
	return newPostStore(s)
}
