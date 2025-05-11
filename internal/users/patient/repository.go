package patient

import "gorm.io/gorm"

type Repository interface {
	Create(user *Patient) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *Patient) error {
	r.db.Create(user)
	return nil
}
