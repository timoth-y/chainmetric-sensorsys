package device

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

const (
	deviceIdentityFile = "device.id"
)

func (d *Device) Init() error {
	var err error

	d.Specs, err = d.DiscoverSpecs(); if err != nil {
		return err
	}

	if id, is := isRegistered(); is {
		if exists, _ := d.client.Contracts.Devices.Exists(id); exists {
			if err := d.client.Contracts.Devices.UpdateSpecs(id, d.Specs); err != nil {
				return err
			}

			shared.Logger.Infof("Device specs been updated in blockchain with id: %s", id)

			return nil
		}

		shared.Logger.Warning("device was removed from network, must re-initialize now")
	}

	qr, err := qrcode.New(d.Specs.Encode(), qrcode.Medium); if err != nil {
		return err
	}

	d.display.PowerOn()

	d.display.DrawImage(qr.Image(d.config.Display.ImageSize))

	shared.Logger.Debug("Subscribing to blockchain...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)

	if err := d.client.Contracts.Devices.Subscribe(ctx, "inserted", func(device models.Device) error {
		if device.Hostname == d.Specs.Hostname {
			defer cancel()

			if err := storeIdentity(device.ID); err != nil {
				shared.Logger.Fatal(errors.Wrap(err, "failed to store device's identity file"))
			}

			shared.Logger.Infof("Device is been registered with id: %s", device.ID)

			return d.client.Contracts.Devices.UpdateSpecs(device.ID, d.Specs)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func isRegistered() (string, bool) {
	id, err := ioutil.ReadFile(deviceIdentityFile); if err != nil {
		if os.IsNotExist(err) {
			return "", false
		}

		shared.Logger.Fatal(errors.Wrap(err, "failed to read device identity file"))
	}

	return string(id), true
}

func storeIdentity(id string) error {
	f, err := os.Create(deviceIdentityFile); if err != nil {
		return err
	}

	if _, err := f.WriteString(id); err != nil {
		return err
	}

	return nil
}
