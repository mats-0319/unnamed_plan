// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.0

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { LoginRes, LoginReq, ListUserRes, ListUserReq, CreateUserRes, CreateUserReq, ModifyUserRes, ModifyUserReq, AuthenticateRes } from "./user.go"
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

    public listUser(page: Pagination): Promise<AxiosResponse<ListUserRes>> {
        let req: ListUserReq = {
            page: page,
        }

        return axiosWrapper.post("/user/list", req)
    }

    public createUser(user_name: string, password: string): Promise<AxiosResponse<CreateUserRes>> {
        let req: CreateUserReq = {
            user_name: user_name,
            password: password,
        }

        return axiosWrapper.post("/user/create", req)
    }

    public modifyUser(nickname: string, password: string, modify_tk_flag: boolean, totp_key: string): Promise<AxiosResponse<ModifyUserRes>> {
        let req: ModifyUserReq = {
            nickname: nickname,
            password: password,
            modify_tk_flag: modify_tk_flag,
            totp_key: totp_key,
        }

        return axiosWrapper.post("/user/modify", req)
    }

    public Authenticate(): Promise<AxiosResponse<AuthenticateRes>> {
        return axiosWrapper.post("/user/authenticate")
    }
}

export const userAxios: UserAxios = new UserAxios()
