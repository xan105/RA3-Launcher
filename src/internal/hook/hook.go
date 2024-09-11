/*
Copyright (c) Anthony Beaumont
This source code is licensed under the MIT License
found in the LICENSE file in the root directory of this source tree.
*/

package hook

import (
  "syscall"
  "golang.org/x/sys/windows"
)

var (
    kernel32 = syscall.NewLazyDLL("kernel32.dll")
    
    pVirtualAllocEx = kernel32.NewProc("VirtualAllocEx")
)

func CreateRemoteThread(pid uintptr, path string) error {

  //Opens a handle to the target process with the needed permissions
  hProcess, err := windows.OpenProcess(
    windows.PROCESS_CREATE_THREAD | 
    windows.PROCESS_VM_OPERATION | 
    windows.PROCESS_VM_WRITE | 
    windows.PROCESS_VM_READ |
    windows.PROCESS_QUERY_INFORMATION,
    false,
    uint32(pid),
  )
  if err != nil {
    return err
  }

 //Allocates virtual memory for the file path
  lpBaseAddress, _, err := pVirtualAllocEx.Call(
    uintptr(hProcess), 
    0, 
    uintptr(len(path)+1), 
    windows.MEM_RESERVE | windows.MEM_COMMIT, 
    windows.PAGE_EXECUTE_READWRITE,
  )
 
  //Converts the file path to type *byte
  lpBuffer, err := windows.BytePtrFromString(path)
  if err != nil {
    return err
  }
 
 //Writes the filename to the previously allocated space
  lpNumberOfBytesWritten:= uintptr(0)
  err = windows.WriteProcessMemory(
    hProcess, 
    lpBaseAddress, 
    lpBuffer, 
    uintptr(len(path)+1), 
    &lpNumberOfBytesWritten,
  )
  if err != nil {
    return err
  }
 
 //Gets a pointer to the LoadLibrary function
  LoadLibAddr, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
  if err != nil {
    return err
  }
 
 //Creates a remote thread that loads the DLL triggering it
  handle, _, err := kernel32.NewProc("CreateRemoteThread").Call(uintptr(hProcess), 0, 0, LoadLibAddr, lpBaseAddress, 0, 0)
  if handle == 0 {
    return err
  }

  defer syscall.CloseHandle(syscall.Handle(handle))
  
  return nil
}