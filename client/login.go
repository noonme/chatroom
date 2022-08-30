package main

//
//import (
//	"encoding/binary"
//	"encoding/json"
//	"fmt"
//	"go_code/chatroom/client/utils"
//	"go_code/chatroom/common/message"
//	"net"
//)
//
////login函数完成登录
//func login(userId int, userPwd string) (err error) {
//	//下一个就要开始定协议...
//	//fmt.Printf("userId = %d userPwd=%s\n",userId,userPwd)
//	//return nil
//
//	//1.连接到服务器
//	conn, err := net.Dial("tcp", "localhost:8889")
//	if err != nil {
//		fmt.Println("net.Dial err=", err)
//		return
//	}
//	//延时关闭
//	defer conn.Close()
//
//	//2.准备通过conn发送消息给服务端
//	var mes message.Message
//	mes.Type = message.LoginMesType
//
//	//3.创建一个LoginMes结构体
//	var loginMes message.LoginMes
//	loginMes.UserId = userId
//	loginMes.UserPwd = userPwd
//
//	//4. 讲loginMes序列化
//	data, err := json.Marshal(loginMes)
//	if err != nil {
//		fmt.Println("json.Marshal err=", err)
//		return
//	}
//
//	//4-5步的必要性，因为message包的Data是string类型，string()无法直接转换loginMes;
//	//这里使用json.Marsh过度一下
//
//	//5. 把data赋给mes.Date字段
//	mes.Data = string(data)
//
//	//6. 将mes （5这里导致mes是string需要再次转换）进行序列化
//	data, err = json.Marshal(mes)
//	if err != nil {
//		fmt.Println("json.Marsh err=", err)
//		return
//	}
//	//7.到这个时候data就是我们需要发送的消息
//
//	//7.1 先把data的长度发送给服务端
//	//将获取到的data长度-->转化为一个表示长度的byte切片？？
//	var pkgLen uint32
//	pkgLen = uint32(len(data))
//	fmt.Println("pkgLen=", pkgLen)
//	var buf [4]byte
//	//fmt.Printf("buf=%v,buf[:4]=%v\n", buf, buf[:4])
//	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //ByteOrder规定了如何将字节序列和 16、32或64比特的无符号整数互相转化。即格式转换
//	fmt.Printf("buf=%v,buf[:4]=%v\n", buf, buf[:4])
//
//	//发送长度
//	n, err := conn.Write(buf[:4])
//	if n != 4 || err != nil {
//		fmt.Println("conn.Write(bytes) failed", err)
//		return
//	}
//
//	//发送消息本身
//	_, err = conn.Write(data)
//	if err != nil {
//		fmt.Println("conn.Write(data) fail", err)
//		return
//	}
//
//	//休眠20s
//	//time.Sleep(20 * time.Second)
//	//fmt.Println("休眠20s...")
//	//fmt.Printf("客户端发送的消息长度=%d 内容=%s", len(data), string(data))
//
//	//这里还需要处理服务器返回的消息，==
//	mes, err = utils.readPkg(conn)
//	if err != nil {
//		fmt.Println("readPkg(conn) err=", err)
//		return
//	}
//
//	//将mes的Data部分反序列化成LoginResMes
//	var loginResMes message.LoginResMes
//	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
//	if loginResMes.Code == 200 {
//		fmt.Println("登录成功")
//	} else if loginResMes.Code == 500 {
//		fmt.Println(loginResMes.Error)
//	}
//
//	return
//}
