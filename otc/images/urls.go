package images

import (
	opentelekomcloud "github.com/cricketlong/otc-sdk-go"
)

func listDetailURL(client *opentelekomcloud.ServiceClient) string {
	return client.ServiceURL("images", "detail")
}

func getURL(client *opentelekomcloud.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}

func deleteURL(client *opentelekomcloud.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}
