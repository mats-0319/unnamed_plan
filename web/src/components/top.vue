<template>
	<div class="color-bg-1 center-hv">
		<div class="top center-hv">
			<div class="t-title" @click="routerLink('home')">Unnamed Plan Web</div>

			<div v-if="flags.wildScreenFlag" class="t-content">
				<div class="tc-item center-hv" @click="routerLink('note')">小纸条</div>

				<div class="tc-item center-hv">
					<a href="https://github.com/mats0319/unnamed_plan" target="_blank">本站代码</a>
				</div>

				<div class="tc-item center-hv">
					<outlined-button v-show="!userStore.isLogin()" @click="beforeOpenLoginDialog()"
						>登录</outlined-button
					>
					<div v-show="userStore.isLogin()">
						<el-dropdown placement="bottom-end">
							<template #default>
								<span> {{ userStore.user.nickname }}&nbsp;<span class="icon">&or;</span> </span>
							</template>

							<template #dropdown>
								<el-dropdown-menu>
									<el-dropdown-item @click="routerLink('pDefault')">个人中心</el-dropdown-item>
									<el-dropdown-item divided @click="exitLogin()">退出登录</el-dropdown-item>
								</el-dropdown-menu>
							</template>
						</el-dropdown>
					</div>
				</div>
			</div>

			<div v-else class="t-content" @click="routerLink('note')">小纸条</div>
		</div>
	</div>

	<el-dialog v-model="showLoginDialog" title="登录">
		<el-form v-model="loginReq" label-width="20%">
			<el-form-item label="用户名">
				<el-input v-model="loginReq.user_name" />
			</el-form-item>

			<el-form-item label="密码">
				<el-input v-model="loginReq.password" show-password />
			</el-form-item>

			<el-form-item label="TOTP Code">
				<el-input v-model="loginReq.totp_code" />
			</el-form-item>

			<el-form-item>
				<outlined-button details="请输入用户名和密码后点击注册" :disabled="!canLoginFlag" @click="register">
					注册新用户
				</outlined-button>
			</el-form-item>
		</el-form>

		<template #footer>
			<el-button @click="showLoginDialog = false">取消</el-button>
			<el-button type="primary" :disabled="!canLoginFlag" @click="login()">登录</el-button>
		</template>
	</el-dialog>
</template>

<script lang="ts" setup>
import { useUserStore } from "@/pinia/user.ts"
import { ref, watch } from "vue"
import { useFlagStore } from "@/pinia/flag.ts"
import OutlinedButton from "@/components/outlined_button.vue"
import { routerLink } from "@/ts/util.ts"
import { LoginReq } from "@/axios/ts/user.go.ts"

let flags = useFlagStore()
let userStore = useUserStore()

let showLoginDialog = ref<boolean>(false)
let canLoginFlag = ref<boolean>(false)
let loginReq = ref<LoginReq>(new LoginReq())

function beforeOpenLoginDialog(): void {
	loginReq.value = new LoginReq()
	showLoginDialog.value = true
}

function login(): void {
	userStore.login(loginReq.value.user_name, loginReq.value.password, loginReq.value.totp_code, () => {
		showLoginDialog.value = false
		routerLink("pDefault")
	})
}

function register(): void {
	userStore.register(loginReq.value.user_name, loginReq.value.password, () => {
		showLoginDialog.value = false
	})
}

function exitLogin(): void {
	userStore.exitLogin()
	routerLink("home")
	localStorage.removeItem("user_id")
	localStorage.removeItem("access_token")
	sessionStorage.removeItem("login_data")
}

watch(
	loginReq,
	(newValue, _) => {
		canLoginFlag.value = newValue.user_name.length > 0 && newValue.password.length > 0
	},
	{ deep: true }
)
</script>

<style lang="less" scoped>
.top {
	@media (min-width: 1280px) {
		width: 80rem;
	}
	width: 40rem;
	height: 6.25rem;

	.t-title {
		@media (min-width: 1280px) {
			width: 24rem;
			padding: 0 3rem;
			font-size: 2.2rem;
		}
		width: 50%;
		font-size: 1.6rem;
	}

	.t-title:hover {
		cursor: pointer;
	}

	.t-content {
		display: flex;
		justify-content: flex-end;
		@media (min-width: 1280px) {
			width: 50rem;
		}
		width: 50%;

		.tc-item {
			width: 8rem;
			padding: 0 1rem;
			font-size: 1.2rem;

			a {
				color: black;
				text-decoration-line: none;
			}

			.el-dropdown {
				font-size: 1.6rem;

				.icon {
					font-size: 70%;
				}
			}
		}

		.tc-item:hover {
			cursor: pointer;
		}
	}
}
</style>
