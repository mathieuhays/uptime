# uptime

![badge](https://github.com/mathieuhays/uptime/actions/workflows/build.yml/badge.svg)

Service tracking websites' uptime. WIP

## Development

### commands

`make install_deps` installs the go package needed for development (i.e sqlc and goose)

`make update_sql` updates the `internal/database` pkg. Should be run when files in sql are updated.

`make up` migration up -- this is automatically run when the server instance boots up

`make down` migration down

`make jwt_secret` generate jwt secret to be pasted in `.env`