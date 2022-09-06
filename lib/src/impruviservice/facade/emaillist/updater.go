package emaillist

import (
	emailListDao "impruviService/dao/emaillist"
)

func Subscribe(email string) error {
	return emailListDao.CreateSubscription(email)
}
