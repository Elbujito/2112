package service

import "github.com/Elbujito/2112/pkg/fx/constants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: constants.FEATURE_SERVICE,
	}
}

func GetClient() *ServiceClient {
	return client
}
