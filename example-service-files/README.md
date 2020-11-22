# Example service files (For systems with systemd based service managers)


## Setting up a dmsg-http server service

First, you must build the [dmsg-http server](https://github.com/skycoin/dmsg/tree/master/examples/dmsgget/dmsg-example-http-server) from Skycoin's [repo](https://github.com/skycoin/dmsg)
- run "go get github.com/skycoin/dmsg" to download the repo to your src go dir (e.g. ~/go/src/)
- Run navagate to your skycoin src dir (e.g. ~/go/src/github.com/skycoin/dmsg/)  "make build"
- Navigate to the ./bin folder and move the dmsg-server binary to a PATH location (e.g. /usr/bin/)
- Generate a keypair for your dmsg-http server by going to ~/go/src/github.com/skycoin/dmsg/examples/dmsgget/gen-keys
    - Run "go run gen-keys.go"
        - **Example:**
        - PK: 02bf0e6c944bed0c08f9567196cddcab52c2c04d4d822e5dc1020f5e6f949c2016
        - SK: 1a9f70958b1d6923e7ac1394626e63033f033f65d1476f832e5f16bb79786b73

    - Save these keys for later use

You now have enough to setup the dmsg-http service. Using the dmsg-http.service file as a template, modify the following line to your preference:

```sh
$ ExecStart=dmsg-http-server --dir [root dir for dmsg-http files] --sk [private key generated goes here]
$ (e.g. ExecStart=dmsg-http-server --dir /srv/dmsg-http-files/ --sk 1a9f70958b1d6923e7ac1394626e63033f033f65d1476f832e5f16bb79786b73)
```

Now copy the modified dmsg-http.service file to your systemd service file location (/etc/systemd/system/). After moving the file to the service location, run: 

```sh
$ systemctl daemon-reload   # reload service files
```

Then run:
```sh
$ systemctl start dmsg-http.service   # start the service
```

Finally run:
```sh
$ systemctl status dmsg-http.service   # see service status
```
The dmsg-http.service should state that it's active and running.

---

## Setting up a dmsg-indexing service

Compile the dmsg-indexer by running:
```sh
$ go build dmsg-indexer.go
```

Now move the dmsg-indexer binary to a location in your PATH (e.g. /usr/bin/). We're now ready to update the service file for the dmsg-indexer. Open up the dmsg-indexer.service file and change the following lines to match your config:



```sh
$ WorkingDirectory=[same dir that the dmsg-http server is set to serve]
$ ExecStart=dmsg-indexer [indexer interval]
$
$ Example:
$ WorkingDirectory=/srv/dmsg-http-files/
$ ExecStart=dmsg-indexer 360 # Run every 360s or 6 mins
```

After the edits to the template have been made, move the dmsg-indexer.service to your systemd service location (/etc/systemd/system/).
After moving the file to the service location, run: 

```sh
$ systemctl daemon-reload   # reload service files
```

Then run:
```sh
$ systemctl start dmsg-indexer.service   # start the service
```

Finally run:
```sh
$ systemctl status dmsg-indexer.service   # see service status
```
The dmsg-indexer.service should state that it's active and running.
