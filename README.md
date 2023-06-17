# magnetico
*Autonomous (self-hosted) BitTorrent DHT search engine suite.*

![Flow of Operations](/doc/operations.svg)

magnetico is the first autonomous (self-hosted) BitTorrent DHT search engine suite that is *designed for end-users*. The suite consists of two packages:

- **magneticod:** is the daemon that crawls the BitTorrent DHT network in the background to discover info hashes and fetches metadata from the peers.
- **magneticow:** is a lightweight web interface to search and to browse the torrents that its counterpart discovered.

Both programs, combined together, allows anyone with a decent Internet connection to access the vast amount of torrents waiting to be discovered within the BitTorrent DHT space, *without relying on any central entity*.

**magnetico** liberates BitTorrent from the yoke of centralised trackers & web-sites and makes it
*truly decentralised*. Finally!

## Features
- Easy installation & minimal requirements:
  - Easy to build golang static binaries.
  - Root access is *not* required to install or to use.
- Near-zero configuration:
  - Both programs work out of the box, and **magneticow** can be used without a web-server too.
  - Detailed, step-by-step manual to guide you through the installation.
- No reliance on any centralised entity:
  - **magneticod** trawls the BitTorrent DHT by "going" from one node to another, and fetches the
    metadata using the nodes without using trackers.
- Resilience:
  - Unlike client-server model that web applications use, P2P networks are *chaotic* and
    **magneticod** is designed to handle all the operational errors accordingly.
    - Currently on paper, wait for the v1.0!
- High performance implementation in Go:
  - **magneticod** utilizes every bit of your resources to discover as many infohashes & metadata as
    possible.
- Built-in lightweight web interface:
  - **magneticow** features a lightweight web interface to help you access the database without
    getting on your way.

### Screenshots

| ![The Homepage](/doc/homepage.png) | ![Searching for torrents](/doc/search.png) | ![Search result](/doc/result.png) |
|:-------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------------------------------------------------------------------:|
|                                                                     __The Homepage__                                                                    |                                                                     __Searching for torrents__                                                                    |                                                     __Viewing the metadata of a torrent__                                                     |

## Why?
BitTorrent, being a distributed P2P file sharing protocol, has long suffered because of the
centralised entities that people depended on for searching torrents (websites) and for discovering
other peers (trackers). Introduction of DHT (distributed hash table) eliminated the need for
trackers, allowing peers to discover each other through other peers and to fetch metadata from the
leechers & seeders in the network. **magnetico** is the finishing move that allows users to search
for torrents in the network, hence removing the need for centralised torrent websites.
