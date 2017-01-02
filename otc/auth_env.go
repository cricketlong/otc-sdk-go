package otc

import (
	"os"

	opentelekomcloud "github.com/cricketlong/otc-sdk-go"
)

var nilOptions = opentelekomcloud.AuthOptions{}

// AuthOptionsFromEnv fills out an identity.AuthOptions structure with the settings found on the various OpenStack
// OS_* environment variables.  The following variables provide sources of truth: OS_AUTH_URL, OS_USERNAME,
// OS_PASSWORD, OS_TENANT_ID, and OS_TENANT_NAME.  Of these, OS_USERNAME, OS_PASSWORD, and OS_AUTH_URL must
// have settings, or an error will result.  OS_TENANT_ID and OS_TENANT_NAME are optional.
func AuthOptionsFromEnv() (opentelekomcloud.AuthOptions, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	username := os.Getenv("OS_USERNAME")
	userID := os.Getenv("OS_USERID")
	password := os.Getenv("OS_PASSWORD")
	tenantID := os.Getenv("OS_TENANT_ID")
	tenantName := os.Getenv("OS_TENANT_NAME")
	domainID := os.Getenv("OS_DOMAIN_ID")
	domainName := os.Getenv("OS_DOMAIN_NAME")

	if authURL == "" {
		err := opentelekomcloud.ErrMissingInput{Argument: "authURL"}
		return nilOptions, err
	}

	if username == "" && userID == "" {
		err := opentelekomcloud.ErrMissingInput{Argument: "username"}
		return nilOptions, err
	}

	if password == "" {
		err := opentelekomcloud.ErrMissingInput{Argument: "password"}
		return nilOptions, err
	}

	ao := opentelekomcloud.AuthOptions{
		IdentityEndpoint: authURL,
		UserID:           userID,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
		TenantName:       tenantName,
		DomainID:         domainID,
		DomainName:       domainName,
	}

	return ao, nil
}
