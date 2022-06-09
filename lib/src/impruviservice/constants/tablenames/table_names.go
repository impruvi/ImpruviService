package tablenames

import "os"

var InvitationCodesTable = os.Getenv("domain") + "-invitation-codes"
var PlayersTable = os.Getenv("domain") + "-players"
var CoachesTable = os.Getenv("domain") + "-coaches"
var DrillsTable = os.Getenv("domain") + "-drills"
var SessionsTable = os.Getenv("domain") + "-sessions"
