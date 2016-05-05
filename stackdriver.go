package stackdriver

import (
	"log"
	"time"

	"github.com/AMeng/stackdriver"
	"github.com/rcrowley/go-metrics"
)

type Config struct {
	APIKey     string
	Prefix     string
	InstanceID string
	Log        *log.Logger
}

func Send(r metrics.Registry, d time.Duration, config Config) {
	client := stackdriver.NewStackdriverClient(config.APIKey)
	for range time.Tick(d) {
		if err := send(r, client, config); err != nil {
			config.Log.Println("stackdriver:", err)
		}
	}
}

func send(r metrics.Registry, client *stackdriver.StackdriverClient, config Config) error {
	m := stackdriver.NewGatewayMessage()
	fillMetrics(r, &m, config)
	return client.Send(m)
}

func fillMetrics(r metrics.Registry, msg *stackdriver.GatewayMessage, config Config) {
	now := time.Now().Unix()
	r.Each(func(name string, i interface{}) {
		if config.Prefix != "" {
			name = config.Prefix + "." + name
		}
		switch m := i.(type) {
		case metrics.Counter:
			msg.CustomMetric(name, config.InstanceID, now, m.Count())
		case metrics.Gauge:
			msg.CustomMetric(name, config.InstanceID, now, m.Value())
		case metrics.GaugeFloat64:
			msg.CustomMetric(name, config.InstanceID, now, m.Value())
		case metrics.Histogram:
			msg.CustomMetric(name, config.InstanceID, now, m.Mean())
		case metrics.Meter:
			msg.CustomMetric(name, config.InstanceID, now, m.Rate1())
		case metrics.Timer:
			msg.CustomMetric(name, config.InstanceID, now, m.Mean())
		}
	})
}
