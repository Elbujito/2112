package gzip

import "github.com/Elbujito/2112/template/go-server/pkg/fx/xconstants"

var client *GzipClient

func init() {
	client = &GzipClient{
		name: xconstants.FEATURE_GZIP,
	}
}

func GetClient() *GzipClient {
	return client
}
