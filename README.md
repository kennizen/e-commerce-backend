# E-Commerce Backend API's (Dummy)

This repository represents an online e-commerce platform backend system written in Go that is solely made to make frontend mockups and designs for generic use cases.  
This repository does not contain complex APIs but simple CRUD APIs as a base.

## Backend Features
- Pre-seeded with 190+ product data 🛒.
- User authentication with JWT 🔐.
- Custom middleware for logging user requests 📰.
- Rotating refresh token strategy 🔑.
- Payload validation ✅.
- Monitoring and logging with Dozzle 🖥️.
- API docs with Swagger 🗎.
- Live reloading with Go AIR 🔃.
- Containerized with Docker 🐳.
- Postgres as the database 🐘.
- PgAdmin (for folks wanting to view the database design and the data in the browser) 👮.

## Steps to get going

Make sure to have Docker installed on your system. If not, you can follow this [link](https://www.docker.com/) to get it installed.
With that done, clone the repository into your local machine.

```bash
git clone https://github.com/kennizen/e-commerce-backend.git
```
The following are the ports used by the application so make sure nothing is already bound to them.

#### Ports - 8080, 5433, 6600, 6601

Once the repository is cloned, start by typing these commands into your terminal. Make sure you are at the root of the project.

```bash
docker compose build
```
```bash
docker compose up
```
Wait for docker to finish setup and you should see after all the logs that. 

🚀🚀🚀 Server started on host port 8080.

If all the above steps are successful then the database is created, product data is seeded into the database, all the necessary migrations are done and you are good to go.

#### For API docs type the following into the browser's url.

```
http://localhost:8080/swagger/index.html#/
```

#### For opening Dozzle visit

```
http://localhost:6600/dozzle/
```

#### For using the PgAdmin visit the following URL

```
http://localhost:6601/login?next=/
```
You will be greeted with a login screen with username and password both of which can be found in the .env file of the project as PGADMIN_DEFAULT_EMAIL and PGADMIN_DEFAULT_PASSWORD.

Once you enter PgAdmin click on `Add New Server` then under `General` provide a name for your connection. Then move to `Connection` tab there `Host name/address` will be `postgres` as the database is running inside docker so this is the container's name. `Maintenance database` will be `ecommerce`. `Username` and `Password` will be `postgres`. Save the connection and you should be connected to the database inside docker. 

So there you have it. Now you can start creating your beautiful frontend mockups for an e-commerce application.
