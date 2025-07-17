module github.com/agentic-layer/model-router-krakend

go 1.24

toolchain go1.24.4

require (
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// must match the latest version in https://github.com/devopsfaith/krakend-ce/blob/v2.10.1/go.sum
replace golang.org/x/sys => golang.org/x/sys v0.31.0

replace go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.33.0

replace go.opentelemetry.io/auto/sdk => go.opentelemetry.io/auto/sdk v1.1.0

replace go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0

replace go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.33.0

replace go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.33.0
