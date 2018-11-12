// light is a port of https://github.com/pimoroni/enviro-phat/blob/master/library/envirophat/tcs3472.py
package light

import (
	"math"

	"github.com/d2r2/go-i2c"
	"github.com/pkg/errors"

	"github.com/zaquestion/go-envirophat"
)

const (
	addr = 0x29
)

var (
	i2cClient *i2c.I2C
	maxCount  int
)

// InitI2C creates a new i2c client and connects to the light sensor. The
// client is then used internally by the package. The caller is responible for
// closing the connection.
func InitI2C() (envirophat.Closer, error) {
	//logger.ChangePackageLogLevel("i2c", logger.PanicLevel)

	var err error
	// TODO: figure out how to detect raspberry pi version to properly select the bus
	// Ref: https://github.com/pimoroni/enviro-phat/blob/master/library/envirophat/i2c_bus.py#L19
	i2cClient, err = i2c.NewI2C(addr, 1)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to i2c bus")
	}
	err = i2cClient.WriteRegU8(RegEnable, RegEnableRGBC|RegEnablePower)
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable light sensor")
	}
	err = SetIntegrationTimeMS(511.2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set integration time")
	}

	return i2cClient, nil
}

type Register = uint8

const (
	RegLux   Register = 0xb4
	RegRed   Register = 0xb6
	RegGreen Register = 0xb8
	RegBlue  Register = 0xba

	RegATime  Register = 0x81
	RegStatus Register = 0x93

	RegControl        Register = 0x8f
	RegControlGain1x  Register = 0x00
	RegControlGain4x  Register = 0x01
	RegControlGain16x Register = 0x02
	RegControlGain60x Register = 0x03

	RegEnable          Register = 0x80
	RegEnableInterrupt Register = 0x10
	RegEnableWait      Register = 0x08
	RegEnableRGBC      Register = 0x02
	RegEnablePower     Register = 0x01
)

// Reading provides a complete look at the sensor data with a varieting of formats
// Raw* values is the raw uint16 data pulled from the sensor
// Normal* values are normalized to 0.0 - 1.0 against RawLux
// Red, Green, and Blue range from 0 - 255
type Reading struct {
	Red   uint8
	Green uint8
	Blue  uint8

	RawLux   uint16
	RawRed   uint16
	RawGreen uint16
	RawBlue  uint16

	NormalRed   float64
	NormalGreen float64
	NormalBlue  float64
}

func checkI2C() error {
	if i2cClient == nil {
		return errors.New("i2cClient is nil, must call InitI2C()")
	}
	return nil
}

// SetIntegrationTimeMS sets the sensor integration time in milliseconds.
// ms is the integration time in milliseconds from 2.4 to 612, in increments of 2.4.
func SetIntegrationTimeMS(ms float64) error {
	if err := checkI2C(); err != nil {
		return err
	}
	if ms < 2.4 || ms > 612.0 {
		return errors.New("Integration time must be between 2.4 and 612ms")
	}
	atime := math.Round(ms / 2.4)
	// update the maxCount for use in MaxCount()
	maxCount = int(math.Min(65535.0, (256.0-atime)*1024.0))
	atime -= 1
	err := i2cClient.WriteRegU8(RegATime, 255-uint8(atime))
	if err != nil {
		return err
	}
	return nil
}

// MaxCount returns the maximum value which can be counted by a channel with
// the chosen integration time
func MaxCount() int {
	return maxCount
}

// Valid reads the status on the light sensor and returns true if the device
// checks out
func Valid() bool {
	if err := checkI2C(); err != nil {
		return false
	}
	status, err := i2cClient.ReadRegU8(RegStatus)
	if err != nil {
		return false
	}
	return status&1 > 0
}

// Read gets the raw lux and rgb values from the sensors and returns a fully
// populated Reading which houses the data in a variety of formats. This is the
// biggest difference from the python library and was done so that all formats
// with minimal reads from the sensor.
func Read() (Reading, error) {
	if err := checkI2C(); err != nil {
		return Reading{}, err
	}
	l, err := i2cClient.ReadRegU16LE(RegLux)
	if err != nil {
		return Reading{}, errors.Wrap(err, "failed reading lux")
	}
	r, err := i2cClient.ReadRegU16LE(RegRed)
	if err != nil {
		return Reading{}, errors.Wrap(err, "failed reading red")
	}
	g, err := i2cClient.ReadRegU16LE(RegGreen)
	if err != nil {
		return Reading{}, errors.Wrap(err, "failed reading green")
	}
	b, err := i2cClient.ReadRegU16LE(RegBlue)
	if err != nil {
		return Reading{}, errors.Wrap(err, "failed reading blue")
	}

	reading := Reading{
		RawLux:   l,
		RawRed:   r,
		RawGreen: g,
		RawBlue:  b,

		NormalRed:   float64(r) / float64(l),
		NormalGreen: float64(g) / float64(l),
		NormalBlue:  float64(b) / float64(l),
	}
	reading.Red = uint8(reading.NormalRed * 255)
	reading.Green = uint8(reading.NormalGreen * 255)
	reading.Blue = uint8(reading.NormalBlue * 255)
	return reading, nil
}
