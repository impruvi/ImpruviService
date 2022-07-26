package subscription

import (
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
)

type GetPaymentMethodsRequest struct {
	Token string `json:"token"`
}

type GetPaymentMethodsResponse struct {
	PaymentMethods []*stripeFacade.PaymentMethod `json:"paymentMethods"`
}

func GetPaymentMethods(request *GetPaymentMethodsRequest) (*GetPaymentMethodsResponse, error) {
	player, err := playerFacade.GetPlayerFromToken(request.Token)
	if err != nil {
		return nil, err
	}

	paymentMethods, err := stripeFacade.GetPaymentMethods(player.StripeCustomerId)
	if err != nil {
		return nil, err
	}
	return &GetPaymentMethodsResponse{PaymentMethods: paymentMethods}, nil
}
