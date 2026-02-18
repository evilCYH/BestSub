import { useMemo, useState } from "react"
import { Input } from "@/src/components/ui/input"
import { SubscriptionColumn } from "./subscription-column"
import { NodesPanel } from "./nodes-panel"
import { useSubs } from "@/src/lib/queries/sub-queries"
import { useNodes } from "@/src/lib/queries/node-queries"
import type { SubResponse } from "@/src/types"

export function NodesPage() {
    const { data: subs = [], isLoading: subsLoading, error: subsError } = useSubs()
    const [selectedSub, setSelectedSub] = useState<SubResponse | null>(null)
    const [query, setQuery] = useState("")

    const filteredSubs = useMemo(() => {
        if (!query.trim()) return subs
        const keyword = query.trim().toLowerCase()
        return subs.filter((sub) => sub.name.toLowerCase().includes(keyword))
    }, [query, subs])

    const { data: nodes = [], isLoading: nodesLoading, error: nodesError } = useNodes(selectedSub?.id ?? null, {
        scope: 'registry',
        status: 'all',
        includeFailed: false,
    })

    return (
        <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
            <div className="flex items-center justify-between px-4 lg:px-6">
                <div>
                    <h1 className="text-2xl font-bold">节点管理</h1>
                    <p className="text-sm text-muted-foreground">选择订阅查看节点详情</p>
                </div>
            </div>

            <div className="px-4 lg:px-6">
                <div className="grid gap-4 lg:grid-cols-[0.7fr_1.6fr] min-w-0">
                    <div className="space-y-4">
                        <Input
                            placeholder="搜索订阅名称"
                            value={query}
                            onChange={(event) => setQuery(event.target.value)}
                        />
                        <SubscriptionColumn
                            subs={filteredSubs}
                            isLoading={subsLoading}
                            error={subsError as Error | null}
                            selectedId={selectedSub?.id ?? null}
                            onSelect={setSelectedSub}
                        />
                    </div>

                    <div className="min-w-0">
                        <NodesPanel
                            selectedSub={selectedSub}
                            nodes={nodes}
                            isLoading={nodesLoading}
                            error={nodesError as Error | null}
                        />
                    </div>
                </div>
            </div>
        </div>
    )
}
