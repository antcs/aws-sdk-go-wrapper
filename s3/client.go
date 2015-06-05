// S3 client setting

package s3

import (
	SDK "github.com/aws/aws-sdk-go/service/s3"

	"github.com/evalphobia/aws-sdk-go-wrapper/auth"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

const (
	s3ConfigSectionName = "s3"
	defaultRegion       = "us-east-1"
	defaultEndpoint     = "http://localhost:4567"
	defaultBucketPrefix = "dev-"
)

// wrapped struct for S3
type AmazonS3 struct {
	buckets map[string]*Bucket
	client  *SDK.S3
}

// Create new AmazonS3 struct
func NewClient() *AmazonS3 {
	s := &AmazonS3{}
	s.buckets = make(map[string]*Bucket)
	region := config.GetConfigValue(s3ConfigSectionName, "region", "")
	awsConf := auth.NewConfig(region)
	endpoint := config.GetConfigValue(s3ConfigSectionName, "endpoint", "")
	switch {
	case endpoint != "":
		awsConf.Endpoint = endpoint
		awsConf.S3ForcePathStyle = true
	case region == "":
		awsConf.Region = defaultRegion
		awsConf.Endpoint = defaultEndpoint
		awsConf.S3ForcePathStyle = true
	}

	s.client = SDK.New(awsConf)
	return s
}

// get bucket
func (s *AmazonS3) GetBucket(bucket string) *Bucket {
	prefix := config.GetConfigValue(s3ConfigSectionName, "prefix", defaultBucketPrefix)
	bucketName := prefix + bucket

	// get the bucket from cache
	b, ok := s.buckets[bucketName]
	if ok {
		return b
	}

	b = &Bucket{}
	b.client = s.client
	b.name = bucketName
	s.buckets[bucketName] = b
	return b
}
