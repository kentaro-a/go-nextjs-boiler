import type {NextPage} from 'next'
import {NextPageContext} from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from 'next'
import Link from "next/link"
import { useRouter } from "next/router"
import RequireSignIn from "../../components/require_signin"

interface Props {
}

const Dashboard: NextPage<Props> = (props) => {
	const router = useRouter()

	const logout = () => {
		router.replace("/user/signin")
	}
	const test = async () => {
		try {
			const res = await fetch(`/api/proxy/user/dashboard`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
			})
			const data = await res.json()
			console.log(data)

		} catch (e) {
			console.log(e)
		}	
	}

	const deleteUser = async () => {
		try {
			const res = await fetch(`/api/proxy/user/delete`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({password: "12345678abc"}),
			})
			const data = await res.json()
			console.log(data)

		} catch (e) {
			console.log(e)
		}	
	}
	
	return (
		<RequireSignIn>
			<div>
				<div>
					dashboard
				</div>

				<div>
					<button type="button" onClick={test}>dashboard</button>
				</div>
				<div>
					<button type="button" onClick={deleteUser}>delete</button>
				</div>
			</div>
		</RequireSignIn>
	)
}


export const getServerSideProps: GetServerSideProps = async (context) => {
	console.log("dashboard")
	// return {
	// 	redirect: {
	//       permanent: false, // 永続的なリダイレクトかどうか
	//       destination: 'https://www.yahoo.co.jp/', // リダイレクト先
	//     },
	// }
		
	return {props: {}}
}


export default Dashboard 

