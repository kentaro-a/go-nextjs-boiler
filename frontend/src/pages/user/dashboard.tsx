import type {NextPage} from 'next'
import {NextPageContext} from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import Link from "next/link"
import { useRouter } from "next/router"
import { useRequireUserSignin } from "../../hooks/useRequireUserSignin"
import { useRecoilState, useResetRecoilState, RecoilRoot } from 'recoil'
import { BackendProxyRequest, ApiResponse } from "../../components/api"
import { signinUserState } from '../../states/signinUserState'
import Me from '../../components/Me'



const Dashboard: NextPage = () => {
	
	useRequireUserSignin()

	const router = useRouter()
	const [signinUser, setSigninUser] = useRecoilState(signinUserState)
	const resetSigninUserState = useResetRecoilState(signinUserState)

	const signout = async () => {
		try {
			const data: ApiResponse = await BackendProxyRequest(`/user/signout`)
			if (data.statusCode === 200) {
				resetSigninUserState()
				router.reload()
			}
			
		} catch (e) {
			console.log(e)
		}

	}

	const deleteUser = async () => {
		try {
			const data = await BackendProxyRequest(`/user/delete`, {
				password: "12345678abc",
			})
			console.log(data)

		} catch (e) {
			console.log(e)
		}	
	}
	
	return (
		<div>
			<div>
				dashboard
			</div>
			<Me />
			<div>
				<button type="button" onClick={signout}>Sign out</button>
			</div>
			<div>
				<button type="button" onClick={deleteUser}>delete</button>
			</div>
		</div>
	)
}





export default Dashboard 

