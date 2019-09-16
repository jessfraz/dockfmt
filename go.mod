module github.com/jessfraz/dockfmt

require (
	github.com/genuinetools/pkg v0.0.0-20181022210355-2fcf164d37cb
	github.com/moby/buildkit v0.6.1
	github.com/sirupsen/logrus v1.4.3-0.20190807103436-de736cf91b92
)

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/containerd/containerd v1.3.0-0.20190507210959-7c1e88399ec0 => github.com/containerd/containerd v1.3.0-rc.1.0.20190916145203-86442dfbb9c7

replace github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v1.4.2-0.20190916154449-92cc603036dd

replace golang.org/x/crypto v0.0.0-20190129210102-0709b304e793 => golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7

go 1.13
