package main

import (
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

const pocpF6MAC = "B0:9C:63:C6:77:5B"

func main() {
	myPhoneMAC, err := bluetooth.ParseMAC(pocpF6MAC)
	if err != nil {
		panic("could not parse MAC address: " + err.Error())
	}

	// Enable BLE interface.
	if err := adapter.Enable(); err != nil {
		panic("could not enable BLE interface: " + err.Error())

	}

	// Start scanning.
	println("scanning...")
	go func() {
		err := adapter.Scan(onScanEvent)
		if err != nil {
			panic("start scan: " + err.Error())
		}
	}()

	// Wait forever.
	btPhone, err := adapter.Connect(bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: myPhoneMAC}}, bluetooth.ConnectionParams{})
	if err != nil {
		panic("could not connect to device: " + err.Error())
	} else {
		println("connected to device", btPhone.Address.String())
	}

}

func onScanEvent(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	if device.Address.String() == pocpF6MAC {
		println("found device:", device.Address.String(), device.RSSI, device.LocalName())
	}
}
