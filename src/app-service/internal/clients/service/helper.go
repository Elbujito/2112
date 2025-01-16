package service

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: xconstants.FEATURE_SERVICE,
	}
}

// GetClient getters
func GetClient() *ServiceClient {
	return client
}
