import path from 'path';
import * as protoLoader from '@grpc/proto-loader';
import * as grpc from '@grpc/grpc-js';

import { ProtoGrpcType } from './rpc/chat';

const packageDefinitions = protoLoader.loadSync(
  path.resolve(process.cwd(), 'proto', 'chat.proto')
);

const proto = grpc.loadPackageDefinition(packageDefinitions) as unknown as ProtoGrpcType;

export const chatClient = new proto.pb.ChatService("host.docker.internal:50052", grpc.credentials.createInsecure());
