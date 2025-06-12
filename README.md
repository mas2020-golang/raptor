# Raptor
Raptor is a CLI application to manage and in a fast and smart way your secrets. It is native app for Mac OS, Linux and Windows.

## Installing raptor

We tested `raptor` on Linux and Mac. To work on Linux you need to get installed one of the following software:
- xclip
- xsel

## Environment variables
Follow the list of the env variables you can use:
- `CRYPTEX_FOLDER`: it is folder where to store the boxes. You can set this variable to override the standard behaviour (default is searching the boxes in the `$HOME/.cryptex/boxes` folder).
- `CRYPTEX_BOX`: setting this variable if you can avoid to pass the `--box` flag
- `RAPTOR_LOGLEVEL`: set it for logging purposes and debugging sessions. Accepted values are: `debug`, `info`, `warn`, `error`. Default level set to error.
- `RAPTOR_TIMEOUT_SEC`: number of seconds of inactivity before the application exits (default 10 mins).
