import type {NextPage} from 'next'
import {NextPageContext} from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from 'next'
import { useRouter } from "next/router"
import React from 'react'
import { BackendProxyRequest, ApiResponse } from '../../components/api'

import { useRecoilState, RecoilRoot, atom } from 'recoil'
import { signinUserState } from '../../states/signinUserState'


 
const signinFormState = atom({
	key: 'signinFormState',
	default: {
		// "mail": "",
		// "password": "",
		"mail": "user1@test.com",
		"password": "12345678abc",
	},
})

const SignIn: NextPage = (props) => {
	const router = useRouter()
	const [signinForm, setSigninForm] = useRecoilState(signinFormState);

	const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const value = e.target.value
		const name = e.target.name
		setSigninForm({...signinForm, [name]: value})
	}

	const signin = async () => {
		try {
			const data: ApiResponse = await BackendProxyRequest(`/user/signin`, signinForm)
			if (data.statusCode === 200) {
				router.push("/user/dashboard")
			}
			
		} catch (e) {
			console.log(e)
		}	
	}


	return (
		<div>
			signin <br />
			<div>
				<div>
					mail: <input type="text" name="mail" value={signinForm.mail} onChange={onChange} />
				</div>
				<div>
					password: <input type="password" name="password" value={signinForm.password} onChange={onChange}/>
				</div>
			</div>
			<div style={{marginTop:"2em"}}>
				Input Data<br />
				<table border="1">
					<tr><td>{signinForm.mail}</td><td>{signinForm.password}</td></tr>
				</table>
			</div>
			<div>
				<button type="button" onClick={signin}>SignIn</button>
			</div>
		</div>	
	)
}



export default SignIn 
