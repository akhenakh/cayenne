// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package cayenne

import (
	"bytes"
	"testing"

	"io"

	. "github.com/smartystreets/assertions"
)

func TestDecode(t *testing.T) {
	a := New(t)

	// Happy flow: uplink
	{
		buf := []byte{
			1, DigitalInput, 255,
			2, DigitalOutput, 100,
			3, AnalogInput, 21, 74,
			4, AnalogOutput, 234, 182,
			5, Luminosity, 1, 244,
			6, Presence, 50,
			7, Temperature, 255, 100,
			8, RelativeHumidity, 160,
			9, Accelerometer, 254, 88, 0, 15, 6, 130,
			10, BarometricPressure, 41, 239,
			11, Gyrometer, 1, 99, 2, 49, 254, 102,
			12, GPS, 7, 253, 135, 0, 190, 245, 0, 8, 106,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))

		msg, err := decoder.DecodeUplink()
		a.So(err, ShouldBeNil)
		a.So(msg.Values()["digital_input_1"], ShouldEqual, 255)
		a.So(msg.Values()["digital_output_2"], ShouldEqual, 100)
		a.So(msg.Values()["analog_input_3"], ShouldEqual, 54.5)
		a.So(msg.Values()["analog_output_4"], ShouldEqual, -54.5)
		a.So(msg.Values()["luminosity_5"], ShouldEqual, 500)
		a.So(msg.Values()["presence_6"], ShouldEqual, 50)
		a.So(msg.Values()["temperature_7"], ShouldEqual, -15.6)
		a.So(msg.Values()["relative_humidity_8"], ShouldEqual, 80)
		a.So(msg.Values()["accelerometer_9"], ShouldResemble, []float32{-0.424, 0.015, 1.666})
		a.So(msg.Values()["barometric_pressure_10"], ShouldEqual, 1073.5)
		a.So(msg.Values()["gyrometer_11"], ShouldResemble, []float32{3.55, 5.61, -4.10})
		a.So(msg.Values()["gps_12"], ShouldResemble, []float32{52.3655, 4.8885, 21.54})
		key, ok := msg.GotLocation()
		a.So(ok, ShouldBeTrue)
		a.So(key, ShouldEqual, "gps_12")
	}

	// Happy flow: downlink
	{
		buf := []byte{
			1, 0, 100,
			2, 234, 182,
			255,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))

		target, err := decoder.DecodeDownlink()
		a.So(err, ShouldBeNil)
		a.So(target.values[1], ShouldEqual, 1)
		a.So(target.values[2], ShouldEqual, -54.5)
	}

	// Invalid data type
	{
		buf := []byte{
			1, 255, 255,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))

		_, err := decoder.DecodeUplink()
		a.So(err, ShouldEqual, ErrInvalidChannel)
	}

	// Not enough data: uplink
	{
		buf := []byte{
			12, GPS, 7, 253, 135, 0, 190,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))

		_, err := decoder.DecodeUplink()
		a.So(err, ShouldEqual, io.ErrUnexpectedEOF)
	}

	// Not enough data: downlink
	{
		buf := []byte{
			12, 1,
		}
		decoder := NewDecoder(bytes.NewBuffer(buf))

		_, err := decoder.DecodeDownlink()
		a.So(err, ShouldEqual, io.ErrUnexpectedEOF)
	}

	// Negative coordinates
	{
		buf := []byte{0x01, GPS, 0x06, 0x76, 0x5f, 0xf2, 0x96, 0x0a, 0x00, 0x03, 0xe8}
		decoder := NewDecoder(bytes.NewBuffer(buf))

		target, err := decoder.DecodeUplink()
		a.So(err, ShouldBeNil)
		a.So(target.values["gps_1"], ShouldResemble, []float32{42.3519, -87.9094, 10})
	}
}
