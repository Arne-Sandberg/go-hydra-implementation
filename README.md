# go-hydra-implementation

This is a simple implementation of [ORY Hydra](https://github.com/ory/hydra) (OpenID Connect certified OAuth2 Server) in Go. It also explores testing Hydra endpoints.

All dependencies are managed by Go 1.11 native package and dependency management.

## Running

You just need to run a Makefile up command:

```sh
make up
```

Then create a client ID to log in before accessing http://127.0.0.1:5000

```sh
docker exec -it hydra \
    hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id go-hydra-implementation \
    --secret secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,offline \
    --callbacks http://127.0.0.1:5000/callback
```


## Testing

Tests automatically create and delete a client ID. Just run:

```sh
make test
```