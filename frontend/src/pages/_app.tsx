import '../styles/globals.css'
import type {AppProps} from 'next/app'
import { RecoilRoot } from 'recoil'
import { signInUserState } from '../states/signInUserState'
import { ApiResponse, BackendProxyRequest } from '../components/api'
import RecoilDebugObserver from '../components/RecoilDebugObserver'


function MyApp({Component, pageProps}: AppProps) {
	return (
		<RecoilRoot>
			{ process.env.NODE_ENV === "development" && <RecoilDebugObserver /> }
			<Component {...pageProps} />
		</RecoilRoot>
	) 
}

export default MyApp
