# Obscurio Proxy
## `a postgresql data hygiene proxy`

## What is Obscurio Proxy?
Obscurio Proxy is a simple proxy server that sits between your application and your database. It is designed to obfuscate sensitive data in your database before it reaches your application (or any other tools with access to your database).

Obscurio Proxy is designed to help you comply with data privacy regulations such as GDPR, HIPAA, and others. It is also useful for protecting sensitive data from unauthorized access by developers, contractors, or other third parties.

## Testing Obscurio Proxy
if you want to evaluate Obscurio Proxy, you can use the following steps to set up a test environment:
```bash
git clone git://github.com/prusiksoftware/monorepo.git
cd monorepo/obscurio
docker-compose up
```
now you will have two postgresql databases running on your machine:
- `localhost:5433` - the original database
- `localhost:5432` - the proxy database

- if you edit the `example-config.yaml` file, you can configure the proxy to obfuscate the data in the `localhost:5432` database before it reaches your application.



