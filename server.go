package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	grpc "google.golang.org/grpc"
	pb "./protobuf"
	"io"
	"net"
)

type authServer struct{
	pb.UnimplementedAUTHServer
}

func (*authServer) AuthLogin(stream pb.AUTH_AuthLoginServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Printf("\n收到的用户名：%s, 密码：%s.", req.Username, req.Password)
		c := make(chan string)
		go str2base64(c, req.Username, req.Password)
		for n := range c {
			resp := &pb.Response{
				Result: string(n),
			}
			stream.Send(resp)
		}
	}
	return nil

}

func newAuthServer() pb.AUTHServer {
	return &authServer{}
}

func str2base64(c chan string, username string, password string) {
	result := base64.StdEncoding.EncodeToString([]byte(username + password))
	fmt.Printf("\nBase64编码后的结果为: %s.", result)
	fmt.Printf("\n--------------------")
	c <- result
	close(c)
}

func main() {
	port := flag.Int("p", 12345, "服务运行端口")
	flag.Parse()

	fmt.Printf("认证服务启动, 运行端口为: %d", *port)
	fmt.Printf("\n--------------------")
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAUTHServer(grpcServer, newAuthServer())
	grpcServer.Serve(conn)
}
