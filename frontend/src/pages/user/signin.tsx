import type {NextPage} from 'next'
import {NextPageContext} from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from 'next'
import Cookies from "js-cookie"
import { useRouter } from "next/router"
import { useState } from 'react'

interface Props {
	dt: string
	data: {} 
}

const SignIn: NextPage<Props> = (props) => {
	const router = useRouter()

	const [values, setValues] = useState({
		mail: '',
		password: '',
	})

	const handleChange = (e) => {
		const value = e.target.value
		const name = e.target.name
		console.log(value, name)
		setValues({...values, [name]: value})
	}

	const signIn = async () => {
		// Cookies.set("signedIn","true")
		// router.replace("/user/dashboard")
		try {
			const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND}/user/signin`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },	
				body: JSON.stringify(values),
			})
			console.log(res)
		} catch (e) {

			console.log(e)
		}	
		
		// const data = await res.json()
		// console.log(data)
	}

	return (
		<div>
			signin <br />
			<div>
				<div>
					mail: <input type="text" name="mail" value={values.mail} onChange={handleChange} />
				</div>
				<div>
					password: <input type="text" name="password" value={values.password} onChange={handleChange}/>
				</div>
			</div>
			<div>
				<button type="button" onClick={signIn}>SignIn</button>
			</div>
		</div>	
	)
}



export default SignIn 
