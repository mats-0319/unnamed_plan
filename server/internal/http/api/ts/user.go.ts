// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v1.0.0

import { Pagination } from "./common.go"

export class LoginReq {
    user_name: string = "";
    password: string = "";
    totp_code: string = "";
}

export class LoginRes {
    nickname: string = "";
    permission: number = 0;
    is_success: boolean = false;
    err: string = "";
}

export class User {
    id: number = 0;
    created_at: number = 0;
    updated_at: number = 0;
    name: string = "";
    nickname: string = "";
    totp_key: string = "";
    is_locked: boolean = false;
    permission: number = 0;
    creator: string = "";
    last_login: number = 0;
}

export class ListUserReq {
    operator: number = 0;
    page: Pagination = new Pagination();
    see_locked: boolean = false;
}

export class ListUserRes {
    user_amount: number = 0;
    list_amount: number = 0;
    users: Array<User> = new Array<User>();
    is_success: boolean = false;
    err: string = "";
}

export class CreateUserReq {
    operator: number = 0;
    user_name: string = "";
    password: string = "";
    permission: number = 0;
}

export class CreateUserRes {
    is_success: boolean = false;
    err: string = "";
}

export class LockUserReq {
    operator: number = 0;
}

export class LockUserRes {
    is_success: boolean = false;
    err: string = "";
}

export class UnlockUserReq {
    operator: number = 0;
}

export class UnlockUserRes {
    is_success: boolean = false;
    err: string = "";
}

export class ModifyInfoReq {
    operator: number = 0;
    nickname: string = "";
    password: string = "";
    modify_tk_flag: boolean = false;
    totp_key: string = "";
}

export class ModifyInfoRes {
    is_success: boolean = false;
    err: string = "";
}

export class AdjustPermissionReq {
    operator: number = 0;
    user_id: number = 0;
    permission: number = 0;
}

export class AdjustPermissionRes {
    is_success: boolean = false;
    err: string = "";
}

export class AuthenticateReq {
    user_id: number = 0;
    password: string = "";
}

export class AuthenticateRes {
    is_success: boolean = false;
    err: string = "";
}
