## Maden

**Maden** is a minimal, lightweight container orchestration tool. It can be used for basic local development. The architecture closely mirrors that of Kubernetes, with:
- pods able to run multiple (Docker) containers; they support pod replicas, affinities/anti-affinities, tolerations, restart policies
- deployments and services; they can be configured through yaml manifests as usual
- schedulers determining how to schedule pods based on available resources (only virtual for now), affinities etc.
- controllers ensuring the state of the system reflects the defined configuration
- an etcd data source storing pods, nodes etc.
- an API server allowing interaction with the Maden resources
- a CLI tool to interact with the API server

### How to use
Maden will be packaged soon. For now, you can use it by following these steps:
1. Ensure you have golang and Docker installed and fetch the repository.
2. Run `docker build -t maden:latest .` and `docker-compose up` to start the server.
3. Run `cd cmd\madencli` and `go build -o madencli.exe` to build the CLI tool.
4. Now you can interact with Maden via commands, for example:
`./madencli.exe apply -f \path-to-your-roout\example_deployments\example_deployment.yaml`
This applies the example deployment from the example_deployments directory. Run `./madencli.exe -h` to see all available commands.

### Status
In mid stages of development.

### Contributing
All contributions are warmly welcomed. Head over to [CONTRIBUTING.md](https://github.com/TudorOrban/Maden/blob/main/CONTRIBUTING.md) for details.