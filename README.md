Write.as
========
[![Build Status](https://travis-ci.org/writeas/nerds.svg)](https://travis-ci.org/writeas/nerds) [![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Public Slack discussion](http://slack.write.as/badge.svg)](http://slack.write.as/)

This is a simple telnet-based interface for publishing text. Users connect and paste / type what they want to publish. Upon indicating that they're finished, a link is generated to access their new post on the web.

![Write.as telnet server](https://github.com/writeas/nerds/raw/master/nerds.png)

## Try it
**Or not :(**. We had to [shut it down](https://twitter.com/writeas__/status/790356847526027264) because it was getting DDoSed too much. But you can still [run it yourself](#run-it-yourself).

```
telnet nerds.write.as
```

## Run it yourself
```
Usage:
  nerds [options]

Options:
  --debug
       Enables garrulous debug logging.
  -o   
       Directory where text files will be stored.
  -s
       Directory where required static files exist (like the banner).
  -h
       Hostname of the server to rsync saved files to.
  -p
       Port to listen on.
```

The default configuration (without any flags) is essentially:

```
nerds -o /var/write -s . -p 2323
```

## How it works
The user's input is simply written to a flat file in a given directory. To provide web access, a web server (sold separately) serves all files in this directory as `plain/text`. That's it!

## License
This project is licensed under the MIT open source license.
