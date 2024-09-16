/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "flag"
  "os"
  "path/filepath"
  "strings"
  "strconv"
  "os/exec"
  "syscall"
  "math/rand"
  "time"
  "log/slog"
  "launcher/internal/regedit"
  "launcher/internal/splash"
  "launcher/internal/hook"
)

type Args struct {
  xres          int 
  yres          int
  xpos          int
  ypos          int
  win           bool
  fullscreen    bool
  noaudio       bool
  noAudioMusic  bool
  silentLogin   bool
  help          bool
  modConfig     string
  replayGame    string
  //ui          bool //UNIMPLEMENTED: Opens the autorun feature otherwise called by inserting the game disc > no gui but splash screen is implemented
  //getPatch    bool //DEPRECATED: Forces check for official updates > currently useless as EA's servers have been shut down
  //runver      int  //DEPRECATED: Select which version of the game to run > does not work with Steam version (ra3game.dat)
}

type Addon struct {
  Path      string  `json:"path"`
  Required  bool    `json:"required"`
}

type Config struct {
  Version     string    `json:"version"`
  Lang        string    `json:"lang"`
  Upheaval    bool      `json:"upheaval"`
  Keygen	    bool      `json:"keygen"`
  Borderless  bool      `json:"borderless"`
  Splash      bool      `json:"splash"`
  SplashImage []string  `json:"splash_image"`
  Addons      []Addon   `json:"addons"`
}

const DEFAULT_REG_PATH = "SOFTWARE\\Electronic Arts\\Electronic Arts\\Red Alert 3"

func parseArgs() (args Args) {
  flag.IntVar(&args.xres, "xres", 0, "Sets resolution width")
  flag.IntVar(&args.yres, "yres", 0, "Sets resolution height")
  flag.IntVar(&args.xpos, "xpos", 0, "Sets horizontal offset of the window")
  flag.IntVar(&args.ypos, "ypos", 0, "Sets vertical offset of the window")
  flag.BoolVar(&args.win, "win", false, "Runs the game in windowed mode")
  flag.BoolVar(&args.fullscreen, "fullscreen", false, "Runs the game in fullscreen mode. Combine with -win for borderless windowed mode")
  flag.BoolVar(&args.noaudio, "noaudio", false, "Disables game audio")
  flag.BoolVar(&args.noAudioMusic, "noAudioMusic", false, "Disables game music")
  flag.BoolVar(&args.silentLogin, "silentLogin", false, "Forces the game to immediately log in to a multiplayer account")
  flag.BoolVar(&args.help, "help", false, "Show list of all arguments")
  flag.StringVar(&args.modConfig, "modConfig", "", "Runs the game with selected mod (has to point to its .skudef file)")
  flag.StringVar(&args.replayGame, "replayGame", "", "Plays replay file")
  flag.Parse()
  
  if args.help { displayHelp() }
  
  return
}

func buildCommandLine(root string, args *Args, config *Config) string {
  
  if(config.Lang == "auto"){
    reglang := regedit.RegQueryStringValue("HKCU", DEFAULT_REG_PATH, "language")
    if len(reglang) > 0 {
        config.Lang = strings.ToLower(reglang)
    } else {
      slog.Warn("Unable to determine game's language... defaulting to \"english\"")
      config.Lang = "english"
    }
  }
  slog.Info("Language: " + config.Lang)

  skudef := "RA3_" + config.Lang + "_" + config.Version + ".SkuDef"
  cmdLine := []string{ "-config", "\"" + filepath.Join(root, skudef) + "\"" }
  
  if len(args.modConfig) > 0 {
    if !filepath.IsAbs(args.modConfig) {
        slog.Info("Looking for mod \"" + args.modConfig + "\"...")
        dirpath := filepath.Join(getUserProfilePath(),"Red Alert 3/Mods", args.modConfig)
        for _, modconfig := range findFilesWithExt(dirpath, ".skudef") {
          if strings.Contains(filepath.Base(modconfig), args.modConfig) {
            slog.Info("Found mod at \"" + modconfig + "\"")
            cmdLine = append(cmdLine, "-modConfig " + "\"" + modconfig + "\"")
            break
          }
        }
    } else {
      cmdLine = append(cmdLine, "-modConfig " + "\"" + args.modConfig + "\"")
    }
  } else if config.Upheaval {
    locations := []string{
      filepath.Join(root, "Upheaval_1.16.SkuDef"),
      filepath.Join(getUserProfilePath(), "Red Alert 3/Mods/Upheaval/Upheaval_1.16.SkuDef")}

    for _, location := range locations {
        if fileExist(location) {
          cmdLine = append(cmdLine, "-modConfig " + "\"" + location + "\"")
          break
        }
    }
  }
  
  if config.Borderless {
    cmdLine = append(cmdLine, "-win", "-fullscreen")
  }

  //passthrough
  if args.xres > 0 && args.yres > 0 {
    cmdLine = append(cmdLine,
      "-xres " + strconv.Itoa(args.xres), 
      "-yres " + strconv.Itoa(args.yres))
  }
  if args.xpos > 0 && args.ypos > 0 {
    cmdLine = append(cmdLine,
      "-xpos " + strconv.Itoa(args.xpos), 
      "-ypos " + strconv.Itoa(args.ypos))
  }
  if args.win && !config.Borderless {
    cmdLine = append(cmdLine, "-win")
  }
  if args.fullscreen && !config.Borderless {
    cmdLine = append(cmdLine, "-fullscreen")
  }
  if args.noaudio {
    cmdLine = append(cmdLine, "-noaudio")
  }
  if args.noAudioMusic {
    cmdLine = append(cmdLine, "-noAudioMusic")
  }
  if args.silentLogin {
    cmdLine = append(cmdLine, "-silentLogin")
  }
  if len(args.replayGame) > 0 {
    cmdLine = append(cmdLine, "-replayGame " + args.replayGame)
  }
  
  return strings.Join(cmdLine, " ")
}

func main(){

  args := parseArgs()
  root := locate()
  
  config, err := readJSON(filepath.Join(root, "RA3.json"))
    if err != nil { alert(err.Error()) }
  
  binary := filepath.Join(root, "/Data/", "ra3_" + config.Version + ".game")
  cmdLine := buildCommandLine(root, &args, &config)
  
  if config.Keygen { keygen() }

  cmd := exec.Command(binary)
  argv := []string{ "\"" + binary + "\"", cmdLine }
  cmd.SysProcAttr = &syscall.SysProcAttr{ CmdLine: strings.Join(argv, " ") } //verbatim arguments
  cmd.Dir = root
  cmd.Env = os.Environ()
  cmd.Stdin = nil
  cmd.Stdout = nil
  cmd.Stderr = nil
  err = cmd.Start()
    if err != nil { alert(err.Error()) }
  
  //splash screen
  exit := make(chan bool)
  if config.Splash {
    splashImage := config.SplashImage[rand.Intn(len(config.SplashImage))]
    if !filepath.IsAbs(splashImage) {
      splashImage = filepath.Join(root, splashImage)
    }
    go splash.CreateWindow(exit, cmd.Process.Pid, splashImage, 640, 480)
  } else {
    go func(exit chan bool){
      exit <- true
    }(exit)
  }
  
  //Addons
  if len(config.Addons) > 0 {
    for _, addon := range config.Addons {
          
      dylib := addon.Path
      if !filepath.IsAbs(dylib) {
        dylib = filepath.Join(root, dylib)
      }
            
      if fileExist(dylib){
        err = hook.CreateRemoteThread(uintptr(cmd.Process.Pid), dylib)
        if err != nil {
          if addon.Required {
            cmd.Process.Kill()
            alert(err.Error())
          } else {
            slog.Error(err.Error())
          }
        }
      }
    }
  }
  
  select {
    case <-exit:
      return
    case <-time.After(time.Second * 10):
      slog.Warn("Timeout")
      return
  }
}