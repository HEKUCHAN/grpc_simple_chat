# gRPCを使用したCLIチャットアプリの作成

## 背景
近年、AI技術の発展に伴って、リアルタイム性が求められるアプリケーションの需要が増えています。特に、AIとの親和性が高いPythonをAI処理部分で使用し、バックエンドや高速な通信処理をGolangやJavaで実装するなどの技術選定が増えています。

また、システムが大規模化するにつれて、異なる言語のシステム間での通信が求められる場面が増えています。従来のREST APIでは、通信のオーバーヘッドが大きく、リアルタイム通信には適していません。そこで、高速かつ効率的なRPC通信を実現する gRPC を活用し、多言語対応のリアルタイムCLIチャットアプリを構築することにしました。

## 目的
- GolangとPythonで動作するgRPCサーバーおよびクライアントを構築し、相互に通信できるようにする
- CLI形式でリアルタイムチャットを行えるようにし、必要に応じて複数クライアントを同時に接続可能にする
- gRPCの多言語対応を活かし、異なる言語のクライアントでもシームレスに通信可能にする

## 手法
### gRPCについての理解を深める
[gRPC](https://grpc.io/)はGoogleにより2015年に開発されたオープンソースなRPCフレームワークです。
HTTP/2をトランスポートとして利用して、[Protocol Buffers](https://protobuf.dev/)などのプラットフォームや言語に依存しない拡張性のある構造化データのシリアライズ手法を利用することで、効率的なデータのエンコード・デコードが実現され、システム間のデータ交換がスムーズに行えます。

また、バックエンドなどのシステム間通信においては、gRPCはProtocol Buffersと連携して、明確に定義されたサービスとメソッドに基づくクライアントとサーバー間の通信インターフェースを自動生成し、効率的なリモートプロシージャ呼び出しを可能にします。
これにより、バックエンド同士の接続実装が容易になり、システムのスケーリングが行いやすくなる特徴があります。

### Protocol Buffersを使用してgRPCをシリアライズ化
![image](https://github.com/user-attachments/assets/160cc4e7-dd7c-4310-bbf6-a3d074cbd618)

ここでは詳細な解説をはぶきますが、ChatServiceという名前でgRPCサービスを定義し、サービス内にChatメゾットを定義し、ChatMessageというユーザー名とメッセージ本文を持っているデータ型をストリーム形式で送受信を行うことを定義しています。
[該当コードページ](https://github.com/HEKUCHAN/grpc_simple_chat/blob/main/protos/chat.proto)

この作成したファイルを使用することにより、これから作成する `Go言語製のサーバー`、`Go言語製のクライアント`、`Python言語製のクライアント`のそれぞれに対して、Protocol Buffersから自動生成されたコードを利用し、統一されたインターフェースで通信を実現します。[
Protocol Buffersのドキュメント](https://protobuf.dev/reference/go/go-generated/)を見ていただけると`protoc`コマンドを使用してそれぞれの言語向けのコードを自動生成することができることがわかります。

それを使用し、それぞれのプラットフォームの言語に合わせて自動生成させました。
- [サーバー用に自動生成させたコード](https://github.com/HEKUCHAN/grpc_simple_chat/tree/main/apps/server/proto)
- [Pythonクライアント用に自動生成させたコード](https://github.com/HEKUCHAN/grpc_simple_chat/tree/main/apps/python_client/proto)
- [Goクライアント用に自動生成させたコード](https://github.com/HEKUCHAN/grpc_simple_chat/tree/main/apps/golang_client/proto)

コードを生成した後に、必要なのはそれぞれのプラットフォームの内部処理を実装することです。

### gRPCサーバーの構築
主に双方向ストリーミングを活用し、サーバーが受信したメッセージを自分以外のクライアントへ配信する仕組みにしています。
サーバーが接続情報を一括管理することで、各クライアントは自分自身の送信メッセージをサーバーに投げるだけで、ほかのすべてのクライアントへメッセージが届くようになります。

このような処理にすることで、今後チャットルームやユーザー認証の概念を追加するのが簡単になると思いました。

[該当コードページ](https://github.com/HEKUCHAN/grpc_simple_chat/tree/main/apps/server)

### Golang製のCLIクライアントの実装
今回の実装では認証機能や複雑のロジックを実装せずにシンプルな構成になりました。
サーバーとの接続が確立したのにユーザー名を入力してもらい、サーバーへのメッセージが送信できるようにしたり、サーバーから返ってきたメッセージをコンソールに表示する仕組みを実装しました。
今後はここに認証用のトークン管理やユーザー情報の登録処理などを追加することで、より実運用に近いチャットクライアントに発展させることが可能だと思います。

[該当コードページ](https://github.com/HEKUCHAN/grpc_simple_chat/tree/main/apps/golang_client)

### Python製のCLIクライアントの実装
PythonクライアントはGoクライアントと大きく異なる部分はありません。
同じ.protoファイルから生成されたコードを利用し、双方向ストリーミングを通じたメッセージ送受信を実現することができています。

GoでもPythonでも、ベースにある通信仕様は同じProtocol Buffers定義に依拠しているため、言語が変わってもサーバーと問題なく通信できることが、gRPCの大きな強みというのを感じることができました。

[該当コードページ](https://github.com/HEKUCHAN/grpc_simple_chat/tree/main/apps/python_client)

## 完成したもの
左側 : Go言語で作成されたサーバー, 右下 : Python言語で作成されたクライアント, 右上 : Go言語で作成されたクライアント

https://github.com/user-attachments/assets/7d27bc5c-f4e1-4ce9-b55e-d94d2410309c

## 感想
今回のCLI上で使えるチャットアプリ作成を通じて、gRPCの多言語対応の強みを実感することができました。
特にProtocol Buffersを使用したコードの自動生成を使用することにより、Goでサーバーを実装しても、Python、Goでクライアントを実装しても、整合性が取れた状態で通信することができ、チーム開発などでもとても活躍すると思いました。

さらに、GoとPythonといった異なる言語を組み合わせることで、性能と開発効率をバランスよく両立できる可能性を感じました。
具体的には、Goの並行処理性能やメモリ効率の良さを活かしてサーバーの負荷を軽減しつつ、Pythonの豊富なライブラリやAIフレームワークを活用することで、チャット内容のリアルタイム解析や機械学習モデルとの組み合わせといった機能追加ができると思います。


今後もgRPCを活用し新しいアプリを開発していこうと思います！
