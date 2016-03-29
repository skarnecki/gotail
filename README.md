#Tail -f in you browser!
Watch file changes in browser. Perfect for tailing logs. Works similar to tail -f but in browser.
HTTPS & Secure WebSockets and Basic Auth supported.

## Usage

gotail \<file path> \<options>

* --number   Starting lines number
* --host     Listening host. Default 0.0.0.0
* --port     Listening port. Default 9001
* --cert     Path to SSL certificate file (HTTPS)
* --key      Path to SSL key file (HTTPS)
* --user     Basic Auth username
* --password Basic Auth password

## Install
``` curl https://raw.githubusercontent.com/skarnecki/gotail/master/install.sh | sh
