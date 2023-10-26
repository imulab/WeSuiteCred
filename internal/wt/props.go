package wt

import "time"

// Properties is the configuration for this package.
type Properties struct {
	SuiteId                string
	SuiteSecret            string
	SuiteAccessTokenLeeway time.Duration
}
