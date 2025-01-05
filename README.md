**chubc** is a non-interactive command line client for [Chub](https://github.com/vchimishuk/chub) audio player.

### Examples
```
$ chubc play "/Candlemass/1992 - Chapter VI/"
$ chubc volume +10
$ chubc pause
$ chubc stop
```

### Build and run
The app can build using standard `go` command.
```
$ go build
$ ./chubc
```
It is also possible to easily build a package for some operation systems. See `dist` folder in the current source distribution.

### Configuration
`chubc` does not require any specific configuration. The only configuration knob available is [Chub](https://github.com/vchimishuk/chub) server host & port target to connect to. See `man chubc` or `chub --help` for details.
