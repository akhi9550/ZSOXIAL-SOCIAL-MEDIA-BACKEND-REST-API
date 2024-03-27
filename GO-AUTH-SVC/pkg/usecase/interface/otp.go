
package interfaces

import "github.com/akhi9550/pkg/utils/models"

type OtpUseCase interface {
	SendOtp(phone string) error
	VerifyOTP(code models.VerifyData) (models.ReponseWithToken, error)
}