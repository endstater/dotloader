package main

import (
	"bufio"
	"os/exec"
	"strings"
	"fmt"
	"syscall"
	"os"
)

type Dotloader struct {
	homedir   string
	conf      map[string]string
}

func NewDotloader()(*Dotloader, error) {
	var d Dotloader
	home, err := os.UserHomeDir()
  if err != nil {
		panic(err)
  }
	d.homedir = home
	d.conf = map[string]string{}
	err = d.BuildConfig()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	err = d.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return &d, nil
}

func (d *Dotloader) ReadConfig() error {
	path := d.homedir + "/.config/dotloader/dotloader.conf"
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Cannot open file: %v", path)
	}
	scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
		equalation := strings.Split(scanner.Text(), "=")
		if len(equalation) <= 1 { continue }
		key := strings.ReplaceAll(equalation[0], " ", "")
		value := strings.ReplaceAll(equalation[1], " ", "")
		d.conf[key] = value
  }
  if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error while scanning: %v", path)
  }
	return nil
}

func (d *Dotloader) GetCopies() ([]string, error) {
	path := d.homedir + "/.config/dotloader/listen-dirs"
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("Cannot open file: %v", path)
	}
	scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
	data := []string{}
  for scanner.Scan() {
		if scanner.Text()[0] == '#' { continue }
		data = append(data, scanner.Text())
  }
  if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error while scanning: %v", path)
  }
	return data, nil
}

func (d *Dotloader) Sync() error {
	repo := d.conf["dotfiles-dir"]
	copies, err := d.GetCopies()
	if copies == nil || repo == "" || err != nil {
		return fmt.Errorf("Error while reading config files")
	}

	script := "#!/sbin/sh\n"
	for _, v := range copies {
		if v[len(v)-1] == '/' {
			v = v[:len(v)-1]
		}
		println(v)
		info, _ := os.Stat(os.ExpandEnv(v))
		if os.Getuid() != int(info.Sys().(*syscall.Stat_t).Uid) {
			script += "sudo "
		}

		vpath := strings.Split(v, "/")
		destination := strings.Join(vpath[:len(vpath)-1], "/")
		if destination[0] != '$' &&
			 destination[0] != '/' {
			destination = "/" + destination
		}
		v = os.ExpandEnv(v)
		repo = os.ExpandEnv(repo)
		if _, err := os.Stat(v); err != nil {
			continue
		}
		_, err = exec.Command("rsync", "-a", "--delete", v, repo).CombinedOutput()

		if err != nil {
			return fmt.Errorf("%v", err)
		}
		script += "rsync -a --delete " + vpath[len(vpath)-1] + " " + destination + "\n"
		os.WriteFile(repo+"/load.sh", []byte(script), 0744)
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
	if _, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func (d *Dotloader) BuildConfig() error {
	baseDir := d.homedir+"/.config/dotloader"
	_, err := os.Stat(baseDir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("Failed to check directory %s: %w", baseDir, err)
	}

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create directory %s: %w", baseDir, err)
	}

	err = os.WriteFile(baseDir+"/listen-dirs", []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %w", baseDir + "/listen-dirs", err)
	}

	err = os.WriteFile(baseDir+"/dotloader.conf", []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %w", baseDir + "/dotloader.conf", err)
	}

	return nil
}

func (d *Dotloader) Help() {
	fmt.Println("\ndotloader <option>")
	fmt.Println("options:")
	fmt.Println("\tsync, s\tread $HOME/.config/dotloader/listen-dirs for dirs, \n\t\tsync them with local dotfiles repo, \n\t\tmake dotfiles-dir/load.sh")
	fmt.Println("\tload, l\tlaunch load.sh that sync local dotfiles repo with system")
	fmt.Println("\tversion, v\tprint version")
	fmt.Println("\thelp, h\tprint this information")
}
