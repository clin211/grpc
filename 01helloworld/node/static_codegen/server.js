import grpc from '@grpc/grpc-js';

import message from './rpc/helloworld_pb';
import service from './rpc/helloworld_grpc_pb';

function sayHello (call, callback) {
    const reply = new message.HelloReply();
    reply.setMessage('Hello ' + call.request.getName());
    callback(null, reply);
}

function sayHelloAgain (call, callback) {
    const reply = new message.HelloReply();
    reply.setMessage('Hello again, ' + call.request.getName());
    callback(null, reply);
}

const server = new grpc.Server();

server.addService(service.Greeter.service, { sayHello, sayHelloAgain });
server.bindAsync('0.0.0.0:50054', grpc.ServerCredentials.createInsecure(), () => {
    server.start();
});
