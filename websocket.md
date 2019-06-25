#websocket消息设计
1. 通过ws:<hostname>:<port>/path?token=<token>&platform=<platform>&id=<id>连接websocket
2. 建立Conn（id, socket, recive ch, send ch, close ch, userid, platform）
3. 将Conn 加入syncMap中以userid+platfrom为键, 以及另外的socket指针为键的syncMap
4. 用户发送一条message(id, serviceId, to, form, data, timestamp, messageType)
5. socket 收到一条消息, socket通过socket指针地址找到对应的Conn， 向Conn的recive ch发送一条消息
6. recive ch处理业务逻辑，通过serviceId, messageType等处理对应的消息 
7. 服务端接受到一条message 通过to, from, serviceId, 和messageType 判断是否有权限发送该消息
8. dispatcher 有 sendToOne sendToManny
9. 有权限之后，通过to和messageType通过syncMap找出所有的socket
10. 调用dispatcher 发送消息
11. 向socket的send ch发送一条消息， send ch发送消息，
