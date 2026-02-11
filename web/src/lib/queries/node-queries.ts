import { useQuery } from '@tanstack/react-query'
import { api } from '@/src/lib/api/client'

const nodeKeys = {
    all: ['nodes'] as const,
    lists: () => [...nodeKeys.all, 'list'] as const,
    listBySub: (subId: number | null) => [...nodeKeys.lists(), { subId }] as const,
}

export function useNodes(subId: number | null) {
    return useQuery({
        queryKey: nodeKeys.listBySub(subId),
        queryFn: () => api.getNodes(subId ?? undefined),
        enabled: subId !== null,
        refetchInterval: 60 * 1000,
        notifyOnChangeProps: ['data', 'error', 'isLoading'],
    })
}

