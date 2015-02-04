Write.as
========
[![Build Status](https://travis-ci.org/writeas/writeas-telnet.svg)](https://travis-ci.org/writeas/writeas-telnet)

This is a simple telnet-based interface for publishing text. Users connect and paste / type what they want to publish. Upon indicating that they're finished, a link is generated to access their new post on the web.

## Try it
```
telnet nerds.write.as
```

## How it works
The user's input is simply written to a flat file in a given directory. To provide web access, a web server (sold separately) serves all files in this directory as `plain/text`. That's it!

## License
This project is licensed under the MIT open source license.
