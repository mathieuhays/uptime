# uptime

![badge](https://github.com/mathieuhays/uptime/actions/workflows/build.yml/badge.svg)

Service tracking websites' uptime. WIP

## TODO

- [ ] Add integration tests using `run()` in the server cmd pkg
- [ ] Set cookie on login, refresh session expiration date?
- [ ] template: login
- [ ] template: sign up
- [ ] template: dashboard

## Development

### commands

`make install_deps` installs the go package needed for development (i.e sqlc and goose)

`make update_sql` updates the `internal/database` pkg. Should be run when files in sql are updated.

`make up` migration up -- this is automatically run when the server instance boots up

`make down` migration down

`make jwt_secret` generate jwt secret to be pasted in `.env`

## Interesting Read

- [How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/#maker-funcs-return-the-handler)