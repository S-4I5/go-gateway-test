module gateway-service

go 1.21rc3

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	example.com/discovery-service v1.0.0
)

replace example.com/discovery-service => ../discovery-service

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231030173426-d783a09b4405 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
