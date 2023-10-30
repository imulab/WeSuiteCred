package suite

import "time"

// Properties is the configuration for this package.
type Properties struct {
	Id                string
	Secret            string
	AccessTokenLeeway time.Duration
}
