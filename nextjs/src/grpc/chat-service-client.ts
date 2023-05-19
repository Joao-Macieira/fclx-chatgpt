import { Metadata } from '@grpc/grpc-js';
import { chatClient } from './client';
import { ChatServiceClient as GrpcChatServiceClient } from './rpc/pb/ChatService';

export class ChatServiceClient {
  private authorization = "123456";
  constructor(private grpcClient: GrpcChatServiceClient) {}

  chatStream(data: { chat_id?: string, user_id: string, message: string }) {
    const metadata = new Metadata();
    metadata.set('Authorization', this.authorization);

    const stream = this.grpcClient.ChatStream({
      chatId: data.chat_id,
      userId: data.user_id,
      userMessage: data.message
    }, metadata);

    stream.on('data', data => {
      console.log(data);
    });

    stream.on('error', error => {
      console.log(error);
    });

    stream.on('end', () => {
      console.log('end');
    });

    return stream;
  }
}

export class ChatServiceClientFactory {
  static create() {
    return new ChatServiceClient(chatClient);
  }
}