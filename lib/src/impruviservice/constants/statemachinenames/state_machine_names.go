package statemachinenames

import "os"

var DynamicRemindersStateMachine = "arn:aws:states:us-west-2:522042996447:stateMachine:" + os.Getenv("DOMAIN") + "-impruvi-service-dynamic-reminder-notification-state-machine"
