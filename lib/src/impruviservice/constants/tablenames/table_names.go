package tablenames

import "os"

var InvitationCodesTable = os.Getenv("DOMAIN") + "-invitation-codes"
var PasswordResetCodesTable = os.Getenv("DOMAIN") + "-password-reset-codes"
var PlayersTable = os.Getenv("DOMAIN") + "-players"
var CoachesTable = os.Getenv("DOMAIN") + "-coaches"
var CoachApplicationsTable = os.Getenv("DOMAIN") + "-coach-applications"
var EmailListSubscriptionsTable = os.Getenv("DOMAIN") + "-email-list-subscriptions"
var DrillsTable = os.Getenv("DOMAIN") + "-drills"
var SessionsTable = os.Getenv("DOMAIN") + "-sessions"
