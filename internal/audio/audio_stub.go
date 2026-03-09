//go:build !windows
// +build !windows

package audio

// PeakReader stub for non-Windows.
type PeakReader struct{}

// ListCaptureDevices is only implemented on Windows.
func ListCaptureDevices() ([]CaptureDevice, error) {
	return nil, nil
}

// OpenPeakReader is only implemented on Windows.
func OpenPeakReader(deviceID string) (*PeakReader, error) {
	return nil, nil
}

// Peak stub.
func (p *PeakReader) Peak() (float32, error) {
	return 0, nil
}

// Close stub.
func (p *PeakReader) Close() error {
	return nil
}

// LinearFromDB stub.
func LinearFromDB(dB float64) float32 {
	return 0
}

// DBFromLinear stub.
func DBFromLinear(linear float32) float64 {
	return -100
}
