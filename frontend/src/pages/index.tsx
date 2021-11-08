import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {GetServerSideProps} from 'next'

const Home: NextPage = (props) => {
	return (
		<div>index</div>	
	)
}

export const getServerSideProps: GetServerSideProps = async (context) => {
	return {props: {}}
}

export default Home
