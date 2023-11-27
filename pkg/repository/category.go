package repository

import (
	"HeadZone/pkg/domain"
	interfaces "HeadZone/pkg/repository/interfaces"
	"errors"
	"strconv"

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

func (p *categoryRespository) CheckCategory(current string) (bool, error) {
	var i int
	err := p.DB.Raw("SELECT COUNT(*) FROM categories WHERE category=?", current).Scan(&i).Error
	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, err
	}

	return true, err
}

func (p *categoryRespository) UpdateCategory(current, new string) (domain.Category, error) {

	// Check the database connection
	if p.DB == nil {
		return domain.Category{}, errors.New("database connection is nil")
	}

	// Update the category
	if err := p.DB.Exec("UPDATE categories SET category = $1 WHERE category = $2", new, current).Error; err != nil {
		return domain.Category{}, err
	}

	// Retrieve the updated category
	var newcat domain.Category
	if err := p.DB.First(&newcat, "category = ?", new).Error; err != nil {
		return domain.Category{}, err
	}

	return newcat, nil
}

func (p *categoryRespository) DeleteCategory(catergoryID string) error {

	id, err := strconv.Atoi(catergoryID)

	if err != nil || id <= 0 {
		return errors.New("invalid category ID")
	}

	if err != nil {
		return errors.New("converting into integers is not happen")
	}

	result := p.DB.Exec("DELETE FROM categories WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("now rows with that id exist")
	}
	return nil
}
