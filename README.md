# toolkit
A small, opinionated, modular toolkit for writing http services in Go

Go 1.16 or above is required because of the usage of embed.FS.

## Packages
Each package is named so it doesn't conflict with the stdlib packages with
similar features. It should be possible to use only these packages to build
services as they wrap libraries to expose a more ergonomic interface.

- db - Mostly a wrapper around jmoiron/sqlx and golang-migrate. It only exposes
  interfaces which take a `context.Context` to ensure people are thinking about
  timeouts.
- web - A wrapper around a number of packages to provide a convenient
  one-stop-shop for interfacing with http services. It wraps net/http,
  encoding/json, and formam. You must choose a router to use. go-chi/chi
  is the author's current recommended choice.

Package ideas:

- Logging - I don't know what this would be called or what it would look like
  yet, so it's on the backburner for now.

Feature ideas:
- db
    - Migrations. Either using golang-migrate, or a custom solution.
- web
    - Routing. Provide a useful router which works well with the included
      handlers.
