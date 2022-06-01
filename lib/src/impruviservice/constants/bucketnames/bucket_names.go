package bucketnames

import "os"

var SubmissionsBucket = os.Getenv("domain") + "-impruvi-submissions"
var FeedbackBucket = os.Getenv("domain") + "-impruvi-feedback"
