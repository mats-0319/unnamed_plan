// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.0

import axios, {AxiosInstance, AxiosResponse, InternalAxiosRequestConfig} from "axios";

export const axiosWrapper: AxiosInstance = axios.create({
    baseURL: import.meta.env.Vite_axios_base_url,
    timeout: 3000,
});

const HttpHeader_UserID = "Unnamed-Plan-User-ID"
const HttpHeader_AccessToken = "Unnamed-Plan-Access-Token"
const StorageName_UserID = "user_id"
const StorageName_AccessToken = "access_token"
const AuthenticateError = "Authentication Error"

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
            if (value.data && value.data.err && (value.data.err as string).includes(AuthenticateError)) {
                // 验证失败，用户id或访问密钥错误
                invalidLoginHandler()
                return value
            }

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
        }
    )
}
