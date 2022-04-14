package goaurestclientprometheus

import (
	"net/url"
	"regexp"
)

var hostnameSanitizer *regexp.Regexp

func ClientNameFromRequestUrl(requestUrl string) string {
	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "unknown"
	}

	hostName := parsedUrl.Hostname()
	return hostnameSanitizer.ReplaceAllString(hostName, "")
}

func OutcomeFromStatus(status int) string {
	if status >= 100 && status < 200 {
		return "INFORMATIONAL"
	} else if status >= 200 && status < 300 {
		return "SUCCESS"
	} else if status >= 300 && status < 400 {
		return "REDIRECTION"
	} else if status >= 400 && status < 500 {
		return "CLIENT_ERROR"
	} else if status >= 500 && status < 600 {
		return "SERVER_ERROR"
	}
	return "UNKNOWN"
}

func SetupCommon() {
	if hostnameSanitizer == nil {
		hostnameSanitizer = regexp.MustCompile("[^a-zA-Z0-9.-]+")
	}
}
