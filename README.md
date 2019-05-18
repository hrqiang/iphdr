# To recreate golang library [wrong endianness in IP header on darwin-amd64.](https://github.com/golang/go/issues/32118#issuecomment-493580817)
* RAW socket(run with sudo) ICMP message send to Google.
* IPHDR byte 3,4 on 64bit OSX and Linux has different order. 
* No reply on OSX, linux works fine. 
