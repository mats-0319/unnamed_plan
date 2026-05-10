// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { Pagination } from "./common.go"
import { RegisterRes, RegisterReq, LoginRes, LoginReq, LoginMFARes, LoginMFAReq, ListUserRes, ListUserReq, ModifyUserRes, ModifyUserReq } from "./user.go"

class UserAxios {
    public register(user_name: string, password: string): Promise<AxiosResponse<RegisterRes>> {
        let req: RegisterReq = {
            user_name: user_name,
            password: password,
        }

        return axiosWrapper.post("/register", req)
    }

    public login(user_name: string, password: string): Promise<AxiosResponse<LoginRes>> {
        let req: LoginReq = {
            user_name: user_name,
            password: password,
        }

        return axiosWrapper.post("/login", req)
    }

    public loginMFA(mfa_token: string, totp_code: string): Promise<AxiosResponse<LoginMFARes>> {
        let req: LoginMFAReq = {
            mfa_token: mfa_token,
            totp_code: totp_code,
        }

        return axiosWrapper.post("/login-mfa", req)
    }

    public listUser(page: Pagination): Promise<AxiosResponse<ListUserRes>> {
        let req: ListUserReq = {
            page: page,
        }

        return axiosWrapper.post("/user/list", req)
    }

    public modifyUser(nickname: string, password: string, enable_mfa: boolean, totp_key: string): Promise<AxiosResponse<ModifyUserRes>> {
        let req: ModifyUserReq = {
            nickname: nickname,
            password: password,
            enable_mfa: enable_mfa,
            totp_key: totp_key,
        }

        return axiosWrapper.post("/user/modify", req)
    }
}

export const userAxios: UserAxios = new UserAxios()
