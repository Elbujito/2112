package cors

import "github.com/Elbujito/2112/pkg/fx/constants"

var client *CorsClient

func init() {
	client = &CorsClient{
		name: constants.FEATURE_CORS,
	}
}

func GetClient() *CorsClient {
	return client
}
