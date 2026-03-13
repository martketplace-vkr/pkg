package s3

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	cmpName = "s3"
)

type Client struct {
	s3  *s3.Client
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Start(_ context.Context) (err error) {
	credsProvider := credentials.NewStaticCredentialsProvider(
		c.cfg.Key,
		c.cfg.Secret,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(c.cfg.Region),
		config.WithCredentialsProvider(credsProvider),
	)
	if err != nil {
		return err
	}

	c.s3 = s3.NewFromConfig(cfg)

	return err
}

func (c *Client) Stop(_ context.Context) (err error) {
	return err
}

func (c *Client) UploadFile(ctx context.Context, key *string, reader io.Reader) error {
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    key,
		Body:   reader,
	})
	if err != nil {
		return err
	}

	return err
}

//func (c *Client) GetUploadedFilesList(ctx context.Context) (files *s3.ListObjectsV2Output, err error) {
//	result, err := c.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
//		Bucket: aws.String(c.cfg.Bucket),
//	})
//	if err != nil {
//		return files, err
//	}
//
//	return result, err
//}

func (c *Client) DownloadFile(
	ctx context.Context,
	filename string,
) (result []byte, err error) {
	file, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	defer file.Body.Close()

	bytes, err := io.ReadAll(file.Body)
	if err != nil {
		return nil, err
	}

	return bytes, err
}

func (c *Client) GetStartTimeout() time.Duration {
	return c.cfg.StartTimeout.Duration
}

func (c *Client) GetStopTimeout() time.Duration {
	return c.cfg.StopTimeout.Duration
}

func (c *Client) GetShutdownDelay() time.Duration {
	return c.cfg.ShutdownDelay.Duration
}

func (c *Client) GetName() string {
	return cmpName
}
