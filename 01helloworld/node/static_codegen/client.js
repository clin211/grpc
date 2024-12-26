import grpc from '@grpc/grpc-js';

import message from './rpc/helloworld_pb';
import service from './rpc/helloworld_grpc_pb';

const client = new service.GreeterClient('0.0.0.0:50054', grpc.credentials.createInsecure());
const request = new message.HelloRequest();
request.setName('clina');
client.sayHello(request, (err, response) => {
    if (err) {
        console.error(err);
        return;
    }
    console.log(`Hello ${response.getMessage()}`);
});

client.sayHelloAgain(request, (err, response) => {
    if (err) {
        console.error(err);
        return;
    }
    console.log(`Hello again ${response.getMessage()}`);
});