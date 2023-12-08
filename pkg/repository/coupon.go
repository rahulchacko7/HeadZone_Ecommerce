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

func (cp *couponRepository) CheckCoupon(coupon string) (bool, error) {

	var count int
	err := cp.DB.Raw("SELECT COUNT(*) FROM coupons WHERE coupon_name = ?", coupon).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (cp *couponRepository) CouponValidity(coupon string) (bool, error) {
	var status bool
	err := cp.DB.Raw("SELECT status FROM coupons WHERE coupon_name = ?", coupon).Scan(status).Error
	if err != nil {
		return false, err
	}
	return status, nil
}

func (cp *couponRepository) MinimumPurchase(coupon string) (int, error) {
	var purchase int
	err := cp.DB.Raw("SELECT minimum_price FROM coupons WHERE coupon_name = ?", coupon).Scan(purchase).Error
	if err != nil {
		return 0, err
	}
	return purchase, nil
}

func (cp *couponRepository) DiscountPercentage(coupon string) (int, error) {
	var discount int
	err := cp.DB.Raw("SELECT discount_percentage FROM coupons WHERE coupon_name = ?", coupon).Scan(discount).Error
	if err != nil {
		return 0, err
	}
	return discount, nil
}

func (cp *couponRepository) UpdateUsedCoupon(coupon string, UserId int) (bool, error) {
	var couponID uint
	err := cp.DB.Raw("SELECT id FROM coupons WHERE coupon = ?", coupon).Scan(&couponID).Error
	if err != nil {
		return false, err
	}

	err = cp.DB.Exec("INSERT INTO used_coupons (coupon_id,user_id,used) VALUES (?, ?, false)", couponID, UserId).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
