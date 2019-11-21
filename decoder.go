// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayenne

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var ErrInvalidChannel = errors.New("cayennelpp: unknown type")

type UplinkMessage struct {
	values map[string]interface{}
}

type DownlinkMessage struct {
	values map[uint8]interface{}
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (msg *UplinkMessage) Values() map[string]interface{} {
	return msg.values
}

func (d *Decoder) DecodeUplink() (*UplinkMessage, error) {
	target := &UplinkMessage{
		values: make(map[string]interface{}),
	}
	buf := make([]byte, 2)
	for {
		_, err := io.ReadFull(d.r, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch buf[1] {
		case DigitalInput:
			err = d.decodeDigitalInput(buf[0], target)
		case DigitalOutput:
			err = d.decodeDigitalOutput(buf[0], target)
		case AnalogInput:
			err = d.decodeAnalogInput(buf[0], target)
		case AnalogOutput:
			err = d.decodeAnalogOutput(buf[0], target)
		case Luminosity:
			err = d.decodeLuminosity(buf[0], target)
		case Presence:
			err = d.decodePresence(buf[0], target)
		case Temperature:
			err = d.decodeTemperature(buf[0], target)
		case RelativeHumidity:
			err = d.decodeRelativeHumidity(buf[0], target)
		case Accelerometer:
			err = d.decodeAccelerometer(buf[0], target)
		case BarometricPressure:
			err = d.decodeBarometricPressure(buf[0], target)
		case Gyrometer:
			err = d.decodeGyrometer(buf[0], target)
		case GPS:
			err = d.decodeGPS(buf[0], target)
		default:
			err = ErrInvalidChannel
		}
		if err != nil {
			return nil, err
		}
	}
	return target, nil
}

func (d *Decoder) DecodeDownlink() (*DownlinkMessage, error) {
	target := &DownlinkMessage{make(map[uint8]interface{})}
	buf := make([]byte, 1)
	for {
		_, err := io.ReadFull(d.r, buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if buf[0] == 0xFF {
			break
		}
		var val int16
		if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		target.values[buf[0]] = float32(val) / 100
	}
	return target, nil
}

func (d *Decoder) decodeDigitalInput(channel uint8, target *UplinkMessage) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("digital_input_%d", channel)] = val
	return nil
}

func (d *Decoder) decodeDigitalOutput(channel uint8, target *UplinkMessage) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("digital_output_%d", channel)] = val
	return nil
}

func (d *Decoder) decodeAnalogInput(channel uint8, target *UplinkMessage) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("analog_input_%d", channel)] = float32(val) / 100
	return nil
}

func (d *Decoder) decodeAnalogOutput(channel uint8, target *UplinkMessage) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}

	target.values[fmt.Sprintf("analog_output_%d", channel)] = float32(val) / 100
	return nil
}

func (d *Decoder) decodeLuminosity(channel uint8, target *UplinkMessage) error {
	var val uint16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("luminosity_%d", channel)] = val
	return nil
}

func (d *Decoder) decodePresence(channel uint8, target *UplinkMessage) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("presence_%d", channel)] = val
	return nil
}

func (d *Decoder) decodeTemperature(channel uint8, target *UplinkMessage) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("temperature_%d", channel)] = float32(val) / 10
	return nil
}

func (d *Decoder) decodeRelativeHumidity(channel uint8, target *UplinkMessage) error {
	var val uint8
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("relative_humidity_%d", channel)] = float32(val) / 2
	return nil
}

func (d *Decoder) decodeAccelerometer(channel uint8, target *UplinkMessage) error {
	var valX, valY, valZ int16
	if err := binary.Read(d.r, binary.BigEndian, &valX); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valY); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valZ); err != nil {
		return err
	}
	target.values[fmt.Sprintf("accelerometer_%d", channel)] = []float32{float32(valX) / 1000, float32(valY) / 1000, float32(valZ) / 1000}
	return nil
}

func (d *Decoder) decodeBarometricPressure(channel uint8, target *UplinkMessage) error {
	var val int16
	if err := binary.Read(d.r, binary.BigEndian, &val); err != nil {
		return err
	}
	target.values[fmt.Sprintf("barometric_pressure_%d", channel)] = float32(val) / 10
	return nil
}

func (d *Decoder) decodeGyrometer(channel uint8, target *UplinkMessage) error {
	var valX, valY, valZ int16
	if err := binary.Read(d.r, binary.BigEndian, &valX); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valY); err != nil {
		return err
	}
	if err := binary.Read(d.r, binary.BigEndian, &valZ); err != nil {
		return err
	}
	target.values[fmt.Sprintf("gyrometer_%d", channel)] = []float32{float32(valX) / 100, float32(valY) / 100, float32(valZ) / 100}

	return nil
}

func (d *Decoder) decodeGPS(channel uint8, target *UplinkMessage) error {
	buf := make([]byte, 9)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return err
	}
	latitude := make([]byte, 4)
	copy(latitude, buf[0:3])
	longitude := make([]byte, 4)
	copy(longitude, buf[3:6])
	altitude := make([]byte, 4)
	copy(altitude, buf[6:9])
	target.values[fmt.Sprintf("gps_%d", channel)] = []float32{
		float32(int32(binary.BigEndian.Uint32(latitude))>>8) / 10000,
		float32(int32(binary.BigEndian.Uint32(longitude))>>8) / 10000,
		float32(int32(binary.BigEndian.Uint32(altitude))>>8) / 100,
	}

	return nil
}
