package service

import xconstants "github.com/Elbujito/2112/lib/fx/xconstants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: xconstants.FEATURE_SERVICE,
	}
}

func GetClient() *ServiceClient {
	return client
}
