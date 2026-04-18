// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

import { Pagination } from "./common.go"

export class RegisterReq {
    user_name: string = "" // unique, also default nickname
    password: string = "" // hex(sha256('text'))
}

export class RegisterRes {}

export class LoginReq {
    user_name: string = ""
    password: string = ""
}

export class LoginRes {
    user_name: string = ""
    nickname: string = ""
    is_admin: boolean = false
    enable_mfa: boolean = false
    mfa_token: string = "" // empty when disable 2fa
}

export class LoginMFAReq {
    mfa_token: string = ""
    totp_code: string = ""
}

export class LoginMFARes {
    user_name: string = ""
    nickname: string = ""
    is_admin: boolean = false
    enable_mfa: boolean = false
}

export class User {
    user_name: string = "" // login name, can't modify
    nickname: string = "" // display name
    created_at: number = 0
    is_admin: boolean = false
    enable_mfa: boolean = false
    last_login: number = 0 // timestamp, unit: milli
}

export class ListUserReq {
    page: Pagination = new Pagination()
}

export class ListUserRes {
    count: number = 0 // 符合查询条件的用户总数
    users: Array<User> = new Array<User>()
}

// ModifyUserReq string类型的属性为空，视为不修改对应字段
export class ModifyUserReq {
    nickname: string = ""
    password: string = ""
    enable_mfa: boolean = false
    totp_key: string = ""
}

export class ModifyUserRes {}
