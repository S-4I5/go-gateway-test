module test-service

go 1.21rc3

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/render v1.0.3
	github.com/go-playground/validator/v10 v10.15.5
	github.com/ilyakaznacheev/cleanenv v1.5.0
	example.com/discovery-service v1.0.0
)

replace example.com/discovery-service => ../discovery-service

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
