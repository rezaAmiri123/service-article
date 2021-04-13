module github.com/rezaAmiri123/service-article

go 1.14

replace github.com/rezaAmiri123/service-user => ../service-user

require (
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/jinzhu/gorm v1.9.16
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rezaAmiri123/service-user v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.uber.org/zap v1.16.0
	google.golang.org/genproto v0.0.0-20210406143921-e86de6bf7a46
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)
