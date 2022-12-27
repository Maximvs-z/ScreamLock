# ScreamLock

This is a PowerShell script designed to monitor the microphone input level and lock the Windows session if the level exceeds a specified threshold. It is intended to be run continuously in the background. 

*While it may potentially be used to address certain behavioral issues related to loud noises while playing games, ie. screaming, it is important to note that I am not a psychologist and cannot guarantee the effectiveness or lack of harm of this script. I cannot be held responsible for any consequences resulting from the use of this script. Use at your own risk.*

<br>
<br>

## Running the PowerShell Script at Startup Using the Task Scheduler

To run this PowerShell script when a local user logs in to Windows, you can use the Task Scheduler to create a task that runs at startup. Here is how you can do this:

Open the Task Scheduler.<br><br>
In the Actions pane, click "Create Basic Task".<br><br>
Follow the prompts to specify a name and description for the task, and then click "Next".<br><br>
Select the "When I log on" trigger, and then click "Next".<br><br>
Select the "Start a program" action, and then click "Next".<br><br>
In the "Program/script" field, enter the path to the PowerShell executable (e.g., "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe").<br><br>
In the "Add arguments (optional)" field, enter the path to your PowerShell script (e.g., "C:\Scripts\MyScript.ps1").<br><br>
Click "Finish" to create the task.<br><br>
This will create a task that runs the specified PowerShell script when the local user logs in to Windows. Note that the script will only run for the user who creates the task, and not for other users on the system.

# Tweaking
The most important value to customize is the $threshold variable, which is currently set to -50. You may need to adjust this value based on your specific requirements."
