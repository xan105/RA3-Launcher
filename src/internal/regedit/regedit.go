/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package regedit

import (
  "golang.org/x/sys/windows/registry"
)

func getHKEY(root string) registry.Key {

  var HKEY registry.Key

  switch root {
    case "HKCU": HKEY = registry.CURRENT_USER
    case "HKLM": HKEY = registry.LOCAL_MACHINE
    case "HKU": HKEY = registry.USERS
    case "HKCC": HKEY = registry.CURRENT_CONFIG
    case "HKCR": HKEY = registry.CLASSES_ROOT
  }
  
  return HKEY
}

func RegQueryStringValue(root string, path string, key string) string { // REG_SZ & REG_EXPAND_SZ

  var result string
  HKEY := getHKEY(root)

  k, _ := registry.OpenKey(HKEY , path, registry.QUERY_VALUE)
  defer k.Close()
  result, _, _ = k.GetStringValue(key)
 
  return result
}

func RegQueryStringValueAndExpand(root string, path string, key string) string { // REG_EXPAND_SZ (expands environment-variable strings)

  var result string
  HKEY := getHKEY(root)

  k, _ := registry.OpenKey(HKEY , path, registry.QUERY_VALUE)
  defer k.Close()
  str, _, _ := k.GetStringValue(key)
  result, _ = registry.ExpandString(str)
     
  return result
}

func RegWriteStringValue(root string, path string, key string, value string) {
  
  HKEY := getHKEY(root)
  
  k, _, _ := registry.CreateKey(HKEY, path, registry.ALL_ACCESS) 
  defer k.Close()
  k.SetStringValue(key, value)
}