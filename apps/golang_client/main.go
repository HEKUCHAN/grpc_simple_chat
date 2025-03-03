package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	chatpb "golang_client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient("localhost:50051", dialOpts...)
	if err != nil {
		log.Fatalf("サーバーへの接続に失敗しました: %v", err)
	}
	defer conn.Close()

	client := chatpb.NewChatServiceClient(conn)

	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("ストリームの作成に失敗しました: %v", err)
	}

	fmt.Print("ユーザー名を入力してください: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Fatalf("メッセージ受信エラー: %v", err)
			}

			if in.User == username {
				continue
			}

			fmt.Print("\r\033[K")
			fmt.Printf("[%s]: %s\n", in.User, in.Message)
			fmt.Print(">> ")
		}
	}()

	fmt.Print(">> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		err := stream.Send(&chatpb.ChatMessage{
			User:    username,
			Message: text,
		})
		if err != nil {
			log.Fatalf("メッセージ送信エラー: %v", err)
		}
		fmt.Print(">> ")
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("ストリームのクローズに失敗しました: %v", err)
	}
}
