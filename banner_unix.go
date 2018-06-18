// +build linux darwin

package main

import (
	"fmt"
)

func bannerGet() string {
	banner := o.magenta.Sprintf(`
    ___________________       __               __    
   / ____/ ____/ ____/ |     / /___  _________/ /____
  / __/ / /_  / /_   | | /| / / __ \/ ___/ __  / ___/
 / /___/ __/ / __/   | |/ |/ / /_/ / /  / /_/ (__  ) 
/_____/_/   /_/      |__/|__/\____/_/   \__,_/____/ `)
	banner += "\n" + `Author:  Morgaine Timms   (@sectorsect)
License: MIT
Warning: Some of the following options when used in combination can
         significantly weaken the pass-phrases generated. 
         You probably know what you are doing though, yeah?`

	banner += fmt.Sprintf("\n\n")
	return banner
}
