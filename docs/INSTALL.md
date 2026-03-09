# ScreamLock — Installation for Parents

## What You Need

- A Windows PC (Windows 10 or 11).
- Two files: `screamlock.exe` and `screamlock-config.exe` (from the [Releases](https://github.com/Maximvs-z/ScreamLock/releases) page or the `build` folder if you built it yourself).

## Step 1: Place the Programs

1. Create a folder that only you (the parent) use, for example:  
   `C:\Programs\ScreamLock`
2. Copy **screamlock.exe** and **screamlock-config.exe** into that folder.

Do **not** put them on the Desktop or in a place the child can easily find or delete.

## Step 2: Choose the Microphone (First-Time Setup)

1. Double-click **screamlock-config.exe** (ScreamLock Config). A small window opens.
2. In **Microphone**, choose the device to monitor (e.g. headset mic) or leave **“(Default microphone)”**.
3. Adjust **Sensitivity (dB)** if needed: more negative = less sensitive (e.g. `-60`); less negative = more sensitive (e.g. `-40`). Default `-50` is a good start.
4. Set **Check every (seconds)** (default 1).
5. Click **Save**. The window closes and the settings are stored for the main app.

*(Advanced: you can instead run `screamlock.exe -list-devices` and edit `config.json` in `%APPDATA%\ScreamLock`.)*

## Step 3: Run at Startup

So that ScreamLock starts when you log on and keeps working after restarts:

**Easiest:** In **ScreamLock Config**, click **"Run at Windows startup"**. This creates a scheduled task for you. Done.

**Manual (optional):** If you prefer to set it up yourself, press **Win + R**, type `taskschd.msc`, press Enter, then **Create Basic Task** → Name: `ScreamLock` → Trigger: **When I log on** → Action: **Start a program** → Program: `C:\Programs\ScreamLock\screamlock.exe` (use your folder) → Finish. In the task list, right-click **ScreamLock** → **Properties** → **Settings** → uncheck **Stop the task if it runs longer than**.

After the next logon, ScreamLock will start automatically. There is no window; it runs in the background.

## Step 4: Verify It's Running

- After logging in, check the log file:  
  `%APPDATA%\ScreamLock\screamlock.log`  
  You should see a line like:  
  `ScreamLock monitoring (threshold -50.00 dB ...)`  
- Optionally, speak loudly into the selected microphone; the PC should lock when the level exceeds the threshold.

## Stopping or Changing Settings

- **Stop:** In Task Scheduler, right-click the ScreamLock task → **Disable** (or **End** if it’s running).
- **Change microphone or sensitivity:**  
  Run **screamlock-config.exe** again, choose the new options, and click Save. Then restart the task or log off and log on so ScreamLock picks up the new config.

## Troubleshooting

See the main [README](../README.md#troubleshooting) for log location, “Open microphone” errors, and device changes.
