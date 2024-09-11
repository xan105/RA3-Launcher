/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import (
  "os"
  "golang.org/x/sys/windows"
  "log/slog"
)

func alert(message string){
  slog.Error(message);
  windows.MessageBox(
    windows.HWND(uintptr(0)),
    windows.StringToUTF16Ptr(message),
    windows.StringToUTF16Ptr("Red Alert 3"),
    windows.MB_OK | windows.MB_ICONERROR)
  os.Exit(1)
}

func displayHelp(){
  windows.MessageBox(
      windows.HWND(uintptr(0)),
      windows.StringToUTF16Ptr(
        "-win\n" +
        "Runs the game in windowed mode\n" +
        "\n" +
        "-fullscreen\n" +
        "Runs the game in fullscreen mode.\n" +
        "Combine with -win for borderless windowed mode\n" +
        "\n" +
        "-modConfig filePath\n" +
        "Runs the game with selected mod (has to point to its .skudef file)\n" +
        "\n" +
        "-replayGame filePath\n" +
        "Plays replay file\n" +
        "\n" +
        "-noaudio\n" +
        "Disables game audio\n" +
        "\n" +
        "-noAudioMusic\n" +
        "Disables game music\n" +
        "\n" +
        "-silentLogin\n" +
        "Forces the game to immediately log in to a multiplayer account\n" +
        "\n" +
        "-xres number\n" +
        "Sets resolution width\n" +
        "\n" +
        "-yres number\n" +
        "Sets resolution height\n" +
        "\n" +
        "-xpos number\n" +
        "Sets horizontal offset of the window\n" +
        "\n" +
        "-ypos number\n" +
        "Sets vertical offset of the window\n" +
        "\n" +
        "-help\n" +
        "Show list of all arguments\n",
      ),
      windows.StringToUTF16Ptr("Red Alert 3"),
      windows.MB_OK,
    )
}