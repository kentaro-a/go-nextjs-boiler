import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import { GetServerSideProps } from 'next'
import { BackendRequest, ApiResponse } from '../components/api'

const HealthCheck: NextPage = (props) => {
	return (<></>)
}

export const getServerSideProps: GetServerSideProps = async (context) => {
	try {
		const data: ApiResponse = await BackendRequest(`/health_check`, {})
		context.res.statusCode = data.status_code

	} catch (e) {
		context.res.statusCode = 500
	}

	return {props: {}}
}

export default HealthCheck 
