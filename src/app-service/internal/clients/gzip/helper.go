package gzip

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

var client *GzipClient

func init() {
	client = &GzipClient{
		name: xconstants.FEATURE_GZIP,
	}
}

func GetClient() *GzipClient {
	return client
}
