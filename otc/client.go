package otc

import (
	"net/url"
	"reflect"
    "strings"

	opentelekomcloud "github.com/cricketlong/otc-sdk-go"
	tokens "github.com/cricketlong/otc-sdk-go/otc/identity/tokens"
)

// NewClient prepares an unauthenticated ProviderClient instance.
// Most users will probably prefer using the AuthenticatedClient function instead.
// This is useful if you wish to explicitly control the version of the identity service that's used for authentication explicitly,
// for example.
func NewClient(endpoint string) (*opentelekomcloud.ProviderClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	hadPath := u.Path != ""
	u.Path, u.RawQuery, u.Fragment = "", "", ""
	base := u.String()

	endpoint = opentelekomcloud.NormalizeURL(endpoint)
	base = opentelekomcloud.NormalizeURL(base)

	if hadPath {
		return &opentelekomcloud.ProviderClient{
			IdentityBase:     base,
			IdentityEndpoint: endpoint,
		}, nil
	}

	return &opentelekomcloud.ProviderClient{
		IdentityBase:     base,
		IdentityEndpoint: "",
	}, nil
}

// AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint specified by options, acquires a token, and
// returns a Client instance that's ready to operate.
// It first queries the root identity endpoint to determine which versions of the identity service are supported, then chooses
// the most recent identity service available to proceed.
func AuthenticatedClient(options opentelekomcloud.AuthOptions) (*opentelekomcloud.ProviderClient, error) {
	client, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate or re-authenticate against the most recent identity service supported at the provided endpoint.
func Authenticate(client *opentelekomcloud.ProviderClient, options opentelekomcloud.AuthOptions) error {
	//versions := []*utils.Version{
		//{ID: v20, Priority: 20, Suffix: "/v2.0/"},
		//{ID: v30, Priority: 30, Suffix: "/v3/"},
	//}

	//chosen, endpoint, err := utils.ChooseVersion(client, versions)
	//if err != nil {
		//return err
	//}

    endpoint := client.IdentityEndpoint
    if !strings.HasSuffix(endpoint, "/") {
        endpoint = endpoint + "/"
    }

	return v3auth(client, endpoint, &options, opentelekomcloud.EndpointOpts{})
}

// AuthenticateV2 explicitly authenticates against the identity v2 endpoint.
//func AuthenticateV2(client *gophercloud.ProviderClient, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	//return v2auth(client, "", options, eo)
//}

/*
func v2auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	v2Client, err := NewIdentityV2(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	v2Opts := tokens2.AuthOptions{
		IdentityEndpoint: options.IdentityEndpoint,
		Username:         options.Username,
		Password:         options.Password,
		TenantID:         options.TenantID,
		TenantName:       options.TenantName,
		AllowReauth:      options.AllowReauth,
		TokenID:          options.TokenID,
	}

	result := tokens2.Create(v2Client, v2Opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v2auth(client, endpoint, options, eo)
		}
	}
	client.TokenID = token.ID
	client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return V2EndpointURL(catalog, opts)
	}

	return nil
}
*/

// AuthenticateV3 explicitly authenticates against the identity v3 service.
func AuthenticateV3(client *opentelekomcloud.ProviderClient, options tokens.AuthOptionsBuilder, eo opentelekomcloud.EndpointOpts) error {
	return v3auth(client, "", options, eo)
}

func v3auth(client *opentelekomcloud.ProviderClient, endpoint string, opts tokens.AuthOptionsBuilder, eo opentelekomcloud.EndpointOpts) error {
	// Override the generated service endpoint with the one returned by the version endpoint.
	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	result := tokens.Create(v3Client, opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	client.TokenID = token.ID

	if opts.CanReauth() {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v3auth(client, endpoint, opts, eo)
		}
	}
	client.EndpointLocator = func(opts opentelekomcloud.EndpointOpts) (string, error) {
		return V3EndpointURL(catalog, opts)
	}

	return nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the v2 identity service.
/*
func NewIdentityV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	endpoint := client.IdentityBase + "v2.0/"
	var err error
	if !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults("identity")
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}, nil
}
*/

// NewIdentityV3 creates a ServiceClient that may be used to access the v3 identity service.
func NewIdentityV3(client *opentelekomcloud.ProviderClient, eo opentelekomcloud.EndpointOpts) (*opentelekomcloud.ServiceClient, error) {
	endpoint := client.IdentityBase + "v3/"
	var err error
	if !reflect.DeepEqual(eo, opentelekomcloud.EndpointOpts{}) {
		eo.ApplyDefaults("identity")
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &opentelekomcloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}, nil
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1 object storage package.
/*
func NewObjectStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("object-store")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
*/

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute package.
func NewComputeV2(client *opentelekomcloud.ProviderClient, eo opentelekomcloud.EndpointOpts) (*opentelekomcloud.ServiceClient, error) {
	eo.ApplyDefaults("compute")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &opentelekomcloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network package.
/*
func NewNetworkV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("network")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       url,
		ResourceBase:   url + "v2.0/",
	}, nil
}
*/

// NewBlockStorageV1 creates a ServiceClient that may be used to access the v1 block storage service.
/*
func NewBlockStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("volume")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
/*

// NewBlockStorageV2 creates a ServiceClient that may be used to access the v2 block storage service.
/*
func NewBlockStorageV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("volumev2")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
*/

// NewSharedFileSystemV2 creates a ServiceClient that may be used to access the v2 shared file system service.
/*
func NewSharedFileSystemV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("sharev2")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
*/

// NewCDNV1 creates a ServiceClient that may be used to access the OpenStack v1
// CDN service.
/*
func NewCDNV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("cdn")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
*/

// NewOrchestrationV1 creates a ServiceClient that may be used to access the v1 orchestration service.
/*
func NewOrchestrationV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("orchestration")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
*/

// NewDBV1 creates a ServiceClient that may be used to access the v1 DB service.
/*
func NewDBV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("database")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
*/

// NewImageServiceV2 creates a ServiceClient that may be used to access the v2 image service.
func NewImageServiceV2(client *opentelekomcloud.ProviderClient, eo opentelekomcloud.EndpointOpts) (*opentelekomcloud.ServiceClient, error) {
	eo.ApplyDefaults("image")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &opentelekomcloud.ServiceClient{ProviderClient: client,
		Endpoint:     url,
		ResourceBase: url + "v2/"}, nil
}
