# Molly's Cafe 🍵

This repo contains the source code for [Molly's Cafe](https://mana.rip/) which, for now is just a status webpage, but I hope I'll be able to realize my goal of making it a neocity !

# Note

## Platform support

The Go code can only be compiled on Linux (maybe MacOS ?) due to :

1. The usage of Go's `syscall` package
2. The usage of a lot of `/proc` / `/sys` file parsing
3. Relying on `systemd` for service management

## Code quality

Some parts of the code are not very well written, mostly because I was sometime lazy / in a hurry.

This section will be removed once I've cleaned up the code.

## Lack of documentation / comments

Coming really soon, I promise !

## Looking for bugs ?

Go ahead, break everything you can (well not too much), I'll be happy to fix it and learn more!

# Exploring the code

If you're interested in some specific part of the code, here's a simple table of content :

- [Netcode](#netcode)
- [Monitoring](#monitoring)
- [2D Rendering](#2d-rendering)

## Netcode <a name="netcode"></a>

The netcode is the part of the code that handles (for the moment) the websocket connections.

- [clients.go](server/socket/clients.go) : Handles connections and handles broadcasting / mutexes
- [proto.go](server/socket/proto.go) : Creation, modification of custom packets
- [packets.go](server/socket/packets.go) : Hashmap of packets, used by watchdogs to quickly edit packets and reflect the changes to the clients
- [main.go](server/main.go) : Handling of HTTP -> Websocket upgrade, along with other private APIs and static file serving

## Monitoring <a name="monitoring"></a>

The monitoring part of the code is where the program is gathering data from other programs / from the kernel.

### Services

The `services` part of the code, is where we gather the state of specific services running through `systemd`, and where we gather the state of `docker` containers.

This part of the code is completely done using polling, so I made the choice of also making it real-time, meaning that any change to a `docker` container or a `systemd` service will be broadcasted instantly.

### Watchdogs

Currently each `watchdog` is a goroutine, and each `watchdog` is started in the `init` function.

A `watchdog` is a goroutine that will poll a specific value, and if the value is different from the previous one, it will update the corresponding packet, and broadcast it to the clients, and then sleep for a specific amount of time.


| File | Description | Remarks |
| :--- | :--- | ---: |
| [Containers.go](server/watchdogs/Containers.go) | Watches every `docker` containers | - |
| [CPUTemp.go](server/watchdogs/CPUTemp.go) | Watches the CPU temperature (compatible with `k10temp` / `coretemp`) | - |
| [DirtyMem.go](server/watchdogs/DirtyMem.go) | Watches the amount of dirty memory | - |
| [DiskSpace.go](server/watchdogs/DiskSpace.go) | Watches the amount of free disk space | Use a translation file instead of hardcoding them |
| [IdleUptime.go](server/watchdogs/IdleUptime.go) | Watches the amount of time the system has been idle | - |
| [InternetSpeed.go](server/watchdogs/InternetSpeed.go) | Downloads a file every hour and measure download speed | Probably cycle through multiple servers to be a good netizen |
| [LoggedUsers.go](server/watchdogs/LoggedUsers.go) | Watches the amount of logged users | Disabled for now, because I don't want to fork `who`, and must parse `/var/run/wtmp` instead |
| [ManualService.go](server/watchdogs/ManualService.go) | Watches the state of a `systemd` service | Use an external file to store the services to watch instead of hardcoding them |
| [MemUsage.go](server/watchdogs/MemUsage.go) | Watches the amount of memory used | - |
| [OpenFiles.go](server/watchdogs/OpenFiles.go) | Watches the amount of opened fds | - |
| [TcpUdp.go](server/watchdogs/TcpUdp.go) | Watches the amount of opened TCP / UDP sockets | - |
| [RunningProcesses.go](server/watchdogs/RunningProcesses.go) | Watches the amount of running processes | - |
| [MonitorSchoolProjects.go](server/watchdogs/MonitorSchoolProjects.go) | Polls a directory containing school projects to dynamically update a table in the front-end | - |
| [PistachePosts.go](server/watchdogs/PistachePosts.go) | Polls a directory containing posts to dynamically update a list of blog-post in the front-end | - |
| [Weather.go](server/watchdogs/Weather.go) | Checks the weather of any configured city every 5 minutes | - |

#### Implementation notes

Every `watchdog` is a goroutine, and every `watchdog` is started in the `init` function.

They're designed to be as simple, efficient and modular as possible (for my use case).

For efficiency sake, I avoid as much as possible to use any forks / execs. But there's still room for improvement! (just too lazy to do it right now, and it works pretty well for now.)

Docker containers states are watched through [docker events](https://docs.docker.com/engine/api/v1.43/#tag/System/operation/SystemEvents).

`systemd` services states are watched through a probably overthought method using `inotify` on the `/run/systemd/units` directory.

#### Flaws

- Handling of errors in the code is poorly done