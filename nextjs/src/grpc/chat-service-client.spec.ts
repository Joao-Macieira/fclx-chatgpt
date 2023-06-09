import { ChatServiceClientFactory } from "./chat-service-client";

describe("ChatSetviceClient", () => {
  test("grpc client", (done) => {
    const chatService = ChatServiceClientFactory.create();

    const stream = chatService.chatStream({
      user_id: '1',
      message: 'Hello world',
      chat_id: null
    });

    stream.on('end', () => {
      done();
    })
  });
});