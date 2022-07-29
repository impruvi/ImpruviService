package stripeevent

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	notificationFacade "impruviService/facade/notification"
	"log"
	"net/http"
	"os"
)

func HandleStripeEvent(apiRequest *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	event, err := webhook.ConstructEvent([]byte(apiRequest.Body), apiRequest.Headers["Stripe-Signature"], os.Getenv("WEB_HOOK_SIGNING_SECRET"))
	if err != nil {
		log.Printf("Failed to parse webhook body json: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}
	}

	log.Printf("Event type: %v\n", event.Type)

	switch event.Type {
	case "customer.subscription.deleted":
		// triggered when subscription is cancelled at the end of the period
		subscription, _ := deserializeSubscription(event.Data.Raw)
		err = handleSubscriptionDeleted(subscription)
	case "invoice.paid":
		// triggered when an invoice is paid either from a subscription being created or
		// automatically renewed
		invoice, _ := deserializeInvoice(event.Data.Raw)
		err = handleInvoicePaid(invoice)
	case "invoice.payment_failed":
		// triggered if payment for subscription fails
		invoice, _ := deserializeInvoice(event.Data.Raw)
		err = handleInvoicePaymentFailed(invoice)
	default:
		// invoice.finalization_failed
		// invoice.payment_action_required
		log.Printf("Unhandled event type: %s. Event: %+v\n", event.Type, event)
		err = notificationFacade.SendUnhandledStripeEventTypeNotifications(event.Type)
	}

	if err != nil {
		log.Printf("Unexpected error: %v\n", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}
}

func deserializeSubscription(raw json.RawMessage) (*stripe.Subscription, error) {
	var subscription stripe.Subscription
	err := json.Unmarshal(raw, &subscription)
	if err != nil {
		log.Printf("Failed to deserialize subscription: %v. raw: %v\n", err, raw)
		return nil, err
	}
	log.Printf("Deserialized subscription: %+v\n", subscription)
	return &subscription, nil
}

func deserializeInvoice(raw json.RawMessage) (*stripe.Invoice, error) {
	var invoice stripe.Invoice
	err := json.Unmarshal(raw, &invoice)
	if err != nil {
		log.Printf("Failed to deserialize invoice: %v. raw: %v\n", err, raw)
		return nil, err
	}
	log.Printf("Deserialized invoice: %+v\n", invoice)
	return &invoice, nil
}
