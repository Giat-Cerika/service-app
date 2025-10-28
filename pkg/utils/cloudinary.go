package utils

import (
	"net/url"
	"path"
	"strings"
)

func ExtractPublicIDFromCloudinaryURL(imageURL string) string {
	parts := strings.Split(imageURL, "/upload/")
	if len(parts) != 2 {
		return ""
	}

	afterUpload := parts[1]

	slashIndex := strings.Index(afterUpload, "/")
	if slashIndex == -1 {
		return ""
	}

	pathWithFile := afterUpload[slashIndex+1:]

	decodedPath, err := url.PathUnescape(pathWithFile)
	if err != nil {
		return ""
	}

	publicID := strings.TrimSuffix(decodedPath, path.Ext(decodedPath))
	return publicID
}
