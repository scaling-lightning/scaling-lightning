---
sidebar_position: 5
---

# Running on github actions

As a scaling lightning network can be setup and interacted with from code it lends itself to being usable for end to end or functional tests running on Github actions (or CI/CD tool of choice).

The following is an example github actions file that utilises minikube to run the cluster:

```yaml
name: Scaling Lightning On MiniKube using GitHub Actions

"on":
  push:
    branches:
      - main

jobs:
  run-sl:
    runs-on: ubuntu-latest
    steps:
      - name: start minikube
        id: minikube
        uses: medyagh/setup-minikube@master
      - name: Setup helmfile
        uses: mamezou-tech/setup-helmfile@v1.2.0
      - name: Install traefik CRDs
        run: |
          helm repo add traefik https://traefik.github.io/charts
          helm repo update
          helm install traefik traefik/traefik -n sl-traefik --create-namespace -f https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/v0.0.33/charts/traefik-values.yml
      - name: Start minikube's loadbalancer tunnel
        run: minikube tunnel &> /dev/null &
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: ">=1.21.0"
      - name: Run example test
        run: go test -run ^TestMainExample$ github.com/scaling-lightning/scaling-lightning/examples/go -count=1 -v -timeout=15m
```

The [example test](https://github.com/scaling-lightning/scaling-lightning/blob/main/examples/go/example_test.go) can be found in our repo under examples/go.
