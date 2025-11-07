package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"fmt"
)

type Dotloader struct {
	homedir   string
	conf      map[string]string
}

/*
TODO:
rsync options in config
git sync
improve project repo
*/

func NewDotloader() *Dotloader {
	var d Dotloader
	home, err := os.UserHomeDir()
  if err != nil {
		panic(err)
  }
	d.homedir = home
	d.conf = map[string]string{}
	d.ReadConfig()
	return &d
}

func (d *Dotloader) ReadConfig() error {
	path := d.homedir + "/.config/dotloader/dotloader.conf"
	file, err := os.Open(path)
	defer file.Close()
	if err != nil{
		return fmt.Errorf("Cannot open file: %v", path)
	}
	scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
		equalation := strings.Split(scanner.Text(),"=")
		if len(equalation) > 1{
			key := strings.ReplaceAll(equalation[0], " ","")
			value := strings.ReplaceAll(equalation[1], " ","")
			d.conf[key] = value
		}
  }
  if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error while scanning: %v",path)
  }
	return nil
}

func (d *Dotloader) GetCopies() ([]string,error){
	path := d.homedir + "/.config/dotloader/listen-dirs"
	file, err := os.Open(path)
	defer file.Close()
	if err != nil{
		return nil, fmt.Errorf("Cannot open file: %v", path)
	}
	scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
	data := []string{}
  for scanner.Scan() {
		data = append(data,scanner.Text())
  }
  if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error while scanning: %v",path)
  }
	return data,nil
}

func (d *Dotloader) Sync() error {
	repo := d.conf["dotfiles-dir"]
	copies, err := d.GetCopies()
	if copies == nil || repo == "" || err != nil{
		return fmt.Errorf("Error while reading config files")
	}

	script := "#!/sbin/sh\n"
	for _, v := range copies {
		if v[len(v)-1] == '/'{
			v = v[:len(v)-1]
		}
		v = os.ExpandEnv(v)
		repo = os.ExpandEnv(repo)
		var err any = exec.Command("rsync","-a","--delete",v,repo).Run().Error()

		if errstr,ok := err.(string); ok{
			return fmt.Errorf("%v",errstr)
		}
		vpath := strings.Split(v, "/")
		destination := "/" + strings.Join(vpath[1:len(vpath)-1],"/")
		script += "rsync -a --delete "+vpath[len(vpath)-1] + " " + destination + "\n"
		os.WriteFile(repo+"/load.sh",[]byte(script),0744)
	}
	return nil
}

func (d *Dotloader) Load() error {
	repo := os.ExpandEnv(d.conf["dotfiles-dir"])
	if repo == "" {
		return fmt.Errorf("Error while reading config files")
	}
	cmd := exec.Command(repo+"/load.sh")
	cmd.Dir = repo
	var err any = cmd.Run().Error()
	if errstr,ok := err.(string); ok{
		return fmt.Errorf("%v",errstr)
	}
	return nil
}

func main(){
	dotloader := NewDotloader()
	
	args := os.Args
	if len(args) < 2{
		dotloader.Sync()
		return
	}

	switch args[1]{
	case "sync":
		err := dotloader.Sync()
		if err != nil{
			fmt.Printf("%v",err)
		}
	case "load":
		err := dotloader.Load()
		if err != nil{
			fmt.Printf("%v",err)
		}
	default:
		fmt.Println("Unknown option:",args[1])
	}
}
