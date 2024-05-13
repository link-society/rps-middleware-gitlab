# RPS Middleware - Gitlab

When running Gitlab behind a reverse proxy, it can be challenging to handle its
misbehaviour regarding the HTTP RFC. This project act as a middleware between
your reverse proxy and the proxied Gitlab instance.

## :necktie: Rules

### Fix API Projects

The Gitlab API for projects expect the repository name (with the namespace) to
be URL encoded:

```
/api/v4/projects/namespace%2Frepository
```

However, most reverse proxies will (rightfully) decode those URLs into:

```
/api/v4/projects/namespace/repository
```

Which will result in a `404 - Not Found` HTTP error.

As such, this middleware will manually re-encode the last part of the URL before
proxying the request to the Gitlab.

## :hammer: Build

Requirements:

 - Go 1.22+

Run:

```
$ make
```

## :memo: Usage

Run (adjust the `REMOTE_URL` variable):

```
$ export REMOTE_URL="https://gitlab.com"
$ ./rps-middleware-gitlab
```

Then, point your reverse proxy to `<YOUR-IP>:8080`.

Or run the Docker image:

```
$ docker run \
    -e REMOTE_URL="https://gitlab.com" \
    -p 8080:8080 \
    linksociety/rps-middleware-gitlab:latest
```

> :construction: Docker image has not been published yet

## :balance_scale: License

This software is distributed in the public domain, see the
[LICENSE.txt](./LICENSE.txt) document for more information.
