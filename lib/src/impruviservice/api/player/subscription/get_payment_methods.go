package subscription

import (
	"impruviService/exceptions"
	playerFacade "impruviService/facade/player"
	stripeFacade "impruviService/facade/stripe"
	"log"
)

type GetPaymentMethodsRequest struct {
	Token string `json:"token"` // TODO: move token into header and pass in second arg to func
}

type GetPaymentMethodsResponse struct {
	PaymentMethods []*stripeFacade.PaymentMethod `json:"paymentMethods"`
}

func GetPaymentMethods(request *GetPaymentMethodsRequest) (*GetPaymentMethodsResponse, error) {
	log.Printf("GetPaymentMethodsRequest: %+v\n", request)
	err := validateGetPaymentMethodsRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid GetPaymentMethodsRequest: %v\n", err)
		return nil, err
	}

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

func validateGetPaymentMethodsRequest(request *GetPaymentMethodsRequest) error {
	if request.Token == "" {
		return exceptions.InvalidRequestError{Message: "Token cannot be null/empty"}
	}

	return nil
}