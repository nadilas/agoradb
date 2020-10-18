# Code first example of agoradb

This packge is an example of a code first approach.

The schema for the database is defined first in the schema.go source file.

### Before running go generate ./...

1. Make sure at least one database node is running
2. Create `test_users` database in the database
3. It is recommended to log in with the cli before running go generate
	```bash
	agoradb login localhost:5750 [--user u1 --password p1]
	```
4. Run `go generate ./...` which will execute two things:
	First it will migrate the database and deploy the schema
	Second it will download the protobuf definition and generate the client code
## Trying out