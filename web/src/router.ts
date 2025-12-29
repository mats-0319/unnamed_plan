import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router"
import { useUserStore } from "@/pinia/user.ts"

const routes: Array<RouteRecordRaw> = [
	{
		path: "/",
		name: "home",
		component: () => import("@/views/home.vue")
	},
	{
		path: "/note",
		name: "note",
		component: () => import("@/views/note.vue")
	},
	{
		path: "/personal-center",
		name: "personalCenter",
		meta: { requireLogin: true },
		component: () => import("@/views/personal_center/home.vue"),
		children: [
			{
				path: "",
				name: "pDefault",
				component: () => import("@/views/personal_center/default.vue")
			},
			{
				path: "list-user",
				name: "pListUser",
				component: () => import("@/views/personal_center/list_user.vue")
			},
			{
				path: "modify-user",
				name: "pModifyUser",
				component: () => import("@/views/personal_center/modify_user.vue")
			},
			{
				path: "note",
				name: "pNote",
				component: () => import("@/views/personal_center/my_note.vue")
			}
		]
	},
	{
		path: "/404", // 考虑后续可能编写404页面，此处预留路由
		name: "notFound",
		redirect: { name: "home" }
	},
	{
		path: "/:pathMatch(.*)*", // 将匹配所有内容并将其放在 `$route.params.pathMatch` 下
		redirect: { name: "notFound" }
	}
]

export const router = createRouter({
	history: createWebHashHistory(),
	routes: routes
})

router.beforeEach((to, _from, next) => {
	let userStore = useUserStore()

	if (!(to.meta && to.meta.requireLogin) || userStore.isLogin()) {
		// 页面不需要登录，或者已经登录
		next()
		return
	}

	next({ path: "/" })
	return
})
