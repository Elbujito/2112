# Project 2112
# Getting Started

1. Clone this repository `git clone git@github.com:elbujito/2112.git`
2. Run `cd 2112`
3. Run `go get`
4. Run `go run . db migrate`
5. Run `go run . db seed`
6. Run `go run . start` to start the server, you should see the following:
```
⇨ http server started on [::]:8081
⇨ http server started on [::]:8080
⇨ http server started on [::]:8079
```
7. List available routes using `go run . info protected-api-routes` and use your favourite API client to test. or use the following to get started and make sure you're up and running.
```bash
curl -H "Accept: application/json" http://127.0.0.1:8081/health/alive
curl -H "Accept: application/json" http://127.0.0.1:8081/health/ready
```

> Recommended: run `go run .` and explore all available options, it should be straightforward.

For more details on running and using the service, scroll down to "[Operations](#operations)" section. 


   ```bash
   git clone git@github.com:elbujito/2112.git
   cd 2112
