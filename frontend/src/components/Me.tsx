

import type {NextPage} from 'next'
import {React} from 'react'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import Link from "next/link"
import { useRouter } from "next/router"
import { useRequireUserSignin } from "../hooks/useRequireUserSignin"
import { useRecoilState, RecoilRoot } from 'recoil'
import { signinUserState } from '../states/signinUserState'

import { atom } from 'recoil';

interface Props {
}

type Form = {
	name: string 
}
const formState = atom<Form>({
	key: 'formState',
	default: {
		"name": "aji",
	},
})


const Me = () => {
	
	const [form, setForm] = useRecoilState(formState);

	const [signinUser, setSigninUser] = useRecoilState(signinUserState);
	
	const update = () => {
		setSigninUser({...signinUser, name: form.name})
	}
	const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setForm({name: e.target.value})
	}
	return (
		<div>
			<div style={{backgroudColor:'red'}}>
				My name is {signinUser?.name}
			</div>
			<div><button onClick={update}>update</button></div>
			<div><input onChange={onChange} value={form.name} /></div>
		</div>
	)
}


export default Me

