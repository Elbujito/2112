package cors

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

var client *CorsClient

func init() {
	client = &CorsClient{
		name: xconstants.FEATURE_CORS,
	}
}

func GetClient() *CorsClient {
	return client
}
