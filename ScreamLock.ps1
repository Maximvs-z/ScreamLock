# Import the System.Windows.Forms namespace
Add-Type -AssemblyName System.Windows.Forms

# Sets the threshold level (in dB). Adjust it to your specific needs.
$threshold = -50

# Continuously check the audio level
while (1 -eq 1) {
  # Get the current audio level
  $audio_level = Get-WmiObject -Class Win32_SoundDevice | Select-Object -ExpandProperty CurrentSamplePeak

  # Check if the audio level is above the threshold
  if ($audio_level -gt $threshold) {
    # Lock the Windows session
    [System.Windows.Forms.SendKeys]::SendWait("{LWIN}+{L}")
  }

  # Wait for 1 second before checking again
  Start-Sleep -Seconds 1

  # Force garbage collection to free up memory
  [System.GC]::Collect()
}
