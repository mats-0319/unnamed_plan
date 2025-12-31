import { router } from "@/router.ts"

// deepCopy 简单的deep copy，没有考虑嵌套对象的情况
export function deepCopy<T extends object>(obj: T): T {
	let res: T = {} as T

	for (let key in obj) {
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
