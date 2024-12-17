// 実行可能なGoアプリケーションであることを宣言する
package main

//docker buildx build --platform linux/amd64 -t us-central1-docker.pkg.dev/grpc-chat-app/grpc-chat-server/grpc-chat-server:v1.0 -f C:\Users\hatte\source\repos\grpc-chat-app-sample\Dockerfile C:\Users\hatte\source\repos\grpc-chat-app-sample
// どの言語でも大抵お馴染みのパッケージインポート
import (
	// クライアントからサーバーへのリクエストのスコープやキャンセル、タイムアウトの管理に使用する
	"context"
	// プログラムの実行中にログ（メッセージ）を記録するためのパッケージ
	"log"
	// ネットワーク関連の機能を提供
	// ここでは、サーバーがリクエストを受け付けるためのネットワークリスナーを作成するためにimport
	"net"

	// gRPCサーバーを作成するためのパッケージ
	"google.golang.org/grpc"
	// 生成されたgRPCコードをインポート
	// pbはこのパッケージを参照するための別名
	pb "grpc-chat-app-sample/gen/api/helloworld"
)

const (
	port = ":8080" // ここを変更
)

// サーバー構造体の定義
// gRPCで定義されたサービスを実装する
type server struct {
	// gRPCのサービスを実装する際に必要なデフォルトの設定を提供（後述）
	pb.UnimplementedGreeterServer
}

// 関数の定義
// server型に属すメソッドで、inという名前の引数（*pb.HelloRequest型のポインタ）を受け取り
// HelloReply型のポインタと、エラー情報を戻り値としている
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// サーバーが指定したポート番号（ここでは50051）でTCP接続を待ち受けるリスナーを作成
	// リスナーは、クライアントからの接続を待機する
	// ちなみに := は変数の宣言と初期化を同時に行う「短縮変数宣言」
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 新しいgRPCサーバーを作成
	// このサーバーがクライアントからのリクエストを処理する
	s := grpc.NewServer()

	// GreeterサービスをgRPCサーバーに登録
	// ここで、&server{}は先ほど定義したserver構造体のインスタンスを渡している
	pb.RegisterGreeterServer(s, &server{})

	// Addr()メソッドは、net.Listenerインターフェースに定義されているメソッド
	// リスナーが現在待ち受けているアドレス（IPアドレスとポート番号の組み合わせ）を返す
	// これにより、サーバーがどのアドレスで接続を待機しているかを確認できる
	log.Printf("server listening at %v", lis.Addr())

	// s.Serve(lis)でサーバーを起動し、クライアントからの接続をリスナーを通じて受け付ける
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
