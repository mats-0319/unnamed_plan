import { defineStore } from "pinia"
import { ListUserRes, LoginRes, ModifyUserRes, RegisterRes, User } from "@/axios/ts/user.go.ts"
import { ref } from "vue"
import CryptoJs from "crypto-js"
import { userAxios } from "@/axios/ts/user.http.ts"
import { log } from "@/ts/log.ts"

export const useUserStore = defineStore("user", () => {
    const user = ref<User>(new User())

    function register(userName: string, password: string, cb: () => void): void {
        const pwdHash = CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex)

        userAxios.register(userName, pwdHash).then(({}: { data: RegisterRes }) => {
            cb()

            log.success("register")

            login(userName, password, () => {})
        })
    }

    function login(userName: string, password: string, cb: () => void): void {
        const pwdHash = CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex)

        userAxios.login(userName, pwdHash).then(({ data }: { data: LoginRes }) => {
            user.value = loginResToUser(data)

            cb()

            log.success("login")
        })
    }

    function modify(nickname: string, password: string, enable2FA: boolean, totpKey: string): void {
        const pwdHash = CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex)

        userAxios.modifyUser(nickname, pwdHash, enable2FA, totpKey).then(({}: { data: ModifyUserRes }) => {
            log.success("modify user")
        })
    }

    function list(pageSize: number, pageNum: number, cb: (amount: number, users: Array<User>) => void): void {
        userAxios.listUser({ size: pageSize, num: pageNum }).then(({ data }: { data: ListUserRes }) => {
            cb(data.amount, data.users)

            log.success("list user")
        })
    }

    function loginResToUser(res: LoginRes): User {
        const userIns = new User()
        userIns.user_name = res.user_name
        userIns.nickname = res.nickname
        userIns.is_admin = res.is_admin
        userIns.enable_2fa = res.enable_2fa

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

    return { user, register, login, modify, list, isLogin, exitLogin }
})
