module github.com/drand/kyber

go 1.17

require (
	github.com/drand/drand v1.2.5
	github.com/drand/kyber-bls12381 v0.2.2
	github.com/jonboulle/clockwork v0.3.0
	github.com/stretchr/testify v1.7.2
	go.dedis.ch/fixbuf v1.0.3
	go.dedis.ch/protobuf v1.0.11
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200608115520-7c474a2e3482 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/drand/kyber-bls12381 => ../kyber-bls12381
