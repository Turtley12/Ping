# Ping
Web app utility to get server list data from Minecraft servers.

## Running yourself
Before compiling and running, make sure you have Go 1.16 or higher installed.
Then run `go run .` to run, or `go build .` to build an executable.

## Config
To modify the default address and port, create a `config.json` file in the same directory as the executable.

Example config:
```json
{
	"address": "localhost",
	"port": 8080
}
```
