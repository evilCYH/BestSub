export interface NodeResponse {
    sub_id: number
    unique_key: number
    name: string
    type: string
    delay: number
    speed_up: number
    speed_down: number
    risk: number
    alive_status: number
    country: string
}

export interface NodeUpdateLog {
    sub_id: number
    created_at: string
    duration_ms: number
    raw_count: number
    candidate: number
    duplicate: number
    invalid: number
    test_failed: number
    accepted: number
    merged: number
    dropped: number
    details?: string[]
}

export interface NodeUpdateLogResponse {
    latest?: NodeUpdateLog
    history: NodeUpdateLog[]
}

// 节点详细测试日志
export interface NodeTestLog {
    id: number
    sub_id: number
    node_name: string
    level: 'info' | 'warn' | 'error'
    message: string
    created_at: string
}

export interface NodeTestLogResponse {
    total: number
    list: NodeTestLog[]
}
