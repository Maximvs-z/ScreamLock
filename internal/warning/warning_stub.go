//go:build !windows
// +build !windows

package warning

// PlayTone is only implemented on Windows.
func PlayTone() {}

// SpeakAsync is only implemented on Windows.
func SpeakAsync(text string) {}

// GracePeriodDuration stub.
const GracePeriodDuration = 0

// RunGracePeriod is only implemented on Windows.
func RunGracePeriod(samplePeak func() float32, thresholdLinear float32) bool {
	return true
}

// RunSequence is only implemented on Windows.
func RunSequence(enableVoice bool, lockFn func()) {
	if lockFn != nil {
		lockFn()
	}
}

// RunSequenceRest is only implemented on Windows.
func RunSequenceRest(enableVoice bool, lockFn func()) {
	if lockFn != nil {
		lockFn()
	}
}
