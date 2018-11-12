go-envirophat
==

`go-envirophat` is a port of the python library [enviro-phat](https://github.com/pimoroni/enviro-phat) with some design liberties taken. It allows you to interact with the [Enviro pHAT](https://shop.pimoroni.com/products/enviro-phat) on the Rasberry Pi.

### Leds
```
import "github.com/zaquestion/go-envirophat/leds"
```
```
leds.On()
```

### Light
```
import "github.com/zaquestion/go-envirophat/light"
```
```
	i2c, err := light.InitI2C()
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	data, err := light.Read()
	if err != nil {
		log.Fatal(err)
	}
	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
```
yields
```
{
  "Red": 161,
  "Green": 108,
  "Blue": 98,
  "RawLux": 1183,
  "RawRed": 749,
  "RawGreen": 505,
  "RawBlue": 456,
  "NormalRed": 0.6331360946745562,
  "NormalGreen": 0.4268808114961961,
  "NormalBlue": 0.38546069315300086
}
```

### Weather
```
	i2c, err := weather.InitI2C()
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	// Read temperature in celsius degree
	t, err := weather.Temperature()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Temperature = %v*C\n", t)
	// Read atmospheric pressure in pascal
	p, err := weather.Pressure()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Pressure = %v Pa\n", p)
	// Read atmospheric altitude in meters above sea level, if we assume
	// that pressure at see level is equal to 101325 Pa.
	a, err := weather.Altitude()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Altitude = %v m\n", a)
```

### TODO:
 - Remaining sensors
 - Support raspberry pi 1?
