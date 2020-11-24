# Album info spotify

This is a simple command line app that shows info of the current playing album listen on spotify desktop app, also have the possibility to search info by enter artistsName and albumName by hand

## Install
NOTE: Works on ubuntu and MacOS.

Go to releases page and download file "spotify-info-*"  for your current operaing system 

https://github.com/ernesto27/album-info-spotify/releases

### Desktop app


Make file executable
```sh
$ chmod +x spotify-info-linux
```

#### Use
```sh
$ ./spotify-info-linux
```

![](https://i.ibb.co/55FsC6h/ui-1.png)

![](https://i.ibb.co/TtczPNW/ui-2.png)



### Command line app
Go to releases page and download file "cmd-spotify-info".

Make file executable
```sh
$ chmod +x cmd-spotify-info
```

#### Use

If you want to have info about a current playing song on spotify desktop app
```sh
$ ./cmd-spotify-info
```

Info about a spececific album

```sh
$ ./cmd-spotify-info "megadeth" "rust in peace"
```

![](https://raw.githubusercontent.com/ernesto27/album-info-spotify/main/screenshots/Screenshot%20from%202020-11-15%2020-15-29.png)
