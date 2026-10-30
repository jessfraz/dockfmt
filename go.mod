module github.com/jessfraz/dockfmt

require (
	github.com/genuinetools/pkg v0.0.0-20181022210355-2fcf164d37cb
	github.com/moby/buildkit v0.12.5
	github.com/sirupsen/logrus v1.9.1
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/containerd/typeurl/v2 v2.1.1 // indirect
	github.com/docker/docker v24.0.0-rc.2.0.20230718135204-8e51b8b59cb8+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gotest.tools/v3 v3.4.0 // indirect
)

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/containerd/containerd v1.3.0-0.20190507210959-7c1e88399ec0 => github.com/containerd/containerd v1.3.0-rc.1.0.20190916145203-86442dfbb9c7

replace github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v1.4.2-0.20190916154449-92cc603036dd

replace golang.org/x/crypto v0.0.0-20190129210102-0709b304e793 => golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7

go 1.20
