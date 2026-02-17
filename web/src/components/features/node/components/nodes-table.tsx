import { Badge } from "@/src/components/ui/badge"
import { Card, CardContent } from "@/src/components/ui/card"
import { Table, TableBody, TableHead, TableHeader, TableRow, TableCell } from "@/src/components/ui/table"
import { InlineLoading } from "@/src/components/ui/loading"
import { formatSpeed } from "@/src/components/features/sub/utils"
import { formatNodeStatus, getRiskClass, getRiskLabel } from "../utils"
import { NODE_STATUS } from "../constants"
import type { NodeResponse } from "@/src/types"

interface NodesTableProps {
    nodes: NodeResponse[]
    isLoading: boolean
    error: Error | null
}

export function NodesTable({ nodes, isLoading, error }: NodesTableProps) {
    if (isLoading) {
        return (
            <Card>
                <CardContent>
                    <InlineLoading message="加载节点列表..." />
                </CardContent>
            </Card>
        )
    }

    if (error) {
        return (
            <Card>
                <CardContent>
                    <div className="text-center py-8 text-destructive">
                        加载失败: {error.message}
                    </div>
                </CardContent>
            </Card>
        )
    }

    if (nodes.length === 0) {
        return (
            <Card>
                <CardContent>
                    <div className="text-center py-8 text-muted-foreground">
                        暂无节点数据
                    </div>
                </CardContent>
            </Card>
        )
    }

    const orderedNodes = nodes.slice().sort((a, b) => {
        const aAlive = (a.alive_status & NODE_STATUS.ALIVE) !== 0
        const bAlive = (b.alive_status & NODE_STATUS.ALIVE) !== 0
        if (aAlive === bAlive) return 0
        return aAlive ? -1 : 1
    })

    return (
        <Card className="min-w-0">
            <CardContent>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>节点名称</TableHead>
                            <TableHead>类型</TableHead>
                            <TableHead>状态</TableHead>
                            <TableHead>延迟</TableHead>
                            <TableHead>上下行</TableHead>
                            <TableHead>风险</TableHead>
                            <TableHead>国家</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {orderedNodes.map((node) => {
                            const statusLabels = formatNodeStatus(node.alive_status)
                            return (
                            <TableRow key={`${node.sub_id}-${node.unique_key}`}>
                                <TableCell className="font-medium">
                                    {node.name || '未命名节点'}
                                    {node.reason ? (
                                        <div className="text-xs text-muted-foreground">原因: {node.reason}</div>
                                    ) : null}
                                </TableCell>
                                    <TableCell>{node.type || 'N/A'}</TableCell>
                                    <TableCell className="space-x-1">
                                        {statusLabels.length === 0 ? (
                                            <Badge variant="outline">非存活</Badge>
                                        ) : (
                                            statusLabels.map((label) => (
                                                <Badge key={label} variant="outline">{label}</Badge>
                                            ))
                                        )}
                                    </TableCell>
                                    <TableCell>{node.delay ? `${node.delay}ms` : 'N/A'}</TableCell>
                                    <TableCell className="text-xs text-muted-foreground">
                                        ↑{formatSpeed(node.speed_up || 0)} ↓{formatSpeed(node.speed_down || 0)}
                                    </TableCell>
                                    <TableCell>
                                        <span className={`text-xs font-medium ${getRiskClass(node.risk || 0)}`}>
                                            {getRiskLabel(node.risk || 0)}
                                        </span>
                                    </TableCell>
                                    <TableCell>{node.country || '未知'}</TableCell>
                                </TableRow>
                            )
                        })}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    )
}
