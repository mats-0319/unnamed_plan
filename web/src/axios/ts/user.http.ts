// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.4

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import {
	RegisterRes,
	RegisterReq,
	LoginRes,
	LoginReq,
	ListUserRes,
	ListUserReq,
	ModifyUserRes,
	ModifyUserReq
} from "./user.go"
import { Pagination } from "./common.go"

class UserAxios {
	public register(user_name: string, password: string): Promise<AxiosResponse<RegisterRes>> {
		let req: RegisterReq = {
			user_name: user_name,
			password: password
		}

		return axiosWrapper.post("/register", req)
	}

	public login(user_name: string, password: string, totp_code: string): Promise<AxiosResponse<LoginRes>> {
		let req: LoginReq = {
			user_name: user_name,
			password: password,
			totp_code: totp_code
		}

		return axiosWrapper.post("/login", req)
	}

	public listUser(page: Pagination): Promise<AxiosResponse<ListUserRes>> {
		let req: ListUserReq = {
			page: page
		}

		return axiosWrapper.post("/user/list", req)
	}

	public modifyUser(
		nickname: string,
		password: string,
		enable_2fa: boolean,
		totp_key: string
	): Promise<AxiosResponse<ModifyUserRes>> {
		let req: ModifyUserReq = {
			nickname: nickname,
			password: password,
			enable_2fa: enable_2fa,
			totp_key: totp_key
		}

		return axiosWrapper.post("/user/modify", req)
	}
}

export const userAxios: UserAxios = new UserAxios()
