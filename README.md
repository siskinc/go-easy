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
1. [x] curd
2. [x] generate soft delete code
3. [x] generate find by xxx unique field code
4. [x] set timestamp to `create at` field
5. [x] set timestamp to `update at` field
6. [x] set timestamp to `delete at` field
7. [ ] generate migrate code by unique index