package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	pb "rpcclient/gencode"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()
	client := pb.NewCalculatorClient(conn)

	Sample(client)
	GetStream(client)
	PutStream(client)
	DoubleStream(client)
}

//简单模式
//客户端发送一次请求，服务端返回一次数据
func Sample(client pb.CalculatorClient) {
	resp1, _ := client.Add(context.Background(), &pb.TwoNum{A: 10, B: 20})
	fmt.Println("普通模式： ", resp1.C, "\r\n")
}

//服务端流
//客户端发送一次请求，服务端返回多次数据
func GetStream(client pb.CalculatorClient) {
	serverPipe, _ := client.GetStream(context.Background(), &pb.TwoNum{A: 10, B: 20})
	for {
		resp2, err := serverPipe.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("服务端流： ", resp2.C)
	}
	fmt.Println("\r\n")
}

//客户端流
//客户端发送多次数据，服务端返回一次数据
func PutStream(client pb.CalculatorClient) {
	clientPipe, _ := client.PutStream(context.Background())
	for i := 1; i <= 100; i++ {
		_ = clientPipe.Send(&pb.OneNum{A: int32(i)})
	}
	resp3, err := clientPipe.CloseAndRecv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("客户端流： ", resp3.C, "\r\n")
}

//双向流
//客户端发送多次数据，服务端也返回多次数据，两者既要接受对方的数据，又要发送数据
func DoubleStream(client pb.CalculatorClient) {
	doublePipe, _ := client.DoubleStream(context.Background())
	ch := make(chan int32, 10)
	go func() {
		for {
			resp4, err := doublePipe.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err.Error())
			}
			ch <- resp4.C
		}
	}()
	go func() {
		for j := 1; j <= 10; j++ {
			time.Sleep(time.Second)
			_ = doublePipe.Send(&pb.TwoNum{A: int32(j), B: int32(j + 1)})
		}
	}()
	for k := 0; k < 10; k++ {
		fmt.Println("双向流： ", <-ch)
	}
}
