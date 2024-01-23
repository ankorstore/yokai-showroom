package repository

import (
	"context"
	"sync"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"gorm.io/gorm"
)

// GopherRepository is the repository to handle the [model.Gopher] model database interactions.
type GopherRepository struct {
	mutex sync.Mutex
	db    *gorm.DB
}

// NewGopherRepository returns a new [GopherRepository].
func NewGopherRepository(db *gorm.DB) *GopherRepository {
	return &GopherRepository{
		db: db,
	}
}

// Find finds a [model.Gopher] by id.
func (r *GopherRepository) Find(ctx context.Context, id int) (*model.Gopher, error) {
	var gopher model.Gopher

	res := r.db.WithContext(ctx).Take(&gopher, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &gopher, nil
}

// FindAll finds all [model.Gopher].
func (r *GopherRepository) FindAll(ctx context.Context) ([]model.Gopher, error) {
	var gophers []model.Gopher

	res := r.db.WithContext(ctx).Find(&gophers)
	if res.Error != nil {
		return nil, res.Error
	}

	return gophers, nil
}

// Create creates a new [model.Gopher].
func (r *GopherRepository) Create(ctx context.Context, gopher *model.Gopher) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	res := r.db.WithContext(ctx).Create(gopher)

	return res.Error
}

// Update updates an existing [model.Gopher].
func (r *GopherRepository) Update(ctx context.Context, gopher *model.Gopher, update *model.Gopher) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	res := r.db.WithContext(ctx).Model(gopher).Updates(update)

	return res.Error
}

// Delete deletes an existing [model.Gopher].
func (r *GopherRepository) Delete(ctx context.Context, gopher *model.Gopher) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	res := r.db.WithContext(ctx).Delete(gopher)

	return res.Error
}
