module github.com/dxvgef/tsing-center

go 1.14

replace (
	github.com/coreos/bbolt => github.com/coreos/bbolt v1.3.3
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.22+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/dxvgef/filter v1.8.1
	github.com/dxvgef/tsing v1.2.2
	github.com/dxvgef/tsing-gateway v0.0.0-20200612110821-3967e76cc753 // indirect
	github.com/mailru/easyjson v0.7.1
	github.com/rs/zerolog v1.19.0
	go.uber.org/atomic v1.6.0 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c
)
