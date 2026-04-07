import { router } from "@/router.ts"

// deepCopy 简单的deep copy，没有考虑嵌套对象的情况
export function deepCopy<T extends object>(obj: T): T {
    const res: T = {} as T

    for (const key in obj) {
        res[key] = obj[key]
    }

    return res
}

export function routerLink(name: string): void {
    router.push({ name: name })
}

export function displayTimestamp(timestamp: number): string {
    if (timestamp == 0) {
        return "无"
    }

    return new Date(timestamp).toLocaleString()
}

export function randomVisitorName(): string {
    const array = new Uint32Array(1)
    crypto.getRandomValues(array)

    return "游客" + array[0].toString().padStart(10, "0").slice(0, 10)
}
