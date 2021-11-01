
const RequireSignIn = ({children}) => {
	const cookies = parseCookies()
	console.log({ cookies }) 
	// const cookies = nookies.get(ctx)

	// const signedIn = Cookies.get("signedIn")
	// if (signedIn !== "true") {
	// 	redirect: {
	// 		permanent: false,
	// 		destination: '/user/signin',
	// 	},
	// } 
	return children
	
}

export default RequireSignIn 
