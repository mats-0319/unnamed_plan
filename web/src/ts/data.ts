export const tips_ModifyUser = "昵称、密码和TOTP密钥字段为空表示不修改"

export const tips_RegisterUser = "请输入用户名和密码后点击注册"

interface FlipResult {
    duration: number;
    steps: number;
}

export function isFlipResult(obj: any):obj is FlipResult {
    return (
        typeof obj === "object" &&
        obj !== null &&
        typeof obj.duration === "number" &&
        typeof obj.steps === "number"
    )
}
