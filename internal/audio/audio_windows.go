//go:build windows
// +build windows

// Package audio uses Windows Core Audio (WASAPI) via go-wca to read microphone peak level.
// Peak is returned as linear 0..1; we convert threshold from dB for comparison.
package audio

import (
	"math"

	"github.com/moutend/go-wca/pkg/wca"
)

// ListCaptureDevices returns all active capture devices. Caller must have called ole.CoInitializeEx before.
func ListCaptureDevices() ([]CaptureDevice, error) {
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil, err
	}
	defer mmde.Release()

	var dc *wca.IMMDeviceCollection
	if err := mmde.EnumAudioEndpoints(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &dc); err != nil {
		return nil, err
	}
	defer dc.Release()

	var count uint32
	if err := dc.GetCount(&count); err != nil {
		return nil, err
	}

	var list []CaptureDevice
	for i := uint32(0); i < count; i++ {
		var mmd *wca.IMMDevice
		if err := dc.Item(i, &mmd); err != nil {
			continue
		}
		id, name, _ := deviceIDAndName(mmd)
		mmd.Release()
		if id != "" {
			list = append(list, CaptureDevice{ID: id, Name: name})
		}
	}
	return list, nil
}

func deviceIDAndName(mmd *wca.IMMDevice) (id string, name string, err error) {
	if err = mmd.GetId(&id); err != nil {
		return "", "", err
	}
	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return id, id, nil
	}
	defer ps.Release()
	// PKEY_Device_FriendlyName - we can try to get it; if not available use ID
	name = id
	// Simplified: many examples read the property store for friendly name. Skip for brevity; name can be ID.
	return id, name, nil
}

// PeakReader reads the current peak level (0.0–1.0) from a capture device.
// Call Close when done.
type PeakReader struct {
	ac  *wca.IAudioClient
	ami *wca.IAudioMeterInformation
}

// OpenPeakReader opens the default capture device or the device with the given ID.
// deviceID empty = default. Caller must have called ole.CoInitializeEx before.
func OpenPeakReader(deviceID string) (*PeakReader, error) {
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil, err
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if deviceID == "" {
		if err := mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmd); err != nil {
			return nil, err
		}
	} else {
		var dc *wca.IMMDeviceCollection
		if err := mmde.EnumAudioEndpoints(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &dc); err != nil {
			return nil, err
		}
		var c uint32
		dc.GetCount(&c)
		for i := uint32(0); i < c; i++ {
			var d *wca.IMMDevice
			if err := dc.Item(i, &d); err != nil {
				continue
			}
			var id string
			d.GetId(&id)
			if id == deviceID {
				mmd = d
				dc.Release()
				break
			}
			d.Release()
		}
		if mmd == nil {
			dc.Release()
			// Device not found (e.g. unplugged); fall back to default
			if err := mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmd); err != nil {
				return nil, err
			}
		}
	}
	defer mmd.Release()

	var ac *wca.IAudioClient
	if err := mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &ac); err != nil {
		return nil, err
	}

	var wfx *wca.WAVEFORMATEX
	if err := ac.GetMixFormat(&wfx); err != nil {
		ac.Release()
		return nil, err
	}

	// 100 ms buffer
	const bufRef = 100 * 10000 // REFERENCE_TIME is 100-nanosecond units
	if err := ac.Initialize(wca.AUDCLNT_SHAREMODE_SHARED, 0, bufRef, 0, wfx, nil); err != nil {
		ac.Release()
		return nil, err
	}

	var ami *wca.IAudioMeterInformation
	if err := ac.GetService(wca.IID_IAudioMeterInformation, &ami); err != nil {
		ac.Release()
		return nil, err
	}

	if err := ac.Start(); err != nil {
		ami.Release()
		ac.Release()
		return nil, err
	}

	return &PeakReader{ac: ac, ami: ami}, nil
}

// Peak returns the current peak level (0.0 to 1.0). Linear scale.
func (p *PeakReader) Peak() (float32, error) {
	var peak float32
	err := p.ami.GetPeakValue(&peak)
	return peak, err
}

// Close stops the stream and releases resources.
func (p *PeakReader) Close() error {
	if p.ac != nil {
		_ = p.ac.Stop()
		p.ac.Release()
		p.ac = nil
	}
	if p.ami != nil {
		p.ami.Release()
		p.ami = nil
	}
	return nil
}

// LinearFromDB converts dB (e.g. -50) to linear scale for comparison with GetPeakValue.
// Formula: 10^(dB/20). For -50 dB -> ~0.00316.
func LinearFromDB(dB float64) float32 {
	if dB <= -100 {
		return 0
	}
	linear := math.Pow(10, dB/20)
	if linear < 0 {
		return 0
	}
	if linear > 1 {
		return 1
	}
	return float32(linear)
}

// DBFromLinear converts linear peak (0..1) to dB for display. Useful for level meters.
func DBFromLinear(linear float32) float64 {
	if linear <= 0 {
		return -100
	}
	return 20 * math.Log10(float64(linear))
}
