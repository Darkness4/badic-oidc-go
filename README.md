# Basic OIDC example in Go

## Usage

Fill a `.env` with the OIDC issuer URL, client secret and client ID:

```shell
CLIENT_SECRET=GOCSPX-0123456789abcdefghijklmnopqr
CLIENT_ID=123456789012-0123456789abcdefghijklmnopqrstuv.apps.googleusercontent.com
OIDC_ISSUER=https://accounts.google.com
```

Launch the server:

```shell
go run main.go
```

Go to the login page: [http://localhost:3000/login](http://localhost:3000/login).

## With Dex

Configure the dex by editing the `dex/config.yaml`. Remove or add providers.

Run the dex server by running the `dex/run.sh` script.

Launch the server:

```shell
go run main.go
```

Go to the login page: [http://localhost:3000/login](http://localhost:3000/login).

## With 389ds and dex

Run the 389ds server by running the `run.sh` script.

Initialize the server:

```shell
docker exec -it 389ds bash

dsconf localhost backend create --suffix dc=example,dc=com --be-name example_backend # Create a backend (a backend is literally a database)
dsidm localhost initialise # Creates examples
# Create a user
dsidm -b "dc=example,dc=com" localhost user create \
  --uid example-user \
  --cn example-user \
  --displayName example-user \
  --homeDirectory "/dev/shm" \
  --uidNumber -1 \
  --gidNumber -1
# Set a user password:
dsidm -b "dc=example,dc=com" localhost user modify \
  example-user add:userPassword:"...."
dsidm -b "dc=example,dc=com" localhost user modify \
  example-user add:mail:example-user@example.com
```

Edit the dex configuration to include LDAP:

```yaml
#config.yaml
#...
connectors:
  - type: ldap
    id: ldap
    name: LDAP
    config:
      host: <your-host-IP>:3389 # EDIT THIS. If you use docker-compose with root, you can set a domain name.
      insecureNoSSL: true
      userSearch:
        baseDN: ou=people,dc=example,dc=com
        username: uid
        idAttr: uid
        emailAttr: mail
        nameAttr: cn
        preferredUsernameAttr: uid
      groupSearch:
        baseDN: ou=groups,dc=example,dc=com
        userMatchers:
          - userAttr: uid
            groupAttr: member
        nameAttr: cn
```

Run the dex server by running the `dex/run.sh` script.

Launch the server:

```shell
go run main.go
```

Go to the login page: [http://localhost:3000/login](http://localhost:3000/login).
