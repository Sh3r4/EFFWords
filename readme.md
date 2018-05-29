# EFF words

```txt
    ___________________       __               __
   / ____/ ____/ ____/ |     / /___  _________/ /____
  / __/ / /_  / /_   | | /| / / __ \/ ___/ __  / ___/
 / /___/ __/ / __/   | |/ |/ / /_/ / /  / /_/ (__  )
/_____/_/   /_/      |__/|__/\____/_/   \__,_/____/
License: MIT
```

EFFWords is a passphrase generator based on the long wordlist from [this EFF post](https://www.eff.org/dice).

It is ported from my npm project easyphrase which adds passphrase generation to webpages, but has much more functionality aimed at administrators.

## Installation

IMO, and because I am a golang fangirl, the best way to get EFFWords is through go get:

```
go get github.com/sectorsect/EFFWords
```

I'm not about to leave non-IT peeps out in the cold though!
I will endeavour to make releases with cross-compiled builds for:

* Windows x86
* Macos   x86
* Linux   x86

That said, since this is just a personal project I am putting up on the internet, I won't have tested it on all the release platforms.

## Usage

``` txt
    ___________________       __               __
   / ____/ ____/ ____/ |     / /___  _________/ /____
  / __/ / /_  / /_   | | /| / / __ \/ ___/ __  / ___/
 / /___/ __/ / __/   | |/ |/ / /_/ / /  / /_/ (__  )
/_____/_/   /_/      |__/|__/\____/_/   \__,_/____/
Author:  Morgaine "sectorsect" Timms
License: MIT
Warning: Some of the following options when used in combination can
         significantly weaken the pass-phrases generated.
         You probably know what you are doing though, yeah?

Things It Does:
  -c, --caps-position int
        Capitalise the word at this index. (default -1)
  -M, --maximum-chars int
        Maximum characters for passphrases. Truncates. (default -1)
  -m, --minimum-chars int
        Minimum characters for passphrases. (default -1)
  -o, --output-to-file string
        Filepath to output to
  -C, --prevent-caps
        Prevents capitalisation in passphrase
  -I, --prevent-int
        Prevents integers in passphrase
  -S, --prevent-special
        Prevents special chars in passphrase
  -q, --quantity int
        Number of passphrases to generate (default 1)
  -R, --random-caps
        Capitalise a random character in each passphrase
  -v, --verbose
        Show all the things
  -w, --wordcount int
        Number of words per passphrase. (default 4)
```

## Disclaimer

I lay no claim to the included wordlist created by the EFF. I am using it in good faith, with full attribution to the Electronic Frontiers Foundation. I believe the list to be licensed under the CC 3.0 Attribution license [https://creativecommons.org/licenses/by/3.0/us/](https://creativecommons.org/licenses/by/3.0/us/)