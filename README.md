Welcome back, Commander 🫡.

Red Alert 3 is one of my favorite game but it's in a sorry state since EA shut down their servers for the game.

This project is an open source re-implementation of the original game launcher `RA3.exe` aimed at addressing and mitigating some of the issues that have emerged since EA servers were shut down.

### Issues due to EA shutting down their servers

1. Game Startup is (very) slow

   The game takes awfully long to start because the official launcher tries to get the latest "Comrade News" from files.ea.com
   whether or not the gui interface was requested (-ui). And wait until the connection times out, before starting the game.

2. LAN play: "CD key is already in use" (%CDKEY%)

   LAN play requires each player to have a different CD Key.
   Steam no longer populates the `%CDKEY%` variable in the registry _(is it due to EA's servers shutdown ?)_.
  
   This is a problem because each time the game is run by Steam; Steam will write "%CDKEY%" as the user's CD key.
   Every (Steam) players will therefore end up with the same CD Key: "%CDKEY%" and are unable to LAN play with each others.

   The registry value is read on game launch. It is uneffective to change it after the game is ran by Steam.

3. Online play: emulation and DLL injection

   You can no longer play online nor co-op without using 3rd party service that emulate gamespy such as [revora/CnC:Online](https://cnc-online.net/en/)

   Requiring an additional launcher to be able to use their service.
 
   While _"Launchers inception"_ (Launcher which starts another Launcher) is despised by many it remains a matter of personal preference. 
   But for me the core issue was that their launcher didn't work with the Steam version nor with Linux/Proton when I tried it.

## Features

- Start Red Alert 3 process almost instantly _(Issue 1)_

- Fix the CD Key registry value if needed for LAN play before starting the game _(Issue 2)_

- Addons: revora/CnC:Online without their launcher _(Issue 3)_

  See [xan105/CnC-Online](https://github.com/xan105/CnC-Online) for more details.
  
- Compatible with 🐧 Linux/Proton

- Splash screen customisation: none, random, ...

## Installation

Copy the files in the game directory and replace `RA3.exe`

#### Steam

> [!NOTE]
> In order to fix the issue with `%CDKEY%` _(Issue 2)_ `RA3.exe` will need elevated privileges.<br/>
> Because the CD Key value is in the registry under `HKLM` which is write access protected.
>
> If you run the game via Steam, you may want to consider always running `RA3.exe` with elevated privileges<br/>
> by right clicking RA3.exe > properties > Compatibility > check "Run this program as an administrator"

<details><summary>Wait... why does it need admin right and Steam does not ?!</summary>
<br/>
Steam has its own windows service running in the background with system privileges (`steamservice.exe`) to do these kind of operations silently without the end user noticing.
</details>

## Command Line Arguments

`RA3.exe --help` to display all arguments.

Most of the orignal `RA3.exe` cmdline arguments were kept:

|flag|type|description|
|----|----|-----------|
|xres|number|Sets resolution width|
|yres|number|Sets resolution height|
|xpos|number|Sets horizontal offset of the window|
|ypos|number|Sets vertical offset of the window|
|win|boolean|Runs the game in windowed mode|
|fullscreen|boolean|Runs the game in fullscreen mode. Combine with -win for borderless windowed mode|
|noaudio|boolean|Disables game audio|
|noAudioMusic|boolean|Disables game music|
|silentLogin|boolean|Forces the game to immediately log in to a multiplayer account|
|modConfig|string|Runs the game with selected mod (has to point to its .skudef file)|
|replayGame|string|Plays replay file|

> [!TIP]
> `-modConfig path`
> 
> If "path" is **not** an absolute path then it will look for any `.skudef` file corresponding in `%Documents%/Red Alert 3/Mods`
>
> eg: `%Documents%/Red Alert 3/Mods/Upheaval/Upheaval_1.16.skudef` > `-modConfig Upheaval`

## Config File

📄 `RA3.json` (required): 

- `lang?: string (auto)`

  Which language to start Red Alert 3 with.

  Default (auto) will query the registry value `HKCU/Software/Electronic Arts/Electronic Arts/Red Alert 3/language`.<br/>
  If no value is set then it defaults to "english".

- `borderless?: boolean (false)`

  Runs Red Alert 3 in borderless fullscreen.<br/>
  This superseeds the cmdline arguments `-win` and `-fullscreen`
 
- `upheaval?: boolean (false)`

  Starts Red Alert 3 with the infamous mod "upheaval" if it's present in the `%GAMEDIR%` or in `%Documents%/Red Alert 3/Mods/Upheaval`

- `keygen?: boolean (true)`

  Check the CD Key value in the registry.<br/>
  If it's empty or equals to "%CDKEY% generates a random CD Key and write it to the registry.

> [!CAUTION]
> Unfortunately requires elevated privileges because the key is located under `HKLM` which is write access protected.

- `splash?: boolean (true)`

  Display a splash screen while the game is loading similar to the original `RA3.exe`

- `splash_image?: []string (["Launcher/splash.bmp"])`

  Splash screen filepath. When more than one, a splash screen is selected at random.

  Either absolute or relative path.<br/>
  _NB: relative to `RA3.exe` and **not** the current working dir_

> [!NOTE]
> Image should be a 640x480 BMP file.

- `addons?: []{ path: string, required?: boolean }`

  List of addons to inject to the game process.<br/>
  When `required` is set to `true` and if the injection failed, alert the user and kill the game process.

  Either absolute or relative path.<br/>
  _NB: relative to `RA3.exe` and **not** the current working dir_

> [!TIP]
> Example: You can use this option to load [xan105/CnC-Online](https://github.com/xan105/CnC-Online).
>
> Restoring the online features of Red Alert 3 without relying on the revora/CnC:Online launcher.

  ```json
  {
    "addons": [
      { "path": "Launcher/opencnconline.dll", "required": true }
    ]
  }
  ```

  <p align="center">
    <img src="https://github.com/xan105/RA3-Launcher/raw/main/screenshot/linux_proton.png">
    <em>Connected to C&C:Online under 🐧 Linux/Proton 9.0-2 (Fedora)</em>
  </p>

## Building

- Golang v1.23
- [go-winres](https://github.com/tc-hib/go-winres) installed in `%PATH%` env var for win32 manifest & cie

Run `build.cmd`<br/>
Output files are located in `./build`
