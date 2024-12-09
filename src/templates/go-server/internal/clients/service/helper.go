package service

import "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: xconstants.FEATURE_SERVICE,
	}
}

func GetClient() *ServiceClient {
	return client
}
