package stepfunction

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	sfnClient "impruviService/clients/stepfunction"
	"impruviService/constants/statemachinenames"
	"log"
)

var stepFunctionClient = sfnClient.GetClient()

func StartExecution(invocationId, input string) error {
	output, err := stepFunctionClient.StartExecution(&sfn.StartExecutionInput{
		Input:           aws.String(input),
		Name:            aws.String(invocationId),
		StateMachineArn: aws.String(statemachinenames.DynamicRemindersStateMachine),
	})
	if err != nil {
		log.Printf("failed to start step function invocation with invocationId: %v, input: %v. error: %v\n", invocationId, input, err)
		return err
	}

	log.Printf("State machine invocation output: %v\n", output)
	return nil
}
