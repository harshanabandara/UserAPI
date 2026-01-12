### Building and running the UserAPI server
#### PostgreSQL instance
To run the PostgreSQL with docker compose follow the steps given below. 
```bash
cd infra/postgres
docker compose up -d
```
This will run a PostgreSQL db instance. If the password needs to be changed, change the value in the 
`password.txt`

When you're ready, start your application by running:
```
docker compose up --build
```
### build the project locally

#### Database query generation
to generate/update database query helpers, run the following commands.
```bash
cd config/database
sqlc generate
```
generated database files will be in `internal/adapters/db/user` directory.

#### SWAG API Documentation
To generate/update api documentation, run the following command.
Documentation will be generated in .docs directory.
```
swag init -g ./cmd/api-server/main.go -d  . --parseInternal
```
#### Build the project
Once the above files are updated, run the below command. 
```bash
cd cmd/api-server
go build .
./api-server
```

#### building a docker image of the server

To build the docker image locally, run the below command at the project root. 
- with multi-arch support:
```
docker buildx build --platform linux/amd64,linux/arm64 -t userapi:latest .
```
- without multi-arch support:
```bash
docker build -t userapi .
```

Following environment variables can be passed into the docker container to override initial values.

| Key         | default   | note                                     |
|-------------|-----------|------------------------------------------|
| PG_HOST     | localhost | the hostname of the db server            |
| PG_PORT     | 5432      | the listening port of the datbase server |
| PG_USER     | postgres  | username for the database                |
| PG_PASSWORD | yaalalabs | password of the database                 |
| PG_DATABASE | userapi   | database name                            |                             
| PG_SSLMODE  | disable   | ssl mode                                 |

if you want to push as you build, run below command. 
```bash
docker buildx build --platform linux/amd64,linux/arm64 -t $dockerhub_username/userapi:latest --push .
```
#### check for linting issues
run below command in the root. 
```
golangci-lint run ./...
```

#### Docker compose deployment
- Navigate to the `infra/deploy` directory and run `docker compose up -d`

#### Kubernetes deployment with microk8s
1. Checkout the GitHub files to a directory.
    ```bash
    mkdir k8s
    cd k8s
    curl -o database.yaml https://raw.githubusercontent.com/harshanabandara/UserAPI/main/infra/deploy/k8s/database.yaml
    curl -o server.yaml https://raw.githubusercontent.com/harshanabandara/UserAPI/main/infra/deploy/k8s/server.yaml
    curl -o kustomization.yaml https://raw.githubusercontent.com/harshanabandara/UserAPI/main/infra/deploy/k8s/kustomization.yaml
    ```
2. change the kustomization.yaml and add the password. 
3. run `microk8s kubectl apply -k ./` to deply the services.
4. access the api usinsg [IP]:30001/users