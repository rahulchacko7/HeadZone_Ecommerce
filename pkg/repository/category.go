package repository

import (
	"HeadZone/pkg/domain"
	interfaces "HeadZone/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type categoryRespository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) interfaces.CategoryRepository {
	return &categoryRespository{DB}
}

func (p *categoryRespository) AddCategory(c domain.Category) (domain.Category, error) {

	var b domain.Category
	err := p.DB.Raw("INSERT INTO categories (category) VALUES (?) RETURNING *", c.Category).Scan(&b).Error
	if err != nil {
		return domain.Category{}, err
	}

	return b, nil
}

func (c *categoryRespository) GetCategories() ([]domain.Category, error) {
	var Model []domain.Category
	err := c.DB.Raw("SELECT * FROM categories").Scan(&Model).Error
	if err != nil {
		return []domain.Category{}, err
	}

	return Model, nil
}
