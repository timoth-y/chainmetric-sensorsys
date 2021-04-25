package sensors

import (
	"fmt"
	"math"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/model/metrics"
)

type MAX44009 struct {
	*peripherals.I2C
}

func NewMAX44009(addr uint16, bus int) sensor.Sensor {
	return &MAX44009{
		I2C: peripherals.NewI2C(addr, bus),
	}
}

func (s *MAX44009) ID() string {
	return "MAX44009"
}

func (s *MAX44009) Init() error {
	if err := s.I2C.Init(); err != nil {
		return err
	}

	if !s.Verify() {
		return fmt.Errorf("driver is not compatiple with specified sensor")
	}

	if err := s.WriteBytes(MAX44009_APP_START); err != nil {
		return err
	}

	return nil
}

func (s *MAX44009) Read() (float64, error) {
	buffer, err := s.ReadBytes(2); if err != nil {
		return math.NaN(), err
	}

	return s.dataToLuminance(buffer), nil
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
	return true
}

func (s *MAX44009) dataToLuminance(d []byte) float64 {
	exponent := int((d[0] & 0xF0) >> 4)
	mantissa := int(((d[0] & 0x0F) << 4) | (d[1] & 0x0F))
	return math.Pow(float64(2), float64(exponent)) * float64(mantissa) * 0.045
}
