import '../styles/globals.css'
import type {AppProps} from 'next/app'

function MyApp({Component, pageProps}: AppProps) {
	const [user, setUser] = useState(null);


	return <Component {...pageProps} />
}

export default MyApp
