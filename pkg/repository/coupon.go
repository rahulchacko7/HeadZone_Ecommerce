package repository

import (
	"HeadZone/pkg/repository/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type couponRepository struct {
	DB *gorm.DB
}

func NewCouponRepository(DB *gorm.DB) interfaces.CouponRepository {
	return &couponRepository{
		DB: DB,
	}
}

func (cp *couponRepository) AddCoupon(CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error) {

	var coupon models.CouponResponse

	query := `
		INSERT INTO coupons (coupon_name, status, discount_percentage, minimum_price)
		VALUES (?, ?, ?, ?)
	`
	result := cp.DB.Exec(query, CouponName, CouponStatus, Discount, MinPurchase)

	if result.Error != nil {
		return coupon, result.Error
	}

	coupon.CouponName = CouponName
	coupon.Status = CouponStatus
	coupon.DiscountPercentage = Discount
	coupon.MinimumPrice = MinPurchase

	return coupon, nil
}

func (cp *couponRepository) GetCopupon() ([]models.CouponResponse, error) {
	var coupon []models.CouponResponse
	err := cp.DB.Raw("SELECT * FROM coupons").Scan(&coupon).Error
	if err != nil {
		return []models.CouponResponse{}, err
	}
	return coupon, nil
}

func (cp *couponRepository) UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error) {
	if cp.DB == nil {
		return models.CouponResponse{}, errors.New("database connection is nil")
	}
	fmt.Println("couponstatus", CouponStatus)
	fmt.Println("id", CId)

	if err := cp.DB.Exec("UPDATE coupons SET coupon_name = ?, status = ?, discount_percentage = ?, minimum_price = ? WHERE id = ?", CouponName, CouponStatus, Discount, MinPurchase, CId).Error; err != nil {
		return models.CouponResponse{}, err
	}

	var updatedCoupon models.CouponResponse
	if err := cp.DB.Table("coupons").First(&updatedCoupon, CId).Error; err != nil {
		return models.CouponResponse{}, err
	}

	return updatedCoupon, nil
}
