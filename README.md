## Table of Contents
- [Commands](#commands)
- [Configuration](#configuration)
- [Installation](#installation)
- [Note](#Note)


## Commands

There are multiple resource that you can manage using the okta-cli

Examples for group resource
```shell
$ okta-cli-client group get --groupId groupId
```
```shell
$ okta-cli-client group create --data '{ "profile": { "description": "test", "name": "Test" }, "type": "OKTA_GROUP"}'
```
```shell
$ okta-cli-client  group replace --groupId groupId --data '{ "profile": { "description": "test", "name": "Test2" }, "type": "OKTA_GROUP"}'
```
```shell
$ okta-cli-client  group delete --groupId groupId
```
```shell
$ okta-cli-client  group lists
```

## Configuration

This library looks for configuration in the following sources:

0. An `okta.yaml` file in a `.okta` folder in the current user's home directory
   (`~/.okta/okta.yaml` or `%userprofile\.okta\okta.yaml`)
0. A `.okta.yaml` file in the application or project's root directory
0. Environment variables

### YAML configuration

When you use an API Token instead of OAuth 2.0 the full YAML configuration
looks like:

```yaml
okta:
  client:
    connectionTimeout: 30 # seconds
    orgUrl: "https://{yourOktaDomain}"
    proxy:
      port: null
      host: null
      username: null
      password: null
    token: {apiToken}
```

When you use OAuth 2.0 the full YAML configuration looks like:

```yaml
okta:
  client:
    connectionTimeout: 30 # seconds
    orgUrl: "https://{yourOktaDomain}"
    proxy:
      port: null
      host: null
      username: null
      password: null
    authorizationMode: "PrivateKey"
    clientId: "{yourClientId}"
    scopes:
      - scope.1
      - scope.2
    privateKey: |
        -----BEGIN RSA PRIVATE KEY-----
        MIIEogIBAAKCAQEAl4F5CrP6Wu2kKwH1Z+CNBdo0iteHhVRIXeHdeoqIB1iXvuv4
        THQdM5PIlot6XmeV1KUKuzw2ewDeb5zcasA4QHPcSVh2+KzbttPQ+RUXCUAr5t+r
        0r6gBc5Dy1IPjCFsqsPJXFwqe3RzUb...
        -----END RSA PRIVATE KEY-----
    privateKeyId: "{JWK key id (kid}" # needed if Okta service application has more then a single JWK registered
    requestTimeout: 0 # seconds
    rateLimit:
      maxRetries: 4
```

### Environment variables

Each one of the configuration values above can be turned into an environment
variable name with the `_` (underscore) character:

* `OKTA_CLIENT_CONNECTIONTIMEOUT`
* `OKTA_CLIENT_TOKEN`
* and so on

## Installation


## Note
In its current form, we only support stdout as output