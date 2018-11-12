package weather

import (
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
	"github.com/pkg/errors"

	"github.com/zaquestion/go-envirophat"
)

const (
	addr = 0x77
)

var (
	sensor *bsbmp.BMP
)

// InitI2C creates a new i2c client and connects to the weather sensor. The
// client is then used internally by the package. The caller is responible for
// closing the connection.
func InitI2C() (envirophat.Closer, error) {
	logger.ChangePackageLogLevel("i2c", logger.PanicLevel)
	logger.ChangePackageLogLevel("bsbmp", logger.PanicLevel)

	// TODO: figure out how to detect raspberry pi version to properly select the bus
	// Ref: https://github.com/pimoroni/enviro-phat/blob/master/library/envirophat/i2c_bus.py#L19
	i2c, err := i2c.NewI2C(addr, 1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to i2c bus")
	}

	sensor, err = bsbmp.NewBMP(bsbmp.BMP280, i2c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bmp280 sensor")
	}

	return i2c, nil
}

// Temperature is returned in Celcius and may be affected by
// nearby sources of heat, including the Pi itself.
func Temperature() (float32, error) {
	t, err := sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return 0, err
	}
	return t, nil
}

// Pressure is returned in Pa units
func Pressure() (float32, error) {
	p, err := sensor.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return 0, err
	}
	return p, nil
}

// Altitude returns the currently alltitude
func Altitude() (float32, error) {
	a, err := sensor.ReadAltitude(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return 0, err
	}
	return a, nil
}
