package main

import (
	"time"

	powerup "github.com/hybridgroup/tinygo-powerup"
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	device  bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)

	airplane *powerup.Airplane
)

func main() {
	wait()

	println("enabling...")
	must("enable BLE interface", adapter.Enable())

	println("start scan...")
	must("start scan", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		println("connected to ", result.Address.String())
	}

	defer device.Disconnect()

	airplane = powerup.NewAirplane(&device)
	err = airplane.Start()
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("motor")
	err = airplane.Throttle(25)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("rudder 0")
	err = airplane.Rudder(0)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("rudder -45")
	err = airplane.Rudder(-45)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("rudder 45")
	err = airplane.Rudder(45)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("rudder 0")
	err = airplane.Rudder(0)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("motor 50")
	err = airplane.Throttle(50)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("motor 75")
	err = airplane.Throttle(75)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("motor 100")
	err = airplane.Throttle(100)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("motor 125")
	err = airplane.Throttle(125)
	if err != nil {
		println(err)
	}

	time.Sleep(3 * time.Second)

	println("stopping")
	airplane.Stop()
}

func scanHandler(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	println("device:", d.Address.String(), d.RSSI, d.LocalName())
	if d.Address.String() == connectAddress() {
		a.StopScan()
		ch <- d
	}
}

func must(action string, err error) {
	if err != nil {
		for {
			println("failed to " + action + ": " + err.Error())
			time.Sleep(time.Second)
		}
	}
}
