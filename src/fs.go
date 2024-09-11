/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "os"
  "io"
  "io/fs"
  "path/filepath"
  "encoding/json"
  "errors"
  "launcher/internal/regedit"
)

func locate() string{
  process, err := os.Executable()
    if err != nil { alert(err.Error()) }
  return filepath.Dir(process)
}

func readJSON(filepath string) (config Config, err error) {

  //default values
  config.Version = "1.12"
  config.Lang = "auto"
  config.Upheaval = false
  config.Keygen = true
  config.Borderless = false
  config.Splash = true
  config.SplashImage = []string{ "Launcher/splash.bmp" }
  config.Addons = []Addon{}

  file, err := os.Open(filepath)
    if err != nil { return }
  defer file.Close()
  
  bytes, err := io.ReadAll(file)
    if err != nil { return }

  err = json.Unmarshal(bytes, &config)
    if err != nil { return }
    
  if config.Version == "" {
    config.Version = "1.12"
  }
  if config.Lang == "" {
    config.Version = "auto"
  }
  
  if len(config.SplashImage) == 0 {
    config.SplashImage = []string{ "Launcher/splash.bmp" }
  }

  return
}

func fileExist(path string) bool {
  target, err := os.Stat(path)
  if err == nil {
    return !target.IsDir()
  }
  if errors.Is(err, os.ErrNotExist) {
    return false
  }
  return false
}

func findFilesWithExt(dirpath string, ext string) []string {
  var matches []string
  filepath.WalkDir(dirpath, func(path string, info fs.DirEntry, e error) error {
      if e != nil { return e }
      if filepath.Ext(info.Name()) == ext {
         matches = append(matches, path)
      }
      return nil
   })
  return matches
}

func getUserProfilePath() string {

  const PATH = "Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\User Shell Folders"
  keys := []string{
    "{F42EE2D3-909F-4907-8871-4C22FC0BF756}", //win10
    "Personal"}
    
  for _, key := range keys {
    value := regedit.RegQueryStringValueAndExpand("HKCU", PATH, key)
    if len(value) > 0 {
      return value
    }
  }

  return filepath.Join(os.Getenv("USERPROFILE"), "Documents")
}