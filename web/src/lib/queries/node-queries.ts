import { useQuery } from '@tanstack/react-query'
import { api } from '@/src/lib/api/client'

const nodeKeys = {
    all: ['nodes'] as const,
    lists: () => [...nodeKeys.all, 'list'] as const,
    listBySub: (subId: number | null) => [...nodeKeys.lists(), { subId }] as const,
}

export function useNodes(
    subId: number | null,
    options?: {
        includeFailed?: boolean
        scope?: 'registry' | 'pool'
        status?: 'alive' | 'dead' | 'init_failed' | 'all'
    }
) {
    return useQuery({
        queryKey: [...nodeKeys.listBySub(subId), options ?? {}],
        queryFn: () => api.getNodes({ subId: subId ?? undefined, ...options }),
        enabled: subId !== null,
        refetchInterval: 60 * 1000,
        notifyOnChangeProps: ['data', 'error', 'isLoading'],
    })
}
