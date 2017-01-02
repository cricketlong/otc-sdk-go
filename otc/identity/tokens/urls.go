package tokens

import (
    opentelekomcloud "github.com/cricketlong/otc-sdk-go"
)

func tokenURL(c *opentelekomcloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
