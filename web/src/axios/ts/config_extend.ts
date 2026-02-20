import { AxiosResponse, InternalAxiosRequestConfig } from "axios"
import { axiosWrapper } from "@/axios/ts/config.ts"
import { log } from "@/ts/log"

const HttpHeader_AccessToken = "Unnamed-Plan-Access-Token"
const StorageName_AccessToken = "access_token"

export function initInterceptors(invalidLoginHandler: () => void): void {
	axiosWrapper.interceptors.request.use((value: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
		value.headers.set(HttpHeader_AccessToken, localStorage.getItem(StorageName_AccessToken))

		return value
	})

	axiosWrapper.interceptors.response.use(
		(value: AxiosResponse) => {
			//@ts-ignore
			const accessToken: string = value.headers.get(HttpHeader_AccessToken)
			if (accessToken && accessToken.length > 0) {
				localStorage.setItem(StorageName_AccessToken, accessToken)
			}

			if (!value.data.is_success) {
				log.fail(value.config.url as string, value.data.err)
				return Promise.reject({ isBusinessError: true })
			}

			return value.data
		},
		(error: any) => {
			if (!!error.isBusinessError) {
				return Promise.reject()
			}

			const code: number = error.response?.status

			if (code == 401) {
				// 验证失败，用户id或访问密钥错误
				invalidLoginHandler()
			}

			log.httpError(code)

			return Promise.reject(error)
		}
	)
}

// export interface AxiosResponse<T = any, D = any, H = {}> {
//   data: T;
//   status: number;
//   statusText: string;
//   headers: H & RawAxiosResponseHeaders | AxiosResponseHeaders;
//   config: InternalAxiosRequestConfig<D>;
//   request?: any;
// }
//
// 根据axios类型定义，我们在res拦截器里拿到的value.data就是我们定义的返回结构：
// type Response struct {
// 	IsSuccess bool   `json:"is_success"`
// 	Err       string `json:"err"`
// 	Data      any    `json:"data"`
// }
//
// 因为res拦截器要求输入/输出参数类型一致，所以不容易直接返回请求的有效属性(value.data.data)
// - `type AxiosResponseInterceptorUse<T> = (onFulfilled?: ((value: T) => T | Promise<T>)`
// 所以解析请求的返回值需要这样写：`({ data }: { data: LoginRes })`，而很难直接写`(data: LoginRes)`
