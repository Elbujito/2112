package gzip

import "github.com/Elbujito/2112/pkg/fx/constants"

var client *GzipClient

func init() {
	client = &GzipClient{
		name: constants.FEATURE_GZIP,
	}
}

func GetClient() *GzipClient {
	return client
}
