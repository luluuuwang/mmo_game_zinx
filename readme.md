go run main.go 开启服务器
client/client.exe,   port:127.0.0.1

2025/10/10: msgChan有缓冲；Writer批量取，一次写；背压策略为去除最旧的，塞入最新的消息
原来的问题：
1.每次sendmsg向msgchan发消息必须等待writer处理完(等待msgchan有位置)，如果同一时间大量消息要发送，则writer处理不完，sendmsg会阻塞
2.每条消息都要write一次，频繁的系统调用，Write:用户态-内核态发送-返回用户态
将Writer监控的msgChan改为有缓冲，以应对高流量时SendMsg一直阻塞等待情况，背压策略为缓冲区满时，拿出最旧的消息，塞入最新的，
并且由于Writer从msgChan取一次消息就要写一次Write，优化为批量取，一次写，减少系统调用次数。
客户端不用改的原因是客户端一直在用ReadFull(8)读头，然后再读体，
批量写只是把多条帧连在一起一起写，帧边界仍由datalen确定。
一次写的时候不用conn.Write()是因为Write会将batch里的n条数据进行n次系统调用，而writev()可以一次性写出多个buffer


2025/10/20：新绑定swagger，将项目暴露出http api接口，通过swagger+gin进行在线调试与契约输出
启动swagger：
go run .  // 启动全部服务
go run ./cmd/adminhttp    // 启动swagger文档
访问   http://localhost:18080/swagger/index.html#/