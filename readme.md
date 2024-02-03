# VOTE ITEMS
This repository hosts the Vote Items application.

You can read about the design of this application in my [Medium article](https://medium.com/@kritwis/golang-clean-architecture-with-demo-e0938e5be02b).

This application serves as a refined project demonstrating Clean Architecture design principles. It includes features such as a reverse proxy, testing, structured code organization, and Swagger documentation.

![App Overview](./application_overview.png)

## Setup

### Traefik Setup

1. Open your hosts file:
   - On Linux or MacOS, the file is `/etc/hosts`.
   - On Windows, the file is `c:\Windows\System32\Drivers\etc\hosts`.

2.  Add the following line to the file: `127.0.0.1 krittawatcode.test`
This will route any requests for `krittawatcode.test` to your local machine.


### Postman Setup
1. Import the VOTE-ITEMS.postman_collection.json file into Postman.
2. Set up your environment variables in Postman as needed.
3. Use the imported collection to test the API endpoints.

### Startup 
- Run your Docker Compose file with docker-compose up.
- Call api login with default user
```
{Email: "admin@mtl.co.th", Password: "adminPassword"},
{Email: "krittawat@mercy.gg", Password: "userPassword"},
```

# API DOCs
can be found at http://krittawatcode.test/swagger/index.html
![API Docs](./api_docs.png)

## END CREDIT
it's such a fun run to do this!<3\
sadly I cant finish everything I want to do but...\
 hopefully u guys have a mercy on me lol\
... Im still working on it to make it purfect meow~!
