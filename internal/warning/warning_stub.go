//go:build !windows
// +build !windows

package warning

// PlayTone is only implemented on Windows.
func PlayTone() {}

// SpeakAsync is only implemented on Windows.
func SpeakAsync(text string) {}

// RunSequence is only implemented on Windows.
func RunSequence(enableVoice bool, lockFn func()) {
	if lockFn != nil {
		lockFn()
	}
}
