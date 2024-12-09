package service

import "github.com/Elbujito/2112/template/go-server/pkg/fx/xconstants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: xconstants.FEATURE_SERVICE,
	}
}

func GetClient() *ServiceClient {
	return client
}
