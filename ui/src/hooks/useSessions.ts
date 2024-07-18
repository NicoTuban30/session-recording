import { useQuery } from '@tanstack/react-query';

export interface Session {
  id: string;
  created: string;
  updated: string;
  origin: string;
  address: string;
  userAgent: string;
  userEmail: string;
  qaId: string;
  qaSessionId: string;
  agoraStreamUrl: string;
}

export default function useSessions() {
  const query = useQuery<Session[]>({
    queryKey: ['sessions'],
    queryFn: () => fetch(`/sessions`).then((res) => res.json()),
  });

  return query;
}
