module ValidateAWSAppConnector

go 1.21.3

replace (
    appconnections => ../../
    cowlibrary => ../../../../../src/cowlibrary
)

require (
    appconnections v0.0.0-00010101000000-000000000000
    cowlibrary v0.0.0-00010101000000-000000000000
    github.com/google/uuid v1.6.0
    gopkg.in/yaml.v3 v3.0.1
)

require (
    github.com/aws/aws-sdk-go v1.43.31 // indirect
    github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
    github.com/charmbracelet/bubbletea v0.24.2 // indirect
    github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
    github.com/dmnlk/stringUtils v0.0.0-20150214151148-aa88c62978f5 // indirect
    github.com/dustin/go-humanize v1.0.1 // indirect
    github.com/gabriel-vasile/mimetype v1.4.2 // indirect
    github.com/go-playground/locales v0.14.1 // indirect
    github.com/go-playground/universal-translator v0.18.1 // indirect
    github.com/go-playground/validator/v10 v10.14.0 // indirect
    github.com/go-resty/resty/v2 v2.7.0 // indirect
    github.com/iancoleman/strcase v0.2.0 // indirect
    github.com/jmespath/go-jmespath v0.4.0 // indirect
    github.com/json-iterator/go v1.1.12 // indirect
    github.com/jwalton/go-supportscolor v1.1.0 // indirect
    github.com/klauspost/compress v1.16.0 // indirect
    github.com/klauspost/cpuid/v2 v2.2.4 // indirect
    github.com/leodido/go-urn v1.2.4 // indirect
    github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
    github.com/mattn/go-isatty v0.0.18 // indirect
    github.com/mattn/go-localereader v0.0.1 // indirect
    github.com/mattn/go-runewidth v0.0.14 // indirect
    github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2 // indirect
    github.com/minio/md5-simd v1.1.2 // indirect
    github.com/minio/minio-go/v7 v7.0.50 // indirect
    github.com/minio/sha256-simd v1.0.0 // indirect
    github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
    github.com/modern-go/reflect2 v1.0.2 // indirect
    github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
    github.com/muesli/cancelreader v0.2.2 // indirect
    github.com/muesli/reflow v0.3.0 // indirect
    github.com/muesli/termenv v0.15.1 // indirect
    github.com/rivo/uniseg v0.2.0 // indirect
    github.com/rs/xid v1.4.0 // indirect
    github.com/sirupsen/logrus v1.9.0 // indirect
    golang.org/x/crypto v0.9.0 // indirect
    golang.org/x/net v0.10.0 // indirect
    golang.org/x/sync v0.1.0 // indirect
    golang.org/x/sys v0.8.0 // indirect
    golang.org/x/term v0.8.0 // indirect
    golang.org/x/text v0.9.0 // indirect
    google.golang.org/protobuf v1.30.0 // indirect
    gopkg.in/ini.v1 v1.67.0 // indirect
    gopkg.in/yaml.v2 v2.4.0 // indirect
)