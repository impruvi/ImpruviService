package bucketnames

import "os"

var SubmissionsBucket = os.Getenv("domain") + "-submissions"
