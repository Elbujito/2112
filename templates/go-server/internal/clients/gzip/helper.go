package gzip

var client *GzipClient

func init() {
	client = &GzipClient{
		name: constants.FEATURE_GZIP,
	}
}

func GetClient() *GzipClient {
	return client
}
