package repository

import (
	"gorm.io/gorm"
)

type Repository[T, K any] struct {
	db *gorm.DB
}

func NewRepository[T, K any](db *gorm.DB) *Repository[T, K] {
	return &Repository[T, K]{db}
}
func (r *Repository[T, K]) Create(entity *T) error {
	return r.db.Create(&entity).Error
}

func (r *Repository[T, K]) Update(entity *T, columns map[string]string) error {
	return r.db.Model(entity).Updates(UpcastMap(columns)).Error

}

func (r *Repository[T, K]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

func (r *Repository[T, K]) FindById(key K) *T {
	var e T
	rows := r.db.Find(&e, key).RowsAffected
	if rows == 0 {
		return nil
	}
	return &e
}

func UpcastMap(data map[string]string) map[string]interface{} {
	columns := make(map[string]interface{}, len(data))
	for k, v := range data {
		columns[k] = v
	}
	return columns
}
