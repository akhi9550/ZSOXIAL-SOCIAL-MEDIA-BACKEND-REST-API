
package interfaces

import "github.com/akhi9550/auth-svc/pkg/utils/models"

type OtpUseCase interface {
	SendOtp(phone string) error
	VerifyOTP(code models.VerifyData) (models.ReponseWithToken, error)
}