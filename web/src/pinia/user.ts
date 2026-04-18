import { defineStore } from "pinia"
import { ListUserRes, LoginRes, LoginMFARes, ModifyUserRes, RegisterRes, User } from "@/axios/ts/user.go.ts"
import { ref } from "vue"
import CryptoJs from "crypto-js"
import { userAxios } from "@/axios/ts/user.http.ts"
import { log } from "@/ts/log.ts"

export const useUserStore = defineStore("user", () => {
    const user = ref<User>(new User())

    function register(userName: string, password: string, cb: () => void): void {
        userAxios.register(userName, calcSHA256(password)).then(({}: { data: RegisterRes }) => {
            cb()

            log.success("register")

            login(userName, password, () => {
            })
        })
    }

    function login(userName: string, password: string, cb: (mfaToken: string) => void): void {
        userAxios.login(userName, calcSHA256(password)).then(({ data }: { data: LoginRes }) => {
            cb(data.mfa_token)

            if (!data.enable_mfa) { // disable MFA
                user.value = loginResToUser(data)
                log.success("login")
            }
        })
    }

    function loginMFA(mfaToken: string, totpCode: string, cb: () => void): void {
        userAxios.loginMFA(mfaToken, totpCode).then(({ data }: { data: LoginMFARes }) => {
            user.value = loginResToUser(data)

            cb()

            log.success("login (enable MFA)")
        })
    }

    function modify(nickname: string, password: string, enable2FA: boolean, totpKey: string): void {
        userAxios.modifyUser(nickname, calcSHA256(password), enable2FA, totpKey).then(({}: { data: ModifyUserRes }) => {
            log.success("modify user")
        })
    }

    function list(pageSize: number, pageNum: number, cb: (count: number, users: Array<User>) => void): void {
        userAxios.listUser({ size: pageSize, num: pageNum }).then(({ data }: { data: ListUserRes }) => {
            cb(data.count, data.users)

            log.success("list user")
        })
    }

    function loginResToUser(res: LoginRes | LoginMFARes): User {
        const userIns = new User()
        userIns.user_name = res.user_name
        userIns.nickname = res.nickname
        userIns.is_admin = res.is_admin
        userIns.enable_mfa = res.enable_mfa

        return userIns
    }

    function isLogin(): boolean {
        return user.value.user_name.length > 0
    }

    function exitLogin(): void {
        user.value = new User()
        localStorage.removeItem("access_token")
        localStorage.removeItem("login_data")
    }

    function calcSHA256(password: string):string{ // internal func
        return password.length > 0 ? CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex) : ""
    }

    return { user, register, login, loginMFA, modify, list, isLogin, exitLogin }
})
