#Starting the Distributed Key-Value Store
#Master Node
The master node is responsible for managing metadata, coordinating nodes, and ensuring data consistency.

#Starting the Master Node
#Navigate to the Project Directory

sh
Copy code
cd path/to/your/project
Run the Master Node

Use the following command to start the master node:

sh
Copy code
go run server.go --port=3000 --shard=2 --replica=false --host=http://localhost
--port=3000: Specifies the port number for the master server.
--shard=2: Specifies the number of shards to be managed.
--replica=false: Indicates that this is not a replica.
--host=http://localhost: Specifies the domain on which the process is running.

Replica Nodes
Replica nodes are responsible for holding copies of the data to provide redundancy and fault tolerance.

#Starting a Replica Node
Navigate to the Project Directory
 Use the following command to start a replica node:

 go run server.go --port=3001 --shard=2 --replica=true --master=http://localhost:3000 --host=http://localhost
--port=3001: Specifies the port number for the replica server.
--shard=2: Specifies the number of shards to be managed.
--replica=true: Indicates that this is a replica.
--master=http://localhost:3000: Specifies the address of the master node.
--host=http://localhost: Specifies the domain on which the process is running.
Configuration Parameters
port: The port number on which the server listens.
shard: The number of shards managed by the server.
replica: A boolean flag indicating if the server is a replica (true) or not (false).
master: The address of the master node (applicable for replica nodes).
host: The domain on which the server process is running.
Example
To run a complete setup with one master and one replica:

Start the Master Node

sh
Copy code
go run server.go --port=3000 --shard=2 --replica=false --host=http://localhost
Start the Replica Node

sh
Copy code
go run server.go --port=3001 --shard=2 --replica=true --master=http://localhost:3000 --host=http://localhost
Monitoring and Logs
Logs will be output to stdout and can be redirected to a file for persistence. Monitoring can be configured through additional endpoints or integrated with external monitoring tools.
 