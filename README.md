# transip-dns-acmetool
Provide DNS hook for the TransIP API, can be used by acmetool

## How to use
Replace the Account Name in the (file)[main.go#L36]
Place a file `transip-priv.key` next to this binary (or see (main.go)[main.go#L37] 

Then move the compiled binary to the directory `/etc/acme/hooks` and enjoy wildcard certificates!




