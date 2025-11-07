# Dotloader <img src="https://pkg.go.dev/static/shared/logo/go-blue.svg" alt="Description" height="24pt">
Simple dotfiles manager working with rsync and git written in go.
## First start
add dirs to `$HOME/.config/dotloader/listen-dirs`
example:
```
$HOME/.config/nvim
$HOME/.config/yay
$HOME/file.txt
/etc/NetworkManager/NetworkManager.conf
```
add repo to `$HOME/.config/dotloader/dotloader.conf`
example:
```
dotfiles-dir = $HOME/projects/my-dotfiles
```
## Usage
```
dotloader <option>
```
options:
- `sync` `s` - read `$HOME/.config/dotloader/listen-dirs` for dirs, sync them with local dotfiles repo, make dotfiles-dir/load.sh
- `load` `l` - launch `load.sh` that sync local dotfiles repo with system
- `version` `v` - print version
- `help` `h` - print manual
## Build and run
```
go build -o out
out/dotloader <option>
```
or
```
go run *go <option>
```
