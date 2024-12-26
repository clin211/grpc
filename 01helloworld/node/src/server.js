// 导入grpc和protoLoader模块
import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';

// 定义proto文件路径
const protoPath = `${process.cwd()}/../proto/helloworld.proto`;

// 加载proto文件，定义proto文件的选项
const packageDefinition = protoLoader.loadSync(protoPath, {
    // keepCase: 保持proto文件中的字段名不变
    keepCase: true,
    // longs: 将long类型转换为string类型
    longs: String,
    // enums: 将枚举类型转换为string类型
    enums: String,
    // defaults: 使用proto文件中的默认值
    defaults: true,
    // oneofs: 支持oneof语法
    oneofs: true
});

// 从packageDefinition中获取helloworld模块
const helloProto = grpc.loadPackageDefinition(packageDefinition).helloworld;

// 定义sayHello函数，处理客户端的sayHello请求
function sayHello (call, callback) {
    // 回调函数，返回响应消息
    callback(null, { message: 'Hello ' + call.request.name });
}

// 定义sayHelloAgain函数，处理客户端的sayHelloAgain请求
function sayHelloAgain (call, callback) {
    // 回调函数，返回响应消息
    callback(null, { message: 'Hello again, ' + call.request.name });
}

// 创建一个grpc服务器实例
const server = new grpc.Server();

// 添加Greeter服务到服务器
server.addService(helloProto.Greeter.service, { sayHello, sayHelloAgain });

// 绑定服务器到0.0.0.0:50053端口
server.bindAsync('0.0.0.0:50053', grpc.ServerCredentials.createInsecure(), () => {
    // 启动服务器
    server.start();
});