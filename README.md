# dmsg-http-cli-browser (Deprecated - Not Working - Switch to VPN+sFTP)


**What is it?**
It’s system to implement a simple and anonymous file navigation/downloading using Skycoin's Skywire network. It builds on top of Skycoin's [dmsg-http](https://github.com/skycoin/dmsg/tree/master/examples/dmsgget/dmsg-example-http-server) and [dmsgget](https://github.com/skycoin/dmsg/blob/master/docs/dmsgget.md) binaries compiled from their [dmsg](https://github.com/skycoin/dmsg) project.


# Indexer (On Server Hosting DMSG-HTTP)
The indexer is setup to index the root directory of the dmsg-http server. The indexer will create an "index" file containing all the files in the directory tree.
The client requests the index file from the server in the following root location dmsgget dmsg://[Server public key]:80/index.


**Indexer features**
- Indexer will scan the root working directory in set intervals (default 30 seconds). This interval can be set manually by providing an integer number in seconds as an argument.
- Does not index empty directories or the index file itself.
- Keyword filter

# Client
The client operates as a simple bookmark manager for public keys that are running dmsg-http services and as a file downloader (dmsgget wrapper). Server public keys are saved locally and can be given a friendly name. After a successfully fetch of the index from the server, the client is able browse the directory structure of the server using a simple interface:

**Example Server List:**

![Server List Example](https://github.com/ThatOneRuffian/dmsg-http-cli-browser/blob/master/README_files/wiki/server_list.png?raw=true)


**Example File Navigation:**

![Server Download Index Example](https://github.com/ThatOneRuffian/dmsg-http-cli-browser/blob/master/README_files/readme_fig1.png?raw=true)


**Client features**
- Server public key management
- Simple UI
- Paginated server/file browsing
- Simultaneous downloads
- Search current directory
- UI scales to terminal size (Optional dependency: tput)

#### ***Windows users looking for skycoin dmsgget, dmsg-server, and gen-key binaries, please check out this forked [skycoin/dmsg](https://github.com/ThatOneRuffian/dmsg/releases) repo for windows compatible binaries.***
---
### ***If you like what you see and you want to see more please donate. Thanks!***

***Skycoin:*** 2ZtDYzoUBESccvK5mzDBHqKaAPvjwVmzESs
