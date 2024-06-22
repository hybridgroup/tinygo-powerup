package powerup

import (
	"errors"

	"tinygo.org/x/bluetooth"
)

type Airplane struct {
	device  *bluetooth.Device
	control *bluetooth.DeviceService
	motor   *bluetooth.DeviceCharacteristic
	rudder  *bluetooth.DeviceCharacteristic

	buf []byte
}

var (
	// BLE services
	// 75b64e51f1814ed1921a476090d80ba7
	powerupService       = bluetooth.NewUUID([16]byte{0x75, 0xb6, 0x4e, 0x51, 0xf1, 0x81, 0x4e, 0xd1, 0x92, 0x1a, 0x47, 0x60, 0x90, 0xd8, 0x0b, 0xa7})
	motorCharacteristic  = bluetooth.NewUUID([16]byte{0x75, 0xb6, 0x4e, 0x51, 0xf1, 0x84, 0x4e, 0xd1, 0x92, 0x1a, 0x47, 0x60, 0x90, 0xd8, 0x0b, 0xa7})
	rudderCharacteristic = bluetooth.NewUUID([16]byte{0x75, 0xb6, 0x4e, 0x51, 0xf1, 0x85, 0x4e, 0xd1, 0x92, 0x1a, 0x47, 0x60, 0x90, 0xd8, 0x0b, 0xa7})
)

// NewAirplane creates a new Powerup Airplane.
func NewAirplane(dev *bluetooth.Device) *Airplane {
	a := &Airplane{
		device: dev,
		buf:    make([]byte, 255),
	}

	return a
}

func (a *Airplane) Start() (err error) {
	srvcs, err := a.device.DiscoverServices([]bluetooth.UUID{
		powerupService,
	})
	if err != nil || len(srvcs) == 0 {
		return errors.New("could not find services")
	}

	a.control = &srvcs[0]
	println("found powerup control service", a.control.UUID().String())

	chars, err := a.control.DiscoverCharacteristics([]bluetooth.UUID{
		motorCharacteristic,
		rudderCharacteristic,
	})
	if err != nil || len(chars) == 0 {
		return errors.New("could not find powerup control characteristic")
	}

	a.motor = &chars[0]
	a.rudder = &chars[1]

	return
}

// Stops stops the Airplane.
func (a *Airplane) Stop() (err error) {
	a.Throttle(0)
	a.Rudder(0)

	return
}

// Throttle sets the throttle of the Airplane.
func (a *Airplane) Throttle(thrust int) (err error) {
	buf := []byte{uint8(thrust)}
	_, err = a.motor.WriteWithoutResponse(buf)

	return err
}

// Rudder sets the rudder of the Airplane.
// angle goes from -45 to 45.
func (a *Airplane) Rudder(angle int) (err error) {
	buf := []byte{byte(angle)}
	_, err = a.rudder.WriteWithoutResponse(buf)

	return err
}
