import grpc
from chat.v1 import chat_pb2
from chat.v1 import chat_pb2_grpc

def main():
    channel = grpc.insecure_channel('localhost:8080')
    # (carter): we'd probably want to wrap/stub out `chat_pb2_grpc`. Offends the eyes, if you ask me.
    client = chat_pb2_grpc.ChatServiceStub(channel)

    user_resp = client.CreateUser(chat_pb2.CreateUserRequest(username="carter"))
    print(f"User: {user_resp.user}")

    server_resp = client.CreateServer(chat_pb2.CreateServerRequest(name="grpc nice"))
    print(f"Server: {server_resp.server}")

    channel_resp = client.CreateChannel(
        chat_pb2.CreateChannelRequest(
            server_id=server_resp.server.id,
            name='channel1'
        )
    )
    print(f"Channel: {channel_resp.channel}")

    msg_resp = client.SendChannelMessage(
        chat_pb2.SendChannelMessageRequest(
            server_id=server_resp.server.id,
            channel_id=channel_resp.channel.id,
            author_user_id=user_resp.user.id,
            content="grpc good make life easier veri nice"
        )
    )
    print(f"Message: {msg_resp.message}")


if __name__ == '__main__':
    main()