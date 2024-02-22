package repository

import (
	"HeadZone/pkg/utils/models"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCheckUserAvailability(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		stub func(mock sqlmock.Sqlmock)
		want bool
	}{
		{
			name: "successful, user available",
			arg:  "rahul2@gmail.com",
			stub: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select count(*) from users where email='rahul2@gmail.com'")).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want: true,
		},
	}
	for _, tt := range tests {
		mockDB, mockSql, _ := sqlmock.New()

		DB, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: mockDB,
		}), &gorm.Config{})

		userRepository := NewUserRepository(DB)
		tt.stub(mockSql)

		result := userRepository.CheckUserAvailability(tt.arg)
		assert.Equal(t, tt.want, result)
	}
}

func TestUserSignUp(t *testing.T) {
	type args struct {
		input models.UserDetails
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mockSQL sqlmock.Sqlmock)
		want       models.UserDetailsResponse
		wantErr    error
	}{
		{
			name: "Successfully user signed up",
			args: args{
				input: models.UserDetails{
					Name:     "Rahul",
					Email:    "rahulchacko888@gmail.com",
					Password: "12345",
					Phone:    "9867327710",
				},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users \(name, email, password, phone\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, name, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("Rahul", "rahulchacko888@gmail.com", "12345", "9867327710").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone"}).
						AddRow(1, "Rahul", "rahulchacko888@gmail.com", "9867327710"))
			},
			want: models.UserDetailsResponse{
				Id:    1,
				Name:  "Rahul",
				Email: "rahulchacko888@gmail.com",
				Phone: "9867327710",
			},
			wantErr: nil,
		},
		{
			name: "Error signing up user",
			args: args{
				input: models.UserDetails{},
			},
			beforeTest: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `INSERT INTO users \(name, email, password, phone\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id, name, email, phone`
				mockSQL.ExpectQuery(expectedQuery).
					WithArgs("", "", "", "").
					WillReturnError(errors.New("email should be unique"))
			},
			want:    models.UserDetailsResponse{},
			wantErr: errors.New("email should be unique"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSql, _ := sqlmock.New()
			defer mockDB.Close()
			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			tt.beforeTest(mockSql)
			u := NewUserRepository(gormDB)
			got, err := u.UserSignUp(tt.args.input)
			assert.Equal(t, tt.wantErr, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.UserSignUp() = %v, want %v", got, tt.want)
			}
		})

	}
}

func Test_FindUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		args    models.UserLogin
		stub    func(sqlmock.Sqlmock)
		want    models.UserSignInResponse
		wantErr error
	}{
		{
			name: "success",
			args: models.UserLogin{
				Email:    "rahulchacko123@gmail.com",
				Password: "1234",
			},
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `^SELECT \* FROM users(.+)$`
				mockSQL.ExpectQuery(expectedQuery).WithArgs("rahulchacko123@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "email", "phone", "password"}).
						AddRow(1, 1, "Rahul", "rahulchacko123@gmail.com", "+917012493965", "4321"))
			},

			want: models.UserSignInResponse{
				Id:       1,
				UserID:   1,
				Name:     "Rahul",
				Email:    "rahulchacko123@gmail.com",
				Phone:    "+917012493965",
				Password: "4321",
			},
			wantErr: nil,
		},
		{
			name: "error",
			args: models.UserLogin{
				Email:    "rahulchacko123@gmail.com",
				Password: "1234",
			},
			stub: func(mockSQL sqlmock.Sqlmock) {

				expectedQuery := `^SELECT \* FROM users(.+)$`

				mockSQL.ExpectQuery(expectedQuery).WithArgs("rahulchacko123@gmail.com").
					WillReturnError(errors.New("new error"))

			},

			want:    models.UserSignInResponse{},
			wantErr: errors.New("error checking user details"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.stub(mockSQL)

			u := NewUserRepository(gormDB)

			result, err := u.FindUserByEmail(tt.args)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}

}

type id struct {
	id int
}

func Test_GetUserDetails(t *testing.T) {
	tests := []struct {
		name    string
		args    id
		stub    func(sqlmock.Sqlmock)
		want    models.UserDetailsResponse
		wantErr error
	}{
		{
			name: "Success",
			args: id{
				id: 1,
			},
			stub: func(mockSQL sqlmock.Sqlmock) {
				expectedQuery := `^select id,name,email,phone from users where id=?`
				mockSQL.ExpectQuery(expectedQuery).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone"}).AddRow(1, "Arun CM", "aruncm@gmale.com", "+1234567890"))
			},
			want: models.UserDetailsResponse{
				Id:    1,
				Name:  "Arun CM",
				Email: "aruncm@gmale.com",
				Phone: "+1234567890", // Corrected phone number
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			gormDB, _ := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})

			tt.stub(mockSQL)

			ad := userDatabase{DB: gormDB}

			result, err := ad.GetUserDetails(tt.args.id)

			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
