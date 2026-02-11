import { NODE_STATUS } from './constants'

export function formatNodeStatus(flags: number): string[] {
    const labels: string[] = []
    if (flags & NODE_STATUS.ALIVE) labels.push('存活')
    if (flags & NODE_STATUS.COUNTRY) labels.push('国家')
    if (flags & NODE_STATUS.TIKTOK) labels.push('TikTok')
    if (flags & NODE_STATUS.TIKTOK_IDC) labels.push('TikTok IDC')
    return labels
}

export function getRiskLabel(risk: number): string {
    if (risk >= 7) return '高风险'
    if (risk >= 4) return '中风险'
    return '低风险'
}

export function getRiskClass(risk: number): string {
    if (risk >= 7) return 'text-red-600'
    if (risk >= 4) return 'text-yellow-600'
    return 'text-green-600'
}

