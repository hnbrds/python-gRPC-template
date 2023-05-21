# python-gRPC-template

## About
The gRPC-Gateway is a plugin of Google protocol buffers compiler. For more details, refer to [gRPC-gateway official GitHub](https://github.com/grpc-ecosystem/grpc-gateway). This repository contains a "**Python** gRPC server + **golang** reverse proxy server template". This helps you handle HTTP(REST) request, Swagger UI, gRPC request in the same port.


## Requirements
- Go, buf
- Python


## Getting Started
1. Make your own proto inside the `proto` folder (remove `sample.proto`)
2. Add your request handlers, and customizey your directory path, ports, etc.
   - Default PORT is set to `13270` for gateway_server, `13271` for python gRPC server.
3. Write a code to test gRPC API inside `grpc_client.py`.
4. run `./setup_env.sh` and check everything is all in place.
5. run `./run_server.sh`
6. Test
   - gRPC API with `python gprc_client.py`
   - RESTful API with your REST API testing tool.
   - Swagger UI in [http://127.0.0.1:{your_port}/swagger-ui/]()

  
## Updating Swagger UI
1. Download the latest swagger package from [Swagger UI official Github](https://github.com/swagger-api/swagger-ui).
2. Copy **dist** folder from the unzipped swagger package to the root of your project folder.
3. Rename directory name to `swagger`
4. Open `swagger-initializer.js` and change `url` value to `swagger.json`
5. After running `./setup_env.sh`, make sure you have created a soft link of swagger.json inside the swagger folder.

## Building Docker image
1. All your required Python pakcage should be listed in `requirements.txt`
2. run `./setup_env.sh`
3. Build your docker image (`docker build -t {your_tag} .`)

## License
- python-gRPC-template is licensed under the MIT License. See LICENSE.txt for more details.