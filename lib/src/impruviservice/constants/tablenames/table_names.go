package tablenames

import "os"

var UsersTable = os.Getenv("domain") + "-users"
var CoachesTable = os.Getenv("domain") + "-coaches"
var DrillsTable = os.Getenv("domain") + "-drills"
var SessionsTable = os.Getenv("domain") + "-sessions"
