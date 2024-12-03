okta-cli-client user writeUserToFile -userId abc123

user is created in ~/.okta/orgid/user@host.json

okta-cli-client envsync pushuser --userdata ~/.okta/dev-000000.okta.com/users/user@host.json
