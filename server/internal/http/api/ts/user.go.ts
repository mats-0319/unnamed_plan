// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v1.0.0

import { ResBase, Pagination } from "./common.go"

export class LoginReq {
    user_name: string = "";
    password: string = "";
    totp_code: string = "";
}

export class LoginRes {
    res: ResBase = new ResBase();
    token: string = "";
}

export class User {
    id: number = 0;
    created_at: number = 0;
    updated_at: number = 0;
    name: string = "";
    nickname: string = "";
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
    res: ResBase = new ResBase();
    user_amount: number = 0;
    list_amount: number = 0;
    users: Array<User> = new Array<User>();
}

export class CreateUserReq {
    operator: number = 0;
    user_name: string = "";
    password: string = "";
    permission: number = 0;
}

export class CreateUserRes {
    res: ResBase = new ResBase();
}

export class LockUserReq {
    operator: number = 0;
}

export class LockUserRes {
    res: ResBase = new ResBase();
}

export class UnlockUserReq {
    operator: number = 0;
}

export class UnlockUserRes {
    res: ResBase = new ResBase();
}

export class ModifyInfoReq {
    operator: number = 0;
    nickname: string = "";
    password: string = "";
    modify_tk_flag: boolean = false;
    totp_key: string = "";
}

export class ModifyInfoRes {
    res: ResBase = new ResBase();
}

export class AdjustPermissionReq {
    operator: number = 0;
    user_id: number = 0;
    permission: number = 0;
}

export class AdjustPermissionRes {
    res: ResBase = new ResBase();
}

export class AuthenticateReq {
    user_id: number = 0;
    password: string = "";
}

export class AuthenticateRes {
    res: ResBase = new ResBase();
}
