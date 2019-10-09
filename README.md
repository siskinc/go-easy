# go-easy
utils for easy coding in golang

## generate

### mongodb
[example](./test/test_generate_mongodb.go)  

command run:

```shell script
export GO111MODULE=on
go build . && go generate test/test_generate_mongodb.go
```

#### functions
1. curd
2. soft delete
3. set timestamp to `create at` field
4. set timestamp to `update at` field
5. set timestamp to `delete at` field