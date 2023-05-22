import { withAuth } from "@/app/api/helpers";
import { prisma } from "@/app/prisma/prisma";
import { NextRequest, NextResponse } from "next/server";

export const GET = withAuth(async (_request: NextRequest, _token, { params }: { params: { chatId: string } }) => {
  const messages = await prisma.message.findMany({
    where: {
      chat_id: params.chatId,
    },
    orderBy: {
      created_at: "asc"
    }
  });

  return NextResponse.json(messages);
});

export const POST = withAuth(async (request: NextRequest, _token, { params }: { params: { chatId: string } }) => {
  const body = await request.json();

  const chat = await prisma.chat.findUniqueOrThrow({
    where: {
      id: params.chatId,
    }
  });

  const messageCreated = await prisma.message.create({
    data: {
      content: body.message,
      chat_id: chat.id
    }
  });

  return NextResponse.json(messageCreated);
});
