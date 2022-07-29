package tablenames

import "os"

var InvitationCodesTable = os.Getenv("DOMAIN") + "-invitation-codes"
var PasswordResetCodesTable = os.Getenv("DOMAIN") + "-password-reset-codes"
var PlayersTable = os.Getenv("DOMAIN") + "-players"
var CoachesTable = os.Getenv("DOMAIN") + "-coaches"
var DrillsTable = os.Getenv("DOMAIN") + "-drills"
var SessionsTable = os.Getenv("DOMAIN") + "-sessions"
