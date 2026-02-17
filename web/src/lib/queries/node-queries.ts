import { useQuery } from '@tanstack/react-query'
import { api } from '@/src/lib/api/client'

const nodeKeys = {
    all: ['nodes'] as const,
    lists: () => [...nodeKeys.all, 'list'] as const,
    listBySub: (subId: number | null) => [...nodeKeys.lists(), { subId }] as const,
}

export function useNodes(subId: number | null, includeFailed = false) {
    return useQuery({
        queryKey: [...nodeKeys.listBySub(subId), { includeFailed }],
        queryFn: () => api.getNodes(subId ?? undefined, includeFailed),
        enabled: subId !== null,
        refetchInterval: 60 * 1000,
        notifyOnChangeProps: ['data', 'error', 'isLoading'],
    })
}
