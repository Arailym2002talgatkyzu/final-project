genPOST:
	  protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative post_db/postpb/*.proto
cleanGeneralExample:
	rm greet/greetpb/*.go

genAUTH:
 protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative authorization/authpb/*.proto

