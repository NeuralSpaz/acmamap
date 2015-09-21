# acmamap
Map of ACMA data
requires go-bindata
https://github.com/jteeuwen/go-bindata

	$ go get github.com/jteeuwen/go-bindata
	$ go get github.com/NeuralSpaz/acmamap
	$ cd $GOPATH/src/github.com/NeuralSpaz/acmamap
	$ go-bindata www/...
	$ go build .
	$ ./acmamap



![alt tag](https://github.com/NeuralSpaz/acmamap/blob/master/radioLights.png)



Live http://acma.delwp.net
