import grpc
import threading
import queue
import sys
from python_client.proto import chat_pb2, chat_pb2_grpc

def request_generator(req_queue):
    """Queueからメッセージを取り出して、ストリームに渡すジェネレーター"""
    while True:
        message = req_queue.get()
        if message is None:
            break
        yield message

def main():
    channel = grpc.insecure_channel('localhost:50051')
    client = chat_pb2_grpc.ChatServiceStub(channel)

    req_queue = queue.Queue()

    responses = client.Chat(request_generator(req_queue))

    username = input("ユーザー名を入力してください: ")

    def receive_messages():
        try:
            for response in responses:
                if response.user == username:
                    continue

                sys.stdout.write("\r\033[K")
                print(f"[{response.user}]: {response.message}")
                sys.stdout.write(">> ")
                sys.stdout.flush()
        except Exception as e:
            print("受信中にエラー:", e)

    threading.Thread(target=receive_messages, daemon=True).start()

    while True:
        text = input(">> ")
        if text.strip() == "exit":
            break
        req_queue.put(chat_pb2.ChatMessage(user=username, message=text))

    req_queue.put(None)

if __name__ == "__main__":
    main()
