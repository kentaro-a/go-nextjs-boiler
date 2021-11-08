import type {NextPage} from 'next'
import {GetStaticProps} from 'next'


interface Props {
	title: string
}


const SSG: NextPage<Props> = ({title}: Props) => {
	return (
		<div>
			{title}
		</div>	
	)
}

export const getStaticProps: GetStaticProps = async (context) => {
	return {props: {"title": "SSG"}}
}

export default SSG 

