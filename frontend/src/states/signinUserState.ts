import { atom } from 'recoil';

// undefined : まだログイン確認が完了していない状態とする
// null      : ログイン確認をした結果、ログインしていなかった状態とする
export type SigninUser = {
	id: number
	name: string 
	mail: string 
	password: string 
}

export const signinUserState = atom<SigninUser>({
	key: 'signinUserState',
	default: null,
})
