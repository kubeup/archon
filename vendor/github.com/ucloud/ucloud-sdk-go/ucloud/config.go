package ucloud

import (
	"os"
	"time"

	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

const (
	DefaultRetries        = -1
	APIVersion     string = "v2"
	APIBaseURL     string = "https://api.ucloud.cn/"
)

var DefaultCredential, _ = auth.LoadKeyPairFromEnv()

var DefaultConfig = NewConfig().
	WithCredentials(&DefaultCredential).
	WithRegion(os.Getenv("UCLOUD_REGION")).
	WithMaxRetries(DefaultRetries).
	WithProjectID("").
	WithSleepDelay(time.Sleep)

// A Config provides service configuration for service clients. By default,
// all clients will use the {defaults.DefaultConfig} structure.
// TODO: max retries and timeout should be added
type Config struct {
	Credentials *auth.KeyPair

	Region string

	MaxRetries int

	ProjectID string

	SleepDelay func(time.Duration)

	HTTPHeader map[string]string
}

// NewConfig returns a new Config pointer that can be chained with builder methods to
// set multiple configuration values inline without using pointers.
//
//     svc := uhost.New(ucloud.NewConfig().WithRegion("cn-bj-2").WithMaxRetries(10))
//
func NewConfig() *Config {
	return &Config{}
}

// WithCredentials sets a config Credentials value returning a Config pointer
// for chaining.
func (c *Config) WithCredentials(creds *auth.KeyPair) *Config {
	c.Credentials = creds
	return c
}

// WithRegion sets a config Region value returning a Config pointer for
// chaining.
func (c *Config) WithRegion(region string) *Config {
	c.Region = region
	return c
}

// WithMaxRetries sets a config MaxRetries value returning a Config pointer
// for chaining.
func (c *Config) WithMaxRetries(max int) *Config {
	c.MaxRetries = max
	return c
}

// WithProjectID sets a config ProjectID value returning a Config pointer
// for chaining
func (c *Config) WithProjectID(projectID string) *Config {
	c.ProjectID = projectID
	return c
}

// WithSleepDelay overrides the function used to sleep while waiting for the
// next retry. Defaults to time.Sleep.
func (c *Config) WithSleepDelay(fn func(time.Duration)) *Config {
	c.SleepDelay = fn
	return c
}

// Merge returns a new Config with the other Config's attribute values merged into
// this Config. If the other Config's attribute is nil it will not be merged into
// the new Config to be returned.
func (c Config) Merge(other *Config) *Config {
	if other == nil {
		return &c
	}

	dst := c

	if other.Credentials != nil {
		dst.Credentials = other.Credentials
	}

	if other.Region != "" {
		dst.Region = other.Region
	}

	if other.SleepDelay != nil {
		dst.SleepDelay = other.SleepDelay
	}

	if other.MaxRetries != -1 {
		dst.MaxRetries = other.MaxRetries
	}

	if len(other.ProjectID) != 0 {
		dst.ProjectID = other.ProjectID
	}

	if len(other.HTTPHeader) != 0 {
		dst.HTTPHeader = other.HTTPHeader
	}

	return &dst
}

// Copy will return a shallow copy of the Config object.
func (c Config) Copy() *Config {
	dst := c
	return &dst
}
