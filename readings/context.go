package readings

import (
	"context"
	"sync"
	"time"

	"github.com/op/go-logging"

	"sensorsys/config"
	"sensorsys/model"
	"sensorsys/readings/receiver"
	"sensorsys/readings/sensor"
)

type Context struct {
	Parent context.Context
	Of     string
	Logger *logging.Logger
	WaitGroup *sync.WaitGroup
	Config config.Config
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Parent: parent,
	}
}

func (c *Context) ForSensor(s sensor.Sensor) *sensor.Context {
	return &sensor.Context{
		Parent: c,
		Logger: c.Logger,
		SensorID: s.ID(),
		Pipe: make(model.MetricReadingsPipe),
	}
}

func (c *Context) ForRequest(metrics []model.Metric) *receiver.Context {
	pipe := make(model.MetricReadingsPipe)
	for _, metric := range metrics {
		pipe[metric] = make(chan model.MetricReading, 3)
	}
	return &receiver.Context {
		Parent: c,
		Logger: c.Logger,
		Pipe: pipe,
	}
}

func (c *Context) SetLogger(logger *logging.Logger) *Context {
	c.Logger = logger
	return c
}

func (c *Context) SetConfig(config config.Config) *Context {
	c.Config = config
	return c
}


func (c *Context) Error(err error) {
	if err != nil {
		c.Logger.Errorf("%v: %v", c.Of, err)
	}
}

func (c *Context) Info(info string) {
	c.Logger.Infof("%v: %v", c.Of, info)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Parent.Deadline()
}

func (c *Context) Done() <- chan struct{} {
	return c.Parent.Done()
}

func (c *Context) Err() error {
	return c.Parent.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.Parent.Value(key)
}