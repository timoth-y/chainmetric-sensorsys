package sensors

import (
	"fmt"
	"math"

	"github.com/d2r2/go-i2c"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings/sensor"
)

const(
	MAX44009_APP_START = 0x03
)

type MAX44009 struct {
	addr uint8
	bus int
	i2c *i2c.I2C
}

func NewMAX44009(addr uint8, bus int) *MAX44009 {
	return &MAX44009{
		addr: addr,
		bus: bus,
	}
}

func (s *MAX44009) ID() string {
	return "MAX44009"
}

func (s *MAX44009) Init() (err error) {
	s.i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	if !s.Verify() {
		return fmt.Errorf("not MAX44009 sensorType")
	}

	_, err = s.i2c.WriteBytes([]byte{MAX44009_APP_START}); if err != nil {
		return
	}

	return
}

func (s *MAX44009) Read() (lux float64, err error) {
	var buffer = make([]byte, 2)
	_, err = s.i2c.ReadBytes(buffer); if err != nil {
		return math.NaN(), err
	}

	lux = dataToLuminance(buffer)
	return
}

func (s *MAX44009) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Luminosity).WriteWithError(s.Read())
}

func (s *MAX44009) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Luminosity,
	}
}

func (s *MAX44009) Verify() bool {
	return true // TODO verify by device ID
}

func (s *MAX44009) Active() bool {
	return s.i2c != nil
}

func (s *MAX44009) Close() error {
	defer s.clean()
	return s.i2c.Close()
}

func dataToLuminance(d []byte) float64 {
	exponent := int((d[0] & 0xF0) >> 4)
	mantissa := int(((d[0] & 0x0F) << 4) | (d[1] & 0x0F))
	return math.Pow(float64(2), float64(exponent)) * float64(mantissa) * 0.045
}

func (s *MAX44009) clean() {
	s.i2c = nil
}
