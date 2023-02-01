package services

import (
	"context"
	"errors"
	"loan_process/common"
	"loan_process/helpers"
	"loan_process/httpserver/repositories"
	"loan_process/httpserver/repositories/models"
	"loan_process/httpserver/request"
	"loan_process/httpserver/responses"
	"net/http"
	"time"
)

type loanApplicationSvc struct {
	customer          repositories.CustomerRepo
	province          repositories.ProvinceRepo
	loanRequest       repositories.CustomerLoanRequestRepo
	dailyLoan         repositories.DailyLoanRequestRepo
	paymentInstalment repositories.PaymentInstallmentRepo
}

func NewLoanApplicationSvc(customer repositories.CustomerRepo, province repositories.ProvinceRepo, loanRequest repositories.CustomerLoanRequestRepo, dailyLoan repositories.DailyLoanRequestRepo, paymentInstallment repositories.PaymentInstallmentRepo) *loanApplicationSvc {
	return &loanApplicationSvc{
		customer:          customer,
		province:          province,
		loanRequest:       loanRequest,
		dailyLoan:         dailyLoan,
		paymentInstalment: paymentInstallment,
	}
}

func (l *loanApplicationSvc) CreateLoanApplication(ctx context.Context, application *request.CreateLoanApplication) *responses.Response {
	_, err := l.customer.FindCustomerByKtpNumber(ctx, application.KtpNumber)
	if err == nil {
		return responses.ErrorResponse(responses.M_BAD_REQUEST, http.StatusBadRequest, errors.New("ktp number already exist"))
	}
	_, err = l.customer.FindCustomerByEmail(ctx, application.Email)
	if err == nil {
		return responses.ErrorResponse(responses.M_BAD_REQUEST, http.StatusBadRequest, errors.New("email already exist"))
	}

	date, err := time.Parse("2006-01-02", application.DateOfBirth)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}
	age := time.Now().Year() - date.Year()
	if age < 17 || age > 80 {
		return responses.ErrorResponse(responses.M_UNPROCESSABLE_ENTITY, http.StatusUnprocessableEntity, errors.New("age must between 17 and 80"))
	}

	province, err := l.province.FindProvinceByName(ctx, application.AddressProvince)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}

	if province == nil || province.Status == false {
		return responses.ErrorResponse(responses.M_UNPROCESSABLE_ENTITY, http.StatusUnprocessableEntity, errors.New("the province not allowed to loan application request"))
	}

	if !common.InArray(common.Tenor, application.Tenor) {
		return responses.ErrorResponse(responses.M_UNPROCESSABLE_ENTITY, http.StatusUnprocessableEntity, errors.New("the tenor not in right value"))
	}

	ktpImageExt := common.GetImageExtension(application.KtpImage.Filename)
	if !common.InAllowedImageExtension(ktpImageExt) {
		return responses.ErrorResponse(responses.M_UNPROCESSABLE_ENTITY, http.StatusUnprocessableEntity, errors.New("ktp image extension not allowed"))
	}
	SelfieImageExt := common.GetImageExtension(application.SelfieImage.Filename)
	if !common.InAllowedImageExtension(SelfieImageExt) {
		return responses.ErrorResponse(responses.M_UNPROCESSABLE_ENTITY, http.StatusUnprocessableEntity, errors.New("selfie image extension not allowed"))
	}

	ktpFileName := "/ktp/" + helpers.RandomString(16) + "." + ktpImageExt
	selfieFileName := "/selfie/" + helpers.RandomString(16) + "." + SelfieImageExt
	err = helpers.UploadFile(ctx, application.KtpImage, "../resources"+ktpFileName)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}
	err = helpers.UploadFile(ctx, application.SelfieImage, "../resources"+selfieFileName)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}

	customer := models.Customer{
		FullName:    application.FullName,
		KtpNumber:   application.KtpNumber,
		Gender:      application.Gender,
		DateOfBirth: date,
		PhoneNumber: application.PhoneNumber,
		Email:       application.Email,
		Nationality: application.Nationality,
		ProvinceId:  province.Id,
		KtpImage:    ktpFileName,
		SelfieImage: selfieFileName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = l.customer.SaveCustomer(ctx, &customer)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}

	customerLoanRequest := models.CustomerLoanRequest{
		CustomerId: customer.Id,
		Amount:     application.LoanAmount,
		Tenor:      application.Tenor,
		Status:     helpers.RandomSelectArrayString(common.StatusLoanRequest),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = l.loanRequest.SaveLoanRequest(ctx, &customerLoanRequest)
	if err != nil {
		return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
	}

	if customerLoanRequest.Status == common.StatusLoanRequest[0] {
		var instalments []models.PaymentInstallment
		installmentAmount := customerLoanRequest.Amount / customerLoanRequest.Tenor
		year := time.Now().Year()
		month := time.Now().Month()
		day := 10
		for i := uint(1); i <= customerLoanRequest.Tenor; i++ {
			if month == 12 && i == 2 {
				month += 1
				year += 1
			}
			dueDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			instalments = append(instalments, models.PaymentInstallment{
				CustomerId: customer.Id,
				Amount:     installmentAmount,
				Status:     "not_paid",
				DueDate:    dueDate,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}
		err = l.paymentInstalment.SavePaymentInstalment(ctx, &instalments)
		if err != nil {
			return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
		}

		dailyRequest, err := l.dailyLoan.FindDailyLoanRequestByDate(ctx, time.Now())
		if err == nil {
			dailyRequest.Request += 1
			err = l.dailyLoan.UpdateDailyLoanRequestById(ctx, dailyRequest)
			if err != nil {
				panic(err)
				return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
			}
		} else {
			dailyLoanRequest := models.DailyLoanRequest{
				CurrentDateRequest: time.Now(),
				Request:            1,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			err = l.dailyLoan.SaveDailyLoanRequest(ctx, &dailyLoanRequest)
			if err != nil {
				return responses.ErrorResponse(responses.M_INTERNAL_SERVER_ERROR, http.StatusInternalServerError, errors.New("internal server error"))
			}
		}
	}

	return responses.SuccessResponse(responses.M_CREATED, http.StatusCreated, "Success create loan application")
}
