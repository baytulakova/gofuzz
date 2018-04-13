#gofuzz

If you don’t have go-fuzz already installed, do that now

$ go get github.com/dvyukov/go-fuzz/go-fuzz
$ go get github.com/dvyukov/go-fuzz/go-fuzz-build

Then 

$GOPATH/bin/go-fuzz-build github.com/baytulakova/gofuzz
$GOPATH/bin/go-fuzz -bin=./arp-fuzz.zip -workdir=workdir
