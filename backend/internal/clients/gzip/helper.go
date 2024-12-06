package gzip

import xconstants "github.com/Elbujito/2112/lib/fx/xconstants"

var client *GzipClient

func init() {
	client = &GzipClient{
		name: xconstants.FEATURE_GZIP,
	}
}

func GetClient() *GzipClient {
	return client
}
