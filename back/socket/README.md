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

## Packet IDs

Since our packets are stored in a Go `byte` (better known as `uint8`), we can only store values between `0` and `255`.

To avoid setting ranges of IDs for each type of packet, I've used some basic bit manipulation tricks to check the type without running any slow math operators (like `%` or `/`), or using any `>=` or `<=` operators.

First two bits of the packet ID are used to define the type of the packet :

- `00` : Generic packet
- `01` : Packet regarding a specific client
- `10` : Reserved (not used yet)
- `11` : Reserved (not used yet)

## Packets

### `0x80` - Welcome

Sent by the server to every clients (except the one who just connected) when a new client connects.

```
Client UUID (36 bytes)
```

### `0x81` - Cya

Sent by the server to every clients (except the one who just disconnected) when a client disconnects.

```
Client UUID (36 bytes)
```