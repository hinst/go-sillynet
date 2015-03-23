package h_sillynet_subsystem

import "h_sillynet"
import "os/exec"
import "path/filepath"
import "strconv"

func GetApplicationDirectory() string {
	var dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
}

type Subsystem struct {
	Port    int
	AppName string
}

func (this *Subsystem) Start() {
	var appPath = GetApplicationDirectory() + filepath.Separator + this.AppName
	var command = exec.Command(appPath, strconv.Itoa(this.Port))
}
