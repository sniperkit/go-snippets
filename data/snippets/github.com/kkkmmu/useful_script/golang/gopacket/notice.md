1. 为了用RAW Socket接收OSPF报文，需要关闭防火墙.
2. 使用recvfrom接收报文的时候buffer 应该是make(buffer, BUFLEN) 而不能 是make(buffer, 0, BUFLEN)，否则的话，内部会认为buffer长度为0.
3. Buffer大小也对报文收发有影响.

