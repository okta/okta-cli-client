# okta-client-cli - The official Okta CLI for the Okta Management API

:warning: The Okta CLI is in beta. We'd love for you to play with it and give us your thoughts, but we don't recommend using it in production applications just yet. We'll be iterating it towards a stable 1.0 release over the next few weeks, based on the feedback we receive. There may be minor interface breaking changes before we stabilize at 1.0. 

The Okta CLI can be used to easily interact with the Okta management API and:

* Create and update users with the [Users API](https://developer.okta.com/docs/api/resources/users)
* Manage groups with the [Groups API](https://developer.okta.com/docs/api/resources/groups)
* Manage applications with the [Apps API](https://developer.okta.com/docs/api/resources/apps)
* Much more!

> Note: In the next few weeks, we'll be working on improving our docs. In the meantime, you can check out the [template file](https://github.com/okta/okta-cli-client/blob/main/template.yaml) for API coverage.

## Table of Contents
- [Release status](#release-status)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage guide](#usage-guide)
- [Notes](#Notes)
- [Need help?](#need-help)

## Release status

This library uses semantic versioning and follows Okta's [library version policy](https://developer.okta.com/code/library-versions/).

| Version | Status                    |
| ------- | ------------------------- |
| 0.x | 🚧 beta |

## Installation

Build the source code locally by executing the following command:

```sh
$ make install
```

## Configuration

The Okta CLI looks for configuration in the following sources:

1. An `okta.yaml` file in a `.okta` folder in the current user's home directory
   (`~/.okta/okta.yaml` or `%userprofile\.okta\okta.yaml`)
1. A `.okta.yaml` file in the application or project's root directory
1. Environment variables

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

## Usage guide

### Manage your Okta resources

#### Get a group by ID

```shell
$ okta-cli-client group get --groupId <GROUP_ID>
```

#### Create a new group

```shell
$ okta-cli-client group create --data '{ "profile": { "description": "test", "name": "Test" }, "type": "OKTA_GROUP"}'
```

#### Update an existing group

```shell
$ okta-cli-client  group replace --groupId <GROUP_ID> --data '{ "profile": { "description": "test", "name": "Test2" }, "type": "OKTA_GROUP"}'
```
#### Delete a group

```shell
$ okta-cli-client  group delete --groupId <GROUP_ID>
```

#### List groups

```shell
$ okta-cli-client group lists
```
#### Assign a group to an application

```sh
$ okta-cli-client applicationGroups assignGroupToApplication --appId <APP_ID> --groupId <GROUP_ID> --data ""
```

#### List groups assigned to an application

```sh
$ okta-cli-client applicationGroups listApplicationGroupAssignments --appId <APP_ID>
```

#### Create a new OIDC application

```sh
$ okta-cli-client application create --data '{"label":"My new services App","name":"oidc_client","signOnMode":"OPENID_CONNECT","settings":{"oauthClient":{"application_type":"service","grant_types":["client_credentials"]}}}'
```

#### Get an application by ID

```sh
$ okta-cli-client application get --appId <APP_ID>
```

#### Deactive and Delete an existing application

```sh
$ okta-cli-client application deactivate --appId <APP_ID>

$ okta-cli-client application delete --appId <APP_ID>
```

#### Create a new user

```sh
$ okta-cli-client user create --data '{"credentials":{"password":{"value":"Hell4W0rld"}},"profile":{"email":"firstname.lastname@gmail.com","firstName":"ExampleFirstName","lastName":"ExampleLastName","login":"firstname.lastname@gmail.com"}}'
```

#### Get a user by ID

```sh
$ okta-cli-client user get --userId <USER_ID>
```

#### Delete a user

```sh
$ okta-cli-client user delete --userId <USER_ID>
```

#### Assign a user to a group

```sh
$ okta-cli-client group assignUserTo --userId <USER_ID> --groupId <GROUP_ID>

```
#### List users assigned to a group

```sh
$ okta-cli-client group listUsers --groupId <GROUP_ID>
```

## Notes
In the Okta CLI current form, we only support `stdout` as output.

## Need help?
 
If you run into problems using the Okta CLI, you can
 
* Ask questions on the [Okta Developer Forums][devforum]
* Post [issues][github-issues] here on GitHub (for code errors)

[devforum]: https://devforum.okta.com/
[github-issues]: https://github.com/okta/okta-client-cli/issues
[github-releases]: https://github.com/okta/okta-client-cli/releases
