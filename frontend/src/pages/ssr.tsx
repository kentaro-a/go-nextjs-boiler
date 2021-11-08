import type {NextPage} from 'next'
import {GetServerSideProps} from 'next'



interface Props {
	title: string
}


const SSR: NextPage<Props> = ({title}: Props) => {
	return (
		<div>
			{title}
		</div>	
	)
}

export const getServerSideProps: GetServerSideProps = async (context) => {
	return {props: {"title": "SSR"}}
}

export default SSR 
