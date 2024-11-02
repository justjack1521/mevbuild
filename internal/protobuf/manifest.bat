cd %GOPATH%/src
protoc --proto_path=github.com/justjack1521/mevpatch/internal/protobuf --go_out=github.com/justjack1521/mevpatch/internal/manifest --go_opt=paths=source_relative manifest.proto