import type {NextPage} from 'next'
import {NextPageContext} from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from 'next'
import Link from "next/link"
import { useRouter } from "next/router"
import Cookies from "js-cookie"
import RequireSignIn from "../../components/require_signin"

interface Props {
	dt: string
	data: {} 
}

const Dashboard: NextPage<Props> = (props) => {
	const router = useRouter()
	const logout = () => {
		Cookies.remove("signedIn")
		router.replace("/user/signin")
	}
	
	return (
		<RequireSignIn>
			<div>
				dashboard
			</div>
		</RequireSignIn>
			
	)
}


export const getServerSideProps: GetServerSideProps = async (context) => {
	const res = await fetch(`${process.env.BACKEND}/user/dashboard`)
	const data = await res.json()
	return {props: {data:data}}
}


export default Dashboard 

