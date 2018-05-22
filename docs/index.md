# Rule Template Helper Functions

`templr` has a number of string and networking helper functions that can be used in the rule templates. An example config can be found [here](https://github.com/gesquive/templr/blob/master/pkg/rules.example.yml).

## slice
`slice` returns the given arguments as an iterable list:

```
slice "8.8.8.8" "8.8.4.4"
```
The above returns `['8.8.8.8', '8.8.4.4']`

## list
`list` returns a comma delimeted list of the given array:
```
list "8.8.8.8" "8.8.4.4"
```
The above returns "`8.8.8.8, 8.8.4.4`"

## rpad
`rpad` pads spaces onto the right side of the given string
```
rpad 10 "hello"
```
The above produces "`hello`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"

## ipfmt
`ipfmt` pads an IP address string based on the type of IP address it is
```
ipfmt "8.8.8.8"
```
The above produces "`8.8.8.8`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"

## lookupHosts
`lookupHosts` returns a list of [HostInfo](https://godoc.org/github.com/gesquive/templr/engine#HostInfo) objects
```
lookupHosts $hostList
```

## lookupIPv4Host
`lookupIPv4Host` returns a list of the given host's IPv4 addresses
```
lookupIPv4Host "google-public-dns-a.google.com"
```
The above returns "`8.8.8.8`"

## lookupIPv6Host
`lookupIPv6Host` returns a list of the given host's IPv6 addresses
```
lookupIPv4Host "google-public-dns-a.google.com"
```
The above returns "`2001:4860:4860::8888`"

## isValidIPv4
`isValidIPv4` returns true if the given address is a valid IPv4 address or IPv4 CIDR range
```
isValidIPv4 "8.8.8.8"
```
The above returns `true`

## isValidIPv6
`isValidIPv6` returns true if the given address is a valid IPv6 address or IPv6 CIDR range
```
isValidIPv6 "2001:4860:4860::8888"
```
The above returns `true`

## isValidIPv4Addr
`isValidIPv4Addr` returns true if the given address is a valid IPv4 address
```
isValidIPv4Addr "8.8.8.8"
```
The above returns `true`

## isValidIPv6Addr
`isValidIPv6Addr` returns true if the given address is a valid IPv6 address
```
isValidIPv6Addr "2001:4860:4860::8888"
```
The above returns `true`

## isValidIPv4CIDR
`isValidIPv4CIDR` returns true if the given address is a valid IPv4 CIDR range
```
isValidIPv4CIDR "10.0.0.0/8"
```
The above returns `true`

## isValidIPv6CIDR
`isValidIPv6CIDR` returns true if the given address is a valid IPv6 CIDR range
```
isValidIPv6CIDR "2001:db8::/32"
```
The above returns `true`

