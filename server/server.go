package main

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	pb "rpcserver/gencode"
	"time"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func main() {
	lis, err := net.Listen("tcp", ":9010")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//简单模式
//返回两个整数相加的结果
func (s *server) Add(ctx context.Context, in *pb.TwoNum) (*pb.Response, error) {
	return &pb.Response{C: in.A + in.B}, nil
}

//服务端流
//返回多次数据，一次是加结果，一次是乘结果，当然还可以有更多返回
func (s *server) GetStream(in *pb.TwoNum, pipe pb.Calculator_GetStreamServer) error {
	_ = pipe.Send(&pb.Response{C: in.A + in.B})
	time.Sleep(time.Second * 2)
	_ = pipe.Send(&pb.Response{C: in.A * in.B})
	return nil
}

//客户端流
//客户端不停发送数据过来，服务端将所有的额数据累加，最后返回总和
func (s *server) PutStream(pipe pb.Calculator_PutStreamServer) error {
	var res int32
	for {
		request, err := pipe.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err.Error())
		}
		res += request.A
	}
	_ = pipe.SendAndClose(&pb.Response{C: res})
	return nil
}

//双向流
//客户端不停的发送数据，服务端将每次的数据相加返回，客户端和服务端都会接受和返回多次数据
func (s *server) DoubleStream(pipe pb.Calculator_DoubleStreamServer) error {
	for {
		request, err := pipe.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		_ = pipe.Send(&pb.Response{C: request.A + request.B})
	}
}
