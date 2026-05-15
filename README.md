# Base App Client

This is a template web application and set of APIs written in go, htmx, tailwind, _hyperscript and postgresql. This project is elite for how boring it is. Dependencies are minimised, all communications between the application and database are done via stored procedures and functions and there is an aversion to technology that is ubiquitous but inappropriate for a straighforward web application and set of APIs.

If you desire to be exposed to the latest, most fashionable technology and become a sophistcated developer then any investment in time you make here will be wasted. If you want to learn how to build a boring, solid web application that will not increase your ability to gain employment but will run without much maintenance then you should stay.

This repository contains the front end for (1) a template web application and that uses go, htmx, _hyperscript and postgresql and (2) a set of APIs. The postgresql scripts are are in the repository [https://github.com/andrewah64/base-app-db/](https://github.com/andrewah64/base-app-db/).

## Philosophy

- The database is the source of truth.
- Any business rule that can be enforced by the database, will be enforced by the database.
- Database functions and stored procedure calls are the only way to interact with the database.
- The technology stack will be as spartan as possible.

## Functionality

- Authentication
    - API
        - User-level API key
    - Web application
        - Username/password + TOTP
        - OIDC
        - Passleys
        - SAML2
- Logging
    - Configurable at runtime at the levels of
        - Endpoint
        - User + endpoint
- Security
    - Role-base access control enforced in the application and database
- Session management
    - Kill HTTP sessions

## Getting started

- All code snippets that follow were tested on Ubuntu 26.04.
- <root> represents the path to the root directory of the repo.
- All commands are to be executed from <root>.

### Clone the repos

```
git clone https://github.com/andrewah64/base-app-client.git
git clone https://github.com/andrewah64/base-app-db.git
```

### Setup database connectivity

Choose the method that will be used to get database credentials.

- AWS secrets manager
    - TBD
- Username / password
    - TBD
- systemd
    - Create a systemd unit file override directory:
        - ```sudo mkdir -p /etc/systemd/system/base-app.service.d/```
    - Generate an encrypted password file
        - ```echo "[Service]" sudo systemd-ask-password -n | sudo systemd-creds encrypt --name=postgres-password -p - - | sudo tee /etc/systemd/system/base-app.service.d/postgresql.conf > /dev/null```
    - Tailor ```<root>/systemd/base-app.service```
        - ```sudo cp <root>/base-app-api.service /etc/systemd/system```
        - ```sudo cp <root>/base-app-web.service /etc/systemd/system```
    - Enable the unit file:
        - ```sudo systemctl daemon-reload```
        - ```sudo systemctl enable --now base-app-api.service```
        - ```sudo systemctl enable --now base-app-web.service```
        - ```sudo systemctl status base-app-api.service```
        - ```sudo systemctl status base-app-web.service```
    - Generate deployable executable files
        - ```go build -ldflags "-s -w" -trimpath -o base-app-api ./cmd/api```
        - ```go build -ldflags "-s -w" -trimpath -o base-app-web ./cmd/web```
    - Move the executable files and pem files to the folder specified in the unit file's 'WorkingDirectory'
        - ```mv base-app-api <WorkingDirectory>```
        - ```mv base-app-web <WorkingDirectory>```
        - ```mv *.pem <WorkingDirectory>```
    - Start the services
        - ```sudo systemctl start base-app-api.service```
        - ```sudo systemctl start base-app-web.service```
