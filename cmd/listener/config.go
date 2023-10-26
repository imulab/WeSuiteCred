package listener

import "github.com/urfave/cli/v2"

type config struct {
	Debug   bool
	MqttUrl string
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
	}
}
