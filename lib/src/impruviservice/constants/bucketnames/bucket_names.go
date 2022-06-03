package bucketnames

import "os"

var SubmissionsBucket = os.Getenv("domain") + "-impruvi-submissions-bucket"
var FeedbackBucket = os.Getenv("domain") + "-impruvi-feedback-bucket"
