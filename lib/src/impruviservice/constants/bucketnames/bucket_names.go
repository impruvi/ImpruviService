package bucketnames

import "os"

var MediaBucket = os.Getenv("domain") + "-impruvi-media-bucket"
