'use client';

import { PropsWithChildren } from 'react';
import { Session } from 'next-auth';
import { SessionProvider as NextAuthSessionProvider } from 'next-auth/react';

type SessionProvidersProps = PropsWithChildren<{
  session: Session | null,
}>

export function SessionProvider(props: SessionProvidersProps) {
  return (
    <NextAuthSessionProvider session={props.session}>
      {props.children}
    </NextAuthSessionProvider>
  );
}