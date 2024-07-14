# Obscurio Proxy
A postgres proxy to remove, replace and obfuscate data from queries.

## Repository Usage
- This project's source code is stored in a private repo. 
- This repo is purely as a holder for the docker repository.
- We intend to open source this project in the future, likely under an AGPL license.

## Example Usage
### Free Usage
```bash
docker run -d -p 5432:5432 \
  -v ./config.yaml:/config.yaml \
  -e OBSCURIO_CONFIG_FILEPATH=/config.yaml \
  -e DATABASE_URI=postgres://user:password@host:port/dbname \
  ghcr.io/prusiksoftware/obscurio:latest
```
with an example `config.yaml`:
```yaml
log_level: debug

database:
  postgres_version: 16.3
  uri_env: DATABASE_URI

clients:
  - name: client 1
    username_env: USERNAME1
    password_env: PASSWORD1
    filters:

      - table: customer
        column: email
        function: hide column

      - table: customer
        column: country
        function: replace
        value: "****"

  - name: second_profile
    username_env: USERNAME2
    password_env: PASSWORD2
    filters:

      - table: customer
        column: country
        function: hide row
        value: Canada
```

### Paid Usage
If you wish to dynamically manage the configuration, and have access to the admin panel, you can use the paid version. (free 30 day trials)
```bash
docker run -d -p 5432:5432 \
  -e OBSCURIO_API_KEY=example_api_key \
  -e DATABASE_URI=postgres://user:password@host:port/dbname \
  ghcr.io/prusiksoftware/obscurio:latest
```

## Configuration
The configuration file is a yaml file with the structure defined at [docs.obscurio.io/configuration](docs.obscurio.io/configuration).

## API
The API is defined at [docs.obscurio.io/api](docs.obscurio.io/api).