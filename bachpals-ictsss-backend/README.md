# BachPals-ICTSSS

## Environment variables
You need to configure all the environment variables available. Look for example file "example.env".
Copy "example.env" into a new file called ".env" and configure your variables.

Double check if the project runs.

## Project set up for development
The main file is in /cmd directory. By running `go run cmd/server.go --help` you can preview all the cli commands.
In order to set up the project correctly you need to do following.

In `hosts` file, point `ikt-stack.internal.uia.no` to `127.0.0.2`. 

1. Run the MongoDB instance, you can check `docker-compose.yaml`, and make sure it is running.
2. Initialize default administrator by running `go run cmd/server.go --init`.
3. In your `hosts` file, make sure your localhost points to the oauth2 domain,`ikt-stack.internal.uia.no`
4. Install and start the frontend project.
5. Run `go run cmd/server.go --serve` to start the HTTP server.

## Commandline
This server provides small cli utility to ease the configuration steps.

Available arguments, with description.
- --reset    Removes all administrators and resets default settings in database.
- --init     Initializes administrators and default settings in database.
- --serve    Starts the http server.
- --help     Shows this page.

### How to use cli arguments?
Run `go run cmd/server.go <argument>`

#### Important
You can only specify one argument when running the command.
In case you provide more than one, only the first one is going to get executed.