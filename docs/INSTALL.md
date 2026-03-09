# ScreamLock — Installation for Parents

## What You Need

- A Windows PC (Windows 10 or 11).
- **Easiest:** Download **ScreamLock-Setup.exe** from [Releases](https://github.com/Maximvs-z/ScreamLock/releases) (single-file installer).  
- **Or** from Releases or the `build` folder: **screamlock.exe**, **screamlock-config.exe**, and **screamlock-setup.exe**.

## Option A: Single-file installer (recommended)

1. Download **ScreamLock-Setup.exe** from the [Releases](https://github.com/Maximvs-z/ScreamLock/releases) page.
2. Right-click it → **Run as administrator**. It installs to `C:\Program Files\ScreamLock` and adds a task to run ScreamLock when you log on.
3. Open **screamlock-config.exe** from that folder (e.g. via Start menu or `C:\Program Files\ScreamLock`) to choose your microphone and settings. Use the **Test level** bar to test without locking.

## Option B: Manual placement, then wizard

## Step 1: Place the Programs

1. Create a folder that only you (the parent) use, for example:  
   `C:\Programs\ScreamLock`
2. Copy **screamlock.exe**, **screamlock-config.exe**, and **screamlock-setup.exe** into that folder.

Do **not** put them on the Desktop or in a place the child can easily find or delete.

## Step 2: Run the Installer or Configure Manually

**Recommended:** Double-click **screamlock-setup.exe**. Click Next → Next → check "Yes, run ScreamLock when I log on" if you want autostart → Finish. The **Microphone & Test** window opens: choose microphone, set sensitivity, and use the **Test level** bar (no lock — for testing only). Save, then Done.

**Or manual:** Double-click **screamlock-config.exe** (ScreamLock Config). A small window opens.
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

## Step 3: Verify It's Running

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
