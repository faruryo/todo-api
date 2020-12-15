
## Develop

### code generate

```
go generate ./...
```

### Local testing

```
kubectl apply -f ./kubernetes
skaffold dev --port-forward
```

http://localhost:8080/playground

## 参考

- [Build a GraphQL API in Golang with MySQL and GORM using Gqlgen | SoberKoder](https://www.soberkoder.com/go-graphql-api-mysql-gorm/)
- [gqlgen + Gorm でUint型の場合エラーになる - Qiita](https://qiita.com/3104k/items/caf17633d4926aee8a84)
- [Golang テスト sqlmock | 実務のGo](https://www.go-lang-programming.com/doc/test/sqlmock)
- [JordanKnott/taskcafe: An open source project management tool with Kanban boards](https://github.com/JordanKnott/taskcafe)
- [alexzimmer96/gqlgen-example: Example of how to structure a GraphQL-Application using Go and gqlgen](https://github.com/alexzimmer96/gqlgen-example)
