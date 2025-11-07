pkgname=dotloader
pkgver=0.1
pkgrel=1
pkgdesc="Simple dotfiles manager working with rsync and git written in go."
arch=('x86_64')
license=('CC BY-NC 4.0')
depends=('rsync' 'git')
makedepends=('go')
source=('main.go' 'go.mod')
sha256sums=('f355af2df0bede5b2623f4066b2f2ed8db7b52bfe83d0c178e76d2180a1243d5'
            'bb2d6c133dfe624a70ea39841194bc3dd8730cd9d6b27292c6d46104e8305cac')

build() {
  go build -o dotloader
}

package() {
  install -Dm755 "$srcdir/dotloader" "$pkgdir/usr/bin/dotloader"
}
