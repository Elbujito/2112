package cors

var client *CorsClient

func init() {
	client = &CorsClient{
		name: constants.FEATURE_CORS,
	}
}

func GetClient() *CorsClient {
	return client
}
