package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

)

var (
	avicap32   = syscall.NewLazyDLL("avicap32.dll")
	proccapCreateCaptureWindowA  = avicap32.NewProc("capCreateCaptureWindowA")
	
	user32  = syscall.NewLazyDLL("user32.dll")
	procSendMessageA = user32.NewProc("SendMessageA")
)

func CaptureWebcam() {
	var name = "WebcamCapture"
	handle, _, _ := proccapCreateCaptureWindowA.Call(uintptr(unsafe.Pointer(&name)), 0, 0, 0, 320, 240, 0, 0)
	procSendMessageA.Call(handle, 0x40A, 0, 0) //WM_CAP_DRIVER_CONNECT
	procSendMessageA.Call(handle, 0x432, 30, 0) //WM_CAP_SET_PREVIEW
	procSendMessageA.Call(handle, 0x43C, 0, 0) //WM_CAP_GRAB_FRAME
	procSendMessageA.Call(handle, 0x41E, 0, 0) //WM_CAP_EDIT_COPY
	procSendMessageA.Call(handle, 0x40B, 0, 0) //WM_CAP_DRIVER_DISCONNECT
	camera, err := os.Create("Image.png")
	if err != nil {
	fmt.Println(err)
		return
	}
	clip, err := readClipboard()
	if err != nil {
	fmt.Println(err)
		return
	}
	_, err = io.Copy(camera, clip)
	if err != nil {
	fmt.Println(err)
		return
	}
	camera.Close()
}

func readClipboard() (io.Reader, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
	fmt.Println(err)
		return nil, err
	}
	f.Close()
	_, err = exec.Command("PowerShell", "-Command", "Add-Type", "-AssemblyName", fmt.Sprintf("System.Windows.Forms;$clip=[Windows.Forms.Clipboard]::GetImage();if ($clip -ne $null) { $clip.Save('%s') };", f.Name())).CombinedOutput()
	if err != nil {
	fmt.Println(err)
		return nil, err
	}
	r := new(bytes.Buffer)
	file, err := os.Open(f.Name())
	if err != nil {
	fmt.Println(err)
		return nil, err
	}
	if _, err := io.Copy(r, file); err != nil {
	fmt.Println(err)
		return nil, err
	}
	file.Close()
	os.Remove(f.Name())
	return r, nil
}

func main() {
	CaptureWebcam()
}