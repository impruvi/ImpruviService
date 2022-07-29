package appversion

import (
	"impruviService/exceptions"
	"log"
	"strconv"
	"strings"
)

type IsAppVersionCompatibleRequest struct {
	Version string `json:"version"`
}

type IsAppVersionCompatibleResponse struct {
	IsCompatible                       bool   `json:"isCompatible"`
	NewVersionPreviewImageFileLocation string `json:"newVersionPreviewImageFileLocation"`
}

const oldestAllowedMajorVersion = 1
const oldestAllowedMinorVersion = 0
const oldestAllowedPatchVersion = 1
const appPreviewFileLocation = "https://prod-impruvi-media-bucket.s3.us-west-2.amazonaws.com/apppreview/1.0.1"

func IsAppVersionCompatible(request *IsAppVersionCompatibleRequest) (*IsAppVersionCompatibleResponse, error) {
	log.Printf("IsAppVersionCompatibleRequest: %+v\n", request)
	err := validateIsAppVersionCompatibleRequest(request)
	if err != nil {
		log.Printf("[WARN] invalid IsAppVersionCompatibleRequest: %v\n", err)
		return nil, err
	}

	versionParts := strings.Split(request.Version, ".")
	majorVersion, _ := strconv.Atoi(versionParts[0])
	minorVersion, _ := strconv.Atoi(versionParts[1])
	patchVersion, _ := strconv.Atoi(versionParts[2])

	var isCompatible = true
	if majorVersion < oldestAllowedMajorVersion {
		isCompatible = false
	} else if majorVersion == oldestAllowedMajorVersion {
		if minorVersion < oldestAllowedMinorVersion {
			isCompatible = false
		} else if minorVersion == oldestAllowedMinorVersion {
			isCompatible = patchVersion >= oldestAllowedPatchVersion
		}
	}

	return &IsAppVersionCompatibleResponse{
		IsCompatible:                       isCompatible,
		NewVersionPreviewImageFileLocation: appPreviewFileLocation,
	}, nil
}

func validateIsAppVersionCompatibleRequest(request *IsAppVersionCompatibleRequest) error {
	if request.Version == "" {
		return exceptions.InvalidRequestError{Message: "Version cannot be null/empty"}
	}
	versionParts := strings.Split(request.Version, ".")
	if len(versionParts) != 3 {
		return exceptions.InvalidRequestError{Message: "Version must follow the pattern <majorVersion>.<minorVersion>.<patchVersion>"}
	}

	_, err := strconv.Atoi(versionParts[0])
	if err != nil {
		log.Printf("Cannot convert %v to major version", versionParts[0])
		return exceptions.InvalidRequestError{Message: "Invalid major version number"}
	}
	_, err = strconv.Atoi(versionParts[1])
	if err != nil {
		log.Printf("Cannot convert %v to minor version", versionParts[1])
		return exceptions.InvalidRequestError{Message: "Invalid minor version number"}
	}
	_, err = strconv.Atoi(versionParts[2])
	if err != nil {
		log.Printf("Cannot convert %v to patch version", versionParts[2])
		return exceptions.InvalidRequestError{Message: "Invalid patch version number"}
	}

	return nil
}
