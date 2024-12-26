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

// 创建一个Greeter客户端实例，连接到localhost:50053端口
// 使用grpc.credentials.createInsecure()创建一个不安全的凭证
const client = new helloProto.Greeter('localhost:50053', grpc.credentials.createInsecure());

// 调用sayHello方法，发送请求到服务器
client.sayHello({ name: 'clina' }, function (err, response) {
    // 输出服务器返回的消息
    console.log('Greeting:', response.message);
});

// 调用sayHelloAgain方法，发送请求到服务器
client.sayHelloAgain({ name: 'clina' }, (err, response) => {
    // 输出服务器返回的消息
    console.log('Greeting:', response.message);
})