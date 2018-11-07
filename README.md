go-envirophat
==

`go-envirophat` is a port of the python library [enviro-phat](https://github.com/pimoroni/enviro-phat) with some design liberties taken. It allows you to interact with the [Enviro pHAT](https://shop.pimoroni.com/products/enviro-phat) on the Rasberry Pi.

leds
```
import "github.com/zaquestion/go-envirophat/leds"
```
```
leds.On()
```

light
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

TODO:
 - Remaining sensors
 - Support raspberry pi 1?
