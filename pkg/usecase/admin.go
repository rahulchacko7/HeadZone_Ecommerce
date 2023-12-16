package usecase

import (
	domain "HeadZone/pkg/domain"
	"fmt"
	"strconv"
	"time"

	"HeadZone/pkg/helper/interfaces"
	repo "HeadZone/pkg/repository/interfaces"
	services "HeadZone/pkg/usecase/interfaces"
	"HeadZone/pkg/utils/models"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepository repo.AdminRepository
	helper          interfaces.Helper
}

func NewAdminUseCase(repo repo.AdminRepository, h interfaces.Helper) services.AdminUseCase {
	return &adminUseCase{
		adminRepository: repo,
		helper:          h,
	}
}

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {

	// getting details of the admin based on the email provided
	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	// compare password from database and that provided from admins
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	var adminDetailsResponse models.AdminDetailsResponse

	//  copy all details except password and sent it back to the front end
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	access, refresh, err := ad.helper.GenerateTokenAdmin(adminDetailsResponse)

	if err != nil {
		return domain.TokenAdmin{}, err
	}

	return domain.TokenAdmin{
		Admin:        adminDetailsResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}

func (ad *adminUseCase) BlockUser(id string) error {

	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return errors.New("invalid id")
	}

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

// unblock user
func (ad *adminUseCase) UnBlockUser(id string) error {

	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return errors.New("invalid id")
	}

	user, err := ad.adminRepository.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

func (ad *adminUseCase) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {

	if page <= 0 {
		return []models.UserDetailsAtAdmin{}, errors.New("invalid page number")
	}

	userDetails, err := ad.adminRepository.GetUsers(page)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}

func (i *adminUseCase) NewPaymentMethod(id string) error {

	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return errors.New("invalid id")
	}

	exists, err := i.adminRepository.CheckIfPaymentMethodAlreadyExists(id)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("payment method already exists")
	}

	err = i.adminRepository.NewPaymentMethod(id)
	if err != nil {
		return err
	}

	return nil

}

func (a *adminUseCase) ListPaymentMethods() ([]domain.PaymentMethod, error) {

	categories, err := a.adminRepository.ListPaymentMethods()
	if err != nil {
		return []domain.PaymentMethod{}, err
	}
	return categories, nil

}

func (a *adminUseCase) DeletePaymentMethod(id int) error {

	if id <= 0 {
		return errors.New("invalid page number")
	}

	err := a.adminRepository.DeletePaymentMethod(id)
	if err != nil {
		return err
	}
	return nil

}

func (ad *adminUseCase) DashBoard() (models.CompleteAdminDashboard, error) {
	userDetails, err := ad.adminRepository.DashBoardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := ad.adminRepository.DashBoardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	orderDetails, err := ad.adminRepository.DashBoardOrder()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	totalRevenue, err := ad.adminRepository.TotalRevenue()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	amountDetails, err := ad.adminRepository.AmountDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
		DashboardOrder:   orderDetails,
		DashboardRevenue: totalRevenue,
		DashboardAmount:  amountDetails,
	}, nil
}
func (ad *adminUseCase) SalesByDate(dayInt int, monthInt int, yearInt int) ([]models.OrderDetailsAdmin, error) {

	if yearInt >= 2020 {
		if monthInt == 0 && dayInt == 0 {

			body, err := ad.adminRepository.SalesByYear(yearInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}
			fmt.Println("body at usecase year", body)
			return body, nil
		} else if monthInt > 0 && monthInt <= 12 && dayInt == 0 {

			body, err := ad.adminRepository.SalesByMonth(yearInt, monthInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}
			fmt.Println("body at usecase month", body)
			return body, nil
		} else if monthInt > 0 && monthInt <= 12 && dayInt > 0 && dayInt <= 31 {

			body, err := ad.adminRepository.SalesByDay(yearInt, monthInt, dayInt)
			if err != nil {
				return []models.OrderDetailsAdmin{}, err
			}
			fmt.Println("body at usecase day", body)
			return body, nil
		}
	}

	return []models.OrderDetailsAdmin{}, errors.New("invalid date parameters")
}

func (ad *adminUseCase) PrintSalesReport(sales []models.OrderDetailsAdmin) (*gofpdf.Fpdf, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(31, 73, 125)
	pdf.CellFormat(0, 20, "Total Sales Report", "0", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 16)
	pdf.SetTextColor(0, 0, 0)

	for _, item := range sales {
		pdf.CellFormat(0, 10, "Product: "+item.ProductName, "0", 1, "L", false, 0, "")
		amount := strconv.FormatFloat(item.TotalAmount, 'f', 2, 64)
		pdf.CellFormat(0, 10, "Amount Sold: $"+amount, "0", 1, "L", false, 0, "")
		pdf.Ln(5)
	}

	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(150, 150, 150)

	pdf.Cell(0, 10, "Generated by HeadZone India Pvt Ltd. - "+time.Now().Format("2006-01-02 15:04:05"))

	return pdf, nil
}

func (ad *adminUseCase) CustomSalesReportByDate(startDate, endDate time.Time) (models.SalesReport, error) {
	orders, err := ad.adminRepository.CustomSalesReportByDate(startDate, endDate)
	if err != nil {
		return models.SalesReport{}, errors.New("report fetching failed")
	}
	return orders, nil
}
