# Sets the threshold for the microphone input level (in dB). Adjust it to your specific requirements.
$threshold = -50

# Get the default audio input device
$inputDevice = Get-WmiObject -Class Win32_SoundDevice | Where-Object {$_.ProductName -match "microphone"}

# Continuously check the input level
while ($true) {

  # Get the current input level
  $inputLevel = $inputDevice.AudioInputMixer.InputGain

  # Check if the input level is above the threshold
  if ($inputLevel -gt $threshold) {
    # Lock the Windows session
    rundll32.exe user32.dll,LockWorkStation
  }

  }

  # Wait for 1 second before checking again
  Start-Sleep -Seconds 1

  # Force garbage collection to free up memory
  [System.GC]::Collect()
}
