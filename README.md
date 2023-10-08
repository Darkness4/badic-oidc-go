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
