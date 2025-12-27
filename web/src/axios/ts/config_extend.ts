import {AxiosResponse, InternalAxiosRequestConfig} from "axios";
import {axiosWrapper} from "@/axios/ts/config.ts";
import {log} from "@/ts/log";

const HttpHeader_UserID = "Unnamed-Plan-User-ID"
const HttpHeader_AccessToken = "Unnamed-Plan-Access-Token"
const StorageName_UserID = "user_id"
const StorageName_AccessToken = "access_token"

export function initInterceptors(invalidLoginHandler: () => void): void {
    axiosWrapper.interceptors.request.use(
        (value: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
            value.headers.set(HttpHeader_UserID, localStorage.getItem(StorageName_UserID))
            value.headers.set(HttpHeader_AccessToken, localStorage.getItem(StorageName_AccessToken))

            return value
        }
    )

    axiosWrapper.interceptors.response.use(
        (value: AxiosResponse): AxiosResponse => {

            //@ts-ignore
            const userID: string = value.headers.get(HttpHeader_UserID)
            if (userID && userID.length > 0) {
                localStorage.setItem(StorageName_UserID, userID)
            }

            //@ts-ignore
            const accessToken: string = value.headers.get(HttpHeader_AccessToken)
            if (accessToken && accessToken.length > 0) {
                localStorage.setItem(StorageName_AccessToken, accessToken)
            }

            return value
        },
        (error: any) => {
            let code = error.response.status

            if (code == 401) { // 验证失败，用户id或访问密钥错误
                invalidLoginHandler()
            }

            log.httpError(code)

            return Promise.reject(error)
        }
    )
}
