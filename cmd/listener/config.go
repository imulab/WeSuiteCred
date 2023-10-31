package listener

import "github.com/urfave/cli/v2"

type config struct {
	Debug       bool
	MqttUrl     string
	SuiteId     string
	SuiteSecret string
}

func (c *config) flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug mode",
			EnvVars:     []string{"WSC_DEBUG"},
			Destination: &c.Debug,
		},
		&cli.StringFlag{
			Name:        "mqtt-url",
			Usage:       "MQTT broker url",
			EnvVars:     []string{"WSC_MQTT_URL"},
			Destination: &c.MqttUrl,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "suite-id",
			Usage:       "Suite id displayed in the admin panel.",
			EnvVars:     []string{"WSC_SUITE_ID"},
			Destination: &c.SuiteId,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "suite-secret",
			Usage:       "Suite secret displayed in the admin panel.",
			EnvVars:     []string{"WSC_SUITE_SECRET"},
			Destination: &c.SuiteSecret,
			Required:    true,
		},
	}
}
