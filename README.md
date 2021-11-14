# GoBLight - An extra simple backlight program for GNU/Linux
The simplicity of this program derives from the fact that it only utilizes `/sys/class/backlight/*` files to control backlight instead of using bloated Xorg drivers.  
Change the driver in `main.go` to your own one.  
This program needs to be run with the SUID bit set to root  
Run `sudo make install` to automatically set SUID bit and perms (binary will be installed at `/usr/bin/goblight`).  
