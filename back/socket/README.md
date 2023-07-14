# Sockets

## File structure

```
$ tree -L 2
.
├── clients.go          <-- Client management, connection, disconnection and broadcasting
├── interpretation.go   <-- Interpretation of received packets
├── packets             <-- Every packet is defined here
│   ├── cya.go          <-- Disconnection packet
│   ├── ids.go          <-- Every packet ID is defined here
│   └── welcome.go      <-- Connection packet
└── README.md           <-- You are here !
```

## Protocol

The protocol I've designed is very basic :

Generic packet :

```
Packet ID (1 byte) | Data (n bytes)
```

Packet regarding a specific client :

```
Packet ID (1 byte) | Client UUID (36 bytes) | Data (n bytes)
```

## Packets

### Welcome

Sent by the server to every clients (except the one who just connected) when a new client connects.

```
Packet ID (1 byte) | Client UUID (36 bytes)
```

### Cya

Sent by the server to every clients (except the one who just disconnected) when a client disconnects.

```
Packet ID (1 byte) | Client UUID (36 bytes)
```