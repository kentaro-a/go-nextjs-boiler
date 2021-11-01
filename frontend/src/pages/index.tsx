import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from 'next'

const Home: NextPage = (props) => {
	console.log(props)
	return (
		<div>index</div>	
	)
}

export const getServerSideProps: GetServerSideProps = async (context) => {
	const res = await fetch(`${process.env.BACKEND}/user/dashboard`)
	const data = await res.json()
	return {props: data}
}

export default Home
