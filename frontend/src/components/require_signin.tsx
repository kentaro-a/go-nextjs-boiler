import {GetServerSideProps} from 'next'


interface Props {
	children?: React.ReactNode
}

const RequireSignIn = ({children}: Props) => {
	return (
		<>
			{children}
		</>
	)
}


export default RequireSignIn 
