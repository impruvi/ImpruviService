package passwordresetcode

import (
	passwordResetCodeDao "impruviService/dao/passwordresetcode"
	"impruviService/util"
	"log"
)

func Exists(email, code string) (bool, error) {
	entries, err := passwordResetCodeDao.GetResetPasswordCodeEntries(email)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func CreateCode(email string) (string, error) {
	code := util.GenerateVerificationCode()
	log.Printf("Code: %v\n", code)
	err := passwordResetCodeDao.PutResetPasswordCodeEntry(&passwordResetCodeDao.PasswordResetCodeEntryDB{
		Email:                   email,
		CreationDateEpochMillis: util.GetCurrentTimeEpochMillis(),
		Code:                    code,
	})
	if err != nil {
		return "", err
	}
	return code, nil
}
