module github.com/AVOXI/ari/v5

replace github.com/CyCoreSystems/ari/v5 => ./

require (
	github.com/CyCoreSystems/ari/v5 v5.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.2
	github.com/inconshreveable/log15 v0.0.0-20171019012758-0decfc6c20d9
	github.com/oklog/ulid v1.3.1
	github.com/rotisserie/eris v0.4.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	google.golang.org/protobuf v1.26.0
)

go 1.13
