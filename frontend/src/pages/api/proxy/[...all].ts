import { NextApiRequest, NextApiResponse } from "next";
import httpProxyMiddleware from 'next-http-proxy-middleware'
import https from 'https'

export const config = {
	api: {
		bodyParser: false
	},
}
// NextJSの/api/proxy/*へのリクエストをbackendにプロキシする
export default (req: NextApiRequest, res: NextApiResponse): Promise<any> => {
	const proxy = httpProxyMiddleware(req, res, {
		target: process.env.BACKEND,
		changeOrigin: true,
		headers: {
		},
		pathRewrite: [
			{
				patternStr: '^/api/proxy',
				replaceStr: ''
			},
		],
	})
	return proxy
}
