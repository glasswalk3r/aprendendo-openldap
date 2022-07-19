package main

import (
  "fmt"
  "flag"
)

const defaultDNSDomain string = "foobar.org"
const defaultBaseDN string = "dc=foobar,dc=org"

func main() {
    var dnsDomain string
    var baseDN string

    flag.StringVar(&dnsDomain, "d", defaultDNSDomain, fmt.Sprintf("Specify the DNS domain to use, default to %s", defaultDNSDomain))
    flag.StringVar(&baseDN, "b", defaultBaseDN, fmt.Sprintf("Specify the base DN, default to %s", defaultBaseDN))

    flag.Parse()

    fmt.Println(dnsDomain, baseDN)
}
