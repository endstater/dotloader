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

func (d *Dotloader) ReadConfig() {
	file, err := os.Open(d.homedir + "/.config/dotloader/dotloader.conf")
	defer file.Close()
	if err != nil{
		panic(err)
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
		panic(err)
  }
}

func (d *Dotloader) GetCopies() []string {
	file, err := os.Open(d.homedir + "/.config/dotloader/listen-dirs")
	defer file.Close()
	if err != nil{
		return nil
	}
	scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
	data := []string{}
  for scanner.Scan() {
		data = append(data,scanner.Text())
  }
  if err := scanner.Err(); err != nil {
		panic(err)
  }
	return data	
}

func (d *Dotloader) Sync()  {
	repo := d.conf["dotfiles-dir"]
	copies := d.GetCopies()
	if copies == nil || repo == ""{
		return
	}
	
	for _, v := range copies {
		if v[len(v)-1] == '/'{
			v = v[:len(v)-1]
		}
		v = os.ExpandEnv(v)
		repo = os.ExpandEnv(repo)
		exec.Command("rsync","-a","--delete",v,repo).Run()
	}
	
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
		dotloader.Sync()
	default:
		fmt.Println("Unknown option:",args[1])
		return
	}
}
