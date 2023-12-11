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

func (cp *couponRepository) AddCoupon(CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error) {

	var coupon models.CouponResponse

	query := `
		INSERT INTO coupons (coupon_name, status, discount_rate)
		VALUES (?, ?, ?)
	`
	result := cp.DB.Exec(query, CouponName, CouponStatus, Discount)

	if result.Error != nil {
		return coupon, result.Error
	}

	coupon.CouponName = CouponName
	coupon.Status = CouponStatus
	coupon.DiscountRate = Discount

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

func (cp *couponRepository) UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int) (models.CouponResponse, error) {
	if cp.DB == nil {
		return models.CouponResponse{}, errors.New("database connection is nil")
	}
	fmt.Println("couponstatus", CouponStatus)
	fmt.Println("id", CId)

	if err := cp.DB.Exec("UPDATE coupons SET coupon_name = ?, status = ?, discount_rate = ? WHERE id = ?", CouponName, CouponStatus, Discount, CId).Error; err != nil {
		return models.CouponResponse{}, err
	}

	var updatedCoupon models.CouponResponse
	if err := cp.DB.Table("coupons").First(&updatedCoupon, CId).Error; err != nil {
		return models.CouponResponse{}, err
	}

	return updatedCoupon, nil
}

func (cp *couponRepository) CheckCoupon(coupon string) (bool, error) {

	var count int
	err := cp.DB.Raw("SELECT COUNT(*) FROM coupons WHERE coupon_name = ?", coupon).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (cp *couponRepository) CheckCouponValid(couponID int) (bool, error) {
	var status bool
	err := cp.DB.Raw("SELECT status FROM coupons WHERE id = ?", couponID).Scan(&status).Error
	if err != nil {
		return false, err
	}
	return status, nil
}

func (cp *couponRepository) FindCouponPrice(couponID int) (int, error) {
	var rate int
	err := cp.DB.Raw("SELECT discount_rate FROM coupons WHERE id = ?", couponID).Scan(&rate).Error
	if err != nil {
		return 0, err
	}
	return rate, nil
}
