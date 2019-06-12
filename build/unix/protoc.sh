cd ../..
protoc --go_out=. protobuf/transaction.proto
protoc --go_out=. protobuf/account.proto
protoc --go_out=. protobuf/contract.proto
