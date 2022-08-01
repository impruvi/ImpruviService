package mediaconvertqueue

import "os"

var MediaconvertQueue = "arn:aws:mediaconvert:us-west-2:522042996447:queues/" + os.Getenv("DOMAIN") + "-impruvi-service-queue"
