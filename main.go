package main

import (
	"maden/pkg/apiserver"
	"maden/pkg/etcd"
)

func main() {
	etcd.InitEtcd()
	apiserver.InitAPIServer()
}
