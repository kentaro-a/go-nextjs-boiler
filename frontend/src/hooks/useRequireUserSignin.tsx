import { useEffect } from 'react';
import { useRouter } from 'next/router';
import { BackendProxyRequest, ApiResponse } from '../components/api';
import { useRecoilState, useResetRecoilState, RecoilRoot } from 'recoil'
import { signinUserState } from '../states/signinUserState'

export const useRequireUserSignin = () => {
	const router = useRouter()
	const [signinUser, setSigninUserState] = useRecoilState(signinUserState);
	const resetSigninUserState = useResetRecoilState(signinUserState);

	useEffect(async ()=>{
		const data: ApiResponse = await BackendProxyRequest(`/user/verify_signin`)
		if (data.statusCode === 200) {
			// グローバルステートのuser情報を更新
			setSigninUserState({...signinUser, ...data.data.user})
		} else {
			// グローバルステートのuser情報をクリア
			resetSigninUserState()
			// リダイレクト
			router.push("/user/signin")
		}
	}, [])
}

