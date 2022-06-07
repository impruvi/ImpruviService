package tablenames

import "os"

var UsersTable = os.Getenv("domain") + "-users"
var DrillsTable = os.Getenv("domain") + "-drills"
var SessionsTable = os.Getenv("domain") + "-sessions"
