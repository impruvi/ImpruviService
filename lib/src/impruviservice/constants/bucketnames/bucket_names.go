package bucketnames

import "os"

var MediaBucket = os.Getenv("DOMAIN") + "-impruvi-media-bucket"
