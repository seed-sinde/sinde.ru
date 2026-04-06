package auth

import "time"

func BuildActionEmailForTest(
	purpose string,
	ttl time.Duration,
	actionURL string,
	publicBaseURL string,
) (string, string, string) {
	return buildActionEmail(purpose, ttl, actionURL, publicBaseURL)
}
