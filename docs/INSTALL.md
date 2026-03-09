# ScreamLock — Installation for Parents

## What You Need

- A Windows PC (Windows 10 or 11).
- One file: `screamlock.exe` (from the [Releases](https://github.com/Maximvs-z/ScreamLock/releases) page or the `build` folder if you built it yourself).

## Step 1: Place the Program

1. Create a folder that only you (the parent) use, for example:  
   `C:\Programs\ScreamLock`
2. Copy `screamlock.exe` into that folder.

Do **not** put it on the Desktop or in a place the child can easily find or delete.

## Step 2: Choose the Microphone (First-Time Setup)

1. Open **Command Prompt** or **PowerShell** (right-click → Run as administrator is not required; your user account is enough).
2. Go to the folder where you put the program, for example:
   ```bat
   cd C:\Programs\ScreamLock
   ```
3. Run:
   ```bat
   screamlock.exe -list-devices
   ```
4. A folder will open showing a file named `devices.txt` (and later `config.json`). The folder is usually:
   `%APPDATA%\ScreamLock`  
   (e.g. `C:\Users\YourName\AppData\Roaming\ScreamLock`).
5. Open `devices.txt` in Notepad. You will see a list of microphones with **ID** and **Name**.
6. If you want to use the **default Windows microphone**, you can leave the config as is (see Step 3).  
   If you want a **specific microphone** (e.g. headset mic), copy the full **ID** line (the long line under "ID:") and keep it for Step 3.

## Step 3: Configure the Program

1. In the same folder as `devices.txt`, open or create `config.json`.
2. If the file does not exist, create it with this content (or copy from `config.example.json` in the project):
   ```json
   {
     "device_id": "",
     "threshold_db": -50,
     "check_interval_seconds": 1
   }
   ```
3. Set the microphone:
   - **Default microphone:** leave `"device_id": ""` as is.
   - **Specific microphone:** paste the ID you copied from `devices.txt` between the quotes:
     ```json
     "device_id": "{0.0.1.00000000}.{GUID-here}",
     ```
4. Adjust sensitivity (optional):
   - `threshold_db`: how loud the sound must be to trigger a lock.  
     - More negative = less sensitive (e.g. `-60` = only very loud sounds).  
     - Less negative = more sensitive (e.g. `-40` = quieter sounds can trigger).  
     - Default `-50` is a good starting point.
5. Save `config.json`.

## Step 4: Run at Startup (Task Scheduler)

So that ScreamLock starts when the computer boots and runs in the background:

1. Press **Win + R**, type `taskschd.msc`, press Enter (Task Scheduler).
2. On the right, click **Create Basic Task**.
3. **Name:** e.g. `ScreamLock`  
   **Description:** e.g. `Monitor microphone and lock when too loud`  
   Next.
4. **Trigger:** **When I log on**  
   Next.
5. **Action:** **Start a program**  
   Next.
6. **Program/script:**  
   `C:\Programs\ScreamLock\screamlock.exe`  
   (use the path where you actually put the file.)
7. **Add arguments:** leave empty.
8. **Start in:**  
   `C:\Programs\ScreamLock`  
   (same folder as the exe.)
9. Next → Finish.
10. In the task list, right-click the new **ScreamLock** task → **Properties**.
11. **General** tab:
    - Check **Run with highest privileges** only if you need it (usually not).
    - Select **Run whether user is logged on or not** if you want it to run before login (advanced); for most users, keep **Run only when user is logged on**.
12. **Settings** tab (optional):
    - Uncheck **Stop the task if it runs longer than** (so it can run indefinitely).
13. OK.

After the next logon, ScreamLock will start automatically. There is no window; it runs in the background.

## Step 5: Verify It’s Running

- After logging in, check the log file:  
  `%APPDATA%\ScreamLock\screamlock.log`  
  You should see a line like:  
  `ScreamLock monitoring (threshold -50.00 dB ...)`  
- Optionally, speak loudly into the selected microphone; the PC should lock when the level exceeds the threshold.

## Stopping or Changing Settings

- **Stop:** In Task Scheduler, right-click the ScreamLock task → **Disable** (or **End** if it’s running).
- **Change microphone or sensitivity:**  
  Edit `%APPDATA%\ScreamLock\config.json`, then either restart the task or log off and log on again so ScreamLock restarts with the new config.

## Troubleshooting

See the main [README](../README.md#troubleshooting) for log location, “Open microphone” errors, and device changes.
