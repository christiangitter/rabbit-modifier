# RabbitMQ Modifier
## Requirements for local testing
### Install Docker Desktop
Go to the [Docker Desktop Website](https://www.example.com) to install the latest version. 

### Pull the RabbitMQ Image
``docker pull rabbitmq:management
``

### Start the Container
 Run the following command to start a RabbitMQ container. The management plugin provides a web-based UI that can be accessed at http://localhost:15672/. The default username and password are guest/guest.

``
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management
``

This command starts a RabbitMQ server and exposes port 5672 for AMQP protocol and port 15672 for the management UI.

