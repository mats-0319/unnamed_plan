// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v1.0.0

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { LoginRes, LoginReq, ListUserRes, ListUserReq, CreateUserRes, CreateUserReq, LockUserRes, LockUserReq, UnlockUserRes, UnlockUserReq, ModifyInfoRes, ModifyInfoReq, AdjustPermissionRes, AdjustPermissionReq, AuthenticateRes, AuthenticateReq } from "./user.go"
import { Pagination } from "./common.go"

class UserAxios {
    public login(user_name: string, password: string, totp_code: string): Promise<AxiosResponse<LoginRes>> {
        let req: LoginReq = {
            user_name: user_name,
            password: password,
            totp_code: totp_code,
        }

        return axiosWrapper.post("/login", req)
    }

    public listUser(operator: number, page: Pagination, see_locked: boolean): Promise<AxiosResponse<ListUserRes>> {
        let req: ListUserReq = {
            operator: operator,
            page: page,
            see_locked: see_locked,
        }

        return axiosWrapper.post("/user/list", req)
    }

    public createUser(operator: number, user_name: string, password: string, permission: number): Promise<AxiosResponse<CreateUserRes>> {
        let req: CreateUserReq = {
            operator: operator,
            user_name: user_name,
            password: password,
            permission: permission,
        }

        return axiosWrapper.post("/user/create", req)
    }

    public lockUser(operator: number): Promise<AxiosResponse<LockUserRes>> {
        let req: LockUserReq = {
            operator: operator,
        }

        return axiosWrapper.post("/user/lock", req)
    }

    public unlockUser(operator: number): Promise<AxiosResponse<UnlockUserRes>> {
        let req: UnlockUserReq = {
            operator: operator,
        }

        return axiosWrapper.post("/user/unlock", req)
    }

    public modifyInfo(operator: number, nickname: string, password: string, modify_tk_flag: boolean, totp_key: string): Promise<AxiosResponse<ModifyInfoRes>> {
        let req: ModifyInfoReq = {
            operator: operator,
            nickname: nickname,
            password: password,
            modify_tk_flag: modify_tk_flag,
            totp_key: totp_key,
        }

        return axiosWrapper.post("/user/modify-info", req)
    }

    public AdjustPermission(operator: number, user_id: number, permission: number): Promise<AxiosResponse<AdjustPermissionRes>> {
        let req: AdjustPermissionReq = {
            operator: operator,
            user_id: user_id,
            permission: permission,
        }

        return axiosWrapper.post("/user/adjust-permission", req)
    }

    public Authenticate(user_id: number, password: string): Promise<AxiosResponse<AuthenticateRes>> {
        let req: AuthenticateReq = {
            user_id: user_id,
            password: password,
        }

        return axiosWrapper.post("/user/authenticate", req)
    }
}

export const userAxios: UserAxios = new UserAxios()
