package main

import (
	"context"
	"flag"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/go-ble/ble/linux"
	"github.com/go-ble/ble/linux/hci/cmd"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	device = flag.String("device", "default", "implementation of ble hci0-hci1-etc...")
	name   = flag.String("name", "LYWSD03MMC", "name of Xiaomi Sensor")
)

func main() {
	// parse cmd line argument
	flag.Parse()
	// create channel for Ctrl+c or other signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Setup BLE
	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// This part is requested only on certain device
	// where MAC Address is not set by Hardware
	if dev, ok := d.(*linux.Device); ok {
		if err := dev.HCI.Send(&cmd.LESetRandomAddress{
			RandomAddress: [6]byte{0xFF, 0x11, 0x22, 0x33, 0x44, 0x55},
		}, nil); err != nil {
			log.Fatalln(err, "can't set random address")
		}
	}

	// Create Cancellation Context
	ctx := ble.WithSigHandler(context.WithCancel(context.Background()))

	// Default to search device with name of MJ_HT_V1 (or specified by user).
	charaUuid := ble.MustParse("181a")
	filter := func(a ble.Advertisement) bool {
		if len(a.ServiceData()) < 1 {
			return false
		}
		chara := a.ServiceData()[0]
		if !charaUuid.Equal(chara.UUID) {
			return false
		}
		// TODO: filter first 3 bytes of MAC
		return true
	}

	log.Printf("Scanning for devices : %s", *name)
	// start Scannig in another Thread
	go ble.Scan(ctx, true, scanAdvertissement, filter)

	// wait for end of program
	c := <-sigs
	log.Printf("Signal Received - %s - Wait Few Seconds: ", c)
	ctx.Done()
	log.Printf("Done completed")
	time.Sleep(5 * time.Second)
}

// Advertissement frame call baclk
func scanAdvertissement(a ble.Advertisement) {
	//log.Printf("scanAdvertissement called. Adv: %#v", a)
	// check all Service Data for 0xfe95 as describe here
	// https://github.com/hannseman/homebridge-mi-hygrothermograph
	for _, chara := range a.ServiceData() {
		if chara.UUID.Equal(ble.MustParse("ae95")) {
			// if not 18 it's not a comple packet
			if len(chara.Data) != 18 {
				return
			}
			// extract only end of frame
			info := chara.Data[len(chara.Data)-4:]
			// extract temp
			temp := float64(int16(info[1])<<8+int16(info[0])) / 10.0
			// extract humidity
			humidity := float64(int16(info[3])<<8+int16(info[2])) / 10.0
			// log it
			log.Printf("%s - %s [ %v °C - %v %%]\n", a.LocalName(), a.Addr().String(), temp, humidity)
			//log.Printf("%s - %s [ %v °C - %v %%]\n [% X] %d\n", a.LocalName(), a.Addr().String(), temp, humidity, chara.Data, len(chara.Data))
		}
		if chara.UUID.Equal(ble.MustParse("181a")) {

			if len(chara.Data) != 13 {
				log.Printf("Invalid length")
			}
			// Byte 0-5 mac in correct order
			// Byte 6-7 Temperature in int16
			temp := int16(chara.Data[6])*256 + int16(chara.Data[7])
			// Byte 8 Humidity in percent
			humidity := chara.Data[8]
			// Byte 9 Battery in percent
			battery := chara.Data[9]
			// Byte 10-11 Battery in mV uint16_t
			batterymV := int16(chara.Data[10])*256 + int16(chara.Data[11])
			// Byte 12 frame packet counter
			log.Printf("Address= %s, RSSI= %d, Temp= %d, Humidity= %d%%, Battery= %d%%, %d mV\n", a.Addr().String(), a.RSSI(), temp, humidity, battery, batterymV)

		}
	}

}
