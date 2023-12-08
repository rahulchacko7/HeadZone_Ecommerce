package usecase

import (
	"HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"
)

type couponUseCase struct {
	couponRepository interfaces.CouponRepository
	orderRespository interfaces.OrderRepository
	cartRepository   interfaces.CartRepository
}

func NewCouponUseCase(repository interfaces.CouponRepository, orderRepo interfaces.OrderRepository, cartRepo interfaces.CartRepository) services.CouponUseCase {
	return &couponUseCase{
		couponRepository: repository,
		orderRespository: orderRepo,
		cartRepository:   cartRepo,
	}
}

func (cp *couponUseCase) AddCoupon(CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error) {
	if Discount <= 0 || MinPurchase <= 0 {
		return models.CouponResponse{}, errors.New("discount or minimum purchase must be a +ve number")
	}
	couponResponse, err := cp.couponRepository.AddCoupon(CouponName, CouponStatus, Discount, MinPurchase)
	if err != nil {
		return models.CouponResponse{}, err
	}
	return couponResponse, nil
}

func (cp *couponUseCase) GetCoupon() ([]models.CouponResponse, error) {
	coupons, err := cp.couponRepository.GetCopupon()
	if err != nil {
		return []models.CouponResponse{}, err
	}
	return coupons, nil
}

func (cp *couponUseCase) UpdateCoupon(CId int, CouponName string, CouponStatus bool, Discount int, MinPurchase float64) (models.CouponResponse, error) {
	if Discount <= 0 || MinPurchase <= 0 || CId <= 0 {
		return models.CouponResponse{}, errors.New("discount or minimum purchase must be a +ve number")
	}
	couponResponse, err := cp.couponRepository.UpdateCoupon(CId, CouponName, CouponStatus, Discount, MinPurchase)
	if err != nil {
		return models.CouponResponse{}, err
	}
	return couponResponse, nil
}
func (repo *couponUseCase) RedeemCoupon(coupon string, UserID int) error {
	CheckCart, err := repo.orderRespository.CartExist(UserID)
	if err != nil {
		return err
	}
	if !CheckCart {
		return errors.New("Cart is empty. Cannot apply coupon")
	}

	CheckCoupon, err := repo.couponRepository.CheckCoupon(coupon)
	if err != nil {
		return err
	}
	if !CheckCoupon {
		return errors.New("Coupon does not exist")
	}

	CouponValidity, err := repo.couponRepository.CouponValidity(coupon)
	if err != nil {
		return err
	}
	if !CouponValidity {
		return errors.New("Coupon is inactive")
	}

	MinimumPurchase, err := repo.couponRepository.MinimumPurchase(coupon)
	if err != nil {
		return err
	}

	totalPriceFromCarts, err := repo.cartRepository.GetTotalPriceFromCart(UserID)
	if err != nil {
		return err
	}

	if totalPriceFromCarts < float64(MinimumPurchase) {
		return errors.New("Coupon cannot be added as the total amount is less than the minimum required for the coupon")
	}

	couponStatus, err := repo.couponRepository.UpdateUsedCoupon(coupon, UserID)
	if err != nil {
		return err
	}

	if couponStatus {
		return nil
	}
	return errors.New("Failed to add the coupon")
}
