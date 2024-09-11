/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package main

import(
  "math/rand"
  "log/slog"
  "launcher/internal/regedit"
)

func randAlphaNumString(length int) string {
  //cf: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
  
  const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  const IdxBits = 6
  const IdxMask = 1<<IdxBits - 1
  const IdxMax = 63 / IdxBits

  bytes := make([]byte, length)
  for i, cache, remain := length-1, rand.Int63(), IdxMax; i >= 0; {
    if remain == 0 {
      cache, remain = rand.Int63(), IdxMax
    }
    if idx := int(cache & IdxMask); idx < len(charset) {
      bytes[i] = charset[idx]
      i--
    }
    cache >>= IdxBits
    remain--
  }
  return string(bytes)
}

func keygen(){ //Elevated perm required
  
  /*
    LAN play requires each player to have a different CD Key.
    Steam no longer populates the %CDKEY% variable in the registry (is it due to EA's servers shutdown ?)
    This is a problem because each time the game is ran by Steam; Steam will write "%CDKEY%" as the user's CD key.
    Every (Steam) players will therefore end up with the same CD Key: "%CDKEY%" and are unable to LAN play with each others.
    The registry value is read on game launch. It is uneffective to change it after the game is ran by Steam.
    
    Since EA shutdown their servers for this game. Having a "real" and legit CD Key is of no concerns here.
    Generate a random alphanum string and write it to the registry if the CD Key is "%CDKEY%" or there is none.
    
    Unfortunately writing to HKLM (registry) requires elevated privileges (UAC/Adming rights).
  */
  
  registration := regedit.RegQueryStringValue("HKLM", DEFAULT_REG_PATH, "Registration")
  if registration == "" {
    registration = DEFAULT_REG_PATH + "\\ergc"
    regedit.RegWriteStringValue("HKLM", DEFAULT_REG_PATH, "Registration", registration) //elevated
  }
  
  regkey := regedit.RegQueryStringValue("HKLM", registration, "")
  
  if regkey == "%CDKEY%" || regkey == "" {
    slog.Warn("RA3 CD Key not found ! This will interfere with LAN play")
    
    key:= randAlphaNumString(20)
    regedit.RegWriteStringValue("HKLM", registration, "", key) //elevated
    slog.Info("Generated random CD Key to enable LAN play")
  } 
}