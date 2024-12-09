package cors

import "github.com/Elbujito/2112/template/go-server/pkg/fx/xconstants"

var client *CorsClient

func init() {
	client = &CorsClient{
		name: xconstants.FEATURE_CORS,
	}
}

func GetClient() *CorsClient {
	return client
}
