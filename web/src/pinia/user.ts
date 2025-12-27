import {defineStore} from "pinia";
import {CreateUserRes, ListUserRes, LoginRes, ModifyUserRes, User} from "@/axios/ts/user.go.ts";
import {ref} from "vue";
import CryptoJs from "crypto-js"
import {userAxios} from "@/axios/ts/user.http.ts";
import {log} from "@/ts/log.ts";

export let useUserStore = defineStore("user", () => {
    let user = ref<User>(new User())

    function register(userName: string, password: string, cb: () => void): void {
        let pwdHash = CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex)

        userAxios.createUser(userName, pwdHash)
            .then(({data}: { data: CreateUserRes }) => {
                if (!data.is_success) {
                    log.fail("register", data.err)
                    return
                }

                cb()

                log.success("register")

                login(userName, password, "", () => {
                })
            })
    }

    function login(userName: string, password: string, totpCode: string, cb: () => void): void {
        let pwdHash = CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex)

        userAxios.login(userName, pwdHash, totpCode)
            .then(({data}: { data: LoginRes }) => {
                if (!data.is_success) {
                    log.fail("login", data.err)
                    return
                }

                user.value = loginResToUser(data)

                cb()

                log.success("login")
            })
    }

    function modify(nickname: string, password: string, modifyTkFlag: boolean, totpKey: string): void {
        let pwdHash = CryptoJs.SHA256(password).toString(CryptoJs.enc.Hex)

        userAxios.modifyUser(nickname, password, modifyTkFlag, totpKey)
            .then(({data}: { data: ModifyUserRes }) => {
                if (!data.is_success) {
                    log.fail("modify user", data.err)
                    return
                }

                log.success("modify user")
            })
    }

    function list(pageSize: number, pageNum: number, cb: (amount: number, users: Array<User>) => void): void {
        userAxios.listUser({size: pageSize, num: pageNum})
            .then(({data}: { data: ListUserRes }) => {
                if (!data.is_success) {
                    log.fail("list user", data.err)
                    return
                }

                cb(data.amount, data.users)

                log.success("list user")
            })
    }

    function loginResToUser(res: LoginRes): User {
        let userIns = new User()
        userIns.id = res.user_id
        userIns.user_name = res.user_name
        userIns.nickname = res.nickname
        userIns.is_admin = res.is_admin

        return userIns
    }

    function isLogin(): boolean {
        return user.value.user_name.length > 0
    }

    function exitLogin(): void {
        user.value = new User()
    }

    return {user, register, login, modify, list, isLogin, exitLogin}
})
