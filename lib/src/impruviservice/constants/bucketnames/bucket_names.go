package bucketnames

import "os"

var SubmissionsBucket = os.Getenv("domain") + "-impruvi-submissions-bucket"
var FeedbackBucket = os.Getenv("domain") + "-impruvi-feedback-bucket"
var DrillsBucket = os.Getenv("domain") + "-impruvi-drills-bucket"
var HeadshotsBucket = os.Getenv("domain") + "-impruvi-headshots-bucket"
