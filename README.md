# RabbitMQ Modifier
## Requirements for local testing
### Install Docker Desktop
Go to the [Docker Desktop Website](https://www.docker.com/products/docker-desktop/) to install the latest version. 

### Pull the RabbitMQ Image
``docker pull rabbitmq:management
``

### Start the Container
 Run the following command to start a RabbitMQ container. The management plugin provides a web-based UI that can be accessed at http://localhost:15672/. The default username and password are guest/guest.

``
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
``

This command starts a RabbitMQ server and exposes port 5672 for AMQP protocol and port 15672 for the management UI.

### Install Erlang
You have to install [Erlang](https://www.erlang.org/downloads)

### Create a queue
Login to the RabbitMQ Management UI (http://localhost:15672/).<br>
Go to `Queues and Streams`.<br>
Click on `Add Queues`.<br>
Name the queue und leave everything by the default settings.

### Publish a json message
Login to the RabbitMQ Management UI (http://localhost:15672/).<br>
Go to `Queues and Streams`.<br>
Select the queue you just created.<br>
Click on `Publish Message` so fire a message in the queue.<br>

## Config the RabbitMQ Modifier
Open the `.env` file to update the RabbitMQ connection URL.

Now you are ready to modify this message via RabbitMQ Modifier.

## Usage
Pull this Repo and navigate inside the project folder with the terminal. Perform a `go mod download` to download the dependencies.<br>
With `go run main.go` you can execute the service.

