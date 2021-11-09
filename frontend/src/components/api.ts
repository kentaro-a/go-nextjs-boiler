import type { NextApiRequest, NextApiResponse } from "next"
import path from "path"

export type ApiData = {
	messages: string[],
	data: {},
} 

export const NewApiData = (messages: string[], data: {}): ApiData => {
	const ret: ApiData = {
		messages: messages,
		data: data,
	}
	return ret 
}

export type ApiStatus = {
	statusCode: number,
}

export const NewApiStatus = (statusCode: number): ApiStatus => {
	const ret: ApiStatus = {
		statusCode: statusCode,
	}
	return ret 
}

export type ApiResponse = ApiStatus & ApiData 

export const NewApiResponse = (s: ApiStatus, d: ApiData): ApiResponse => {
	return {
		...s,
		...d,
	}
}

const _request = async (url: string, options: {}): Promise<ApiResponse> => {
	const res = await fetch(url, options)
	const _json = await res.json() 
	const data: {} = _json.data || {}
	const messages: string[] = _json.messages || []
	const apiData: ApiData = NewApiData(messages, data)
	const apiStatus: ApiStatus = NewApiStatus(_json.status_code)
	const apiRes = NewApiResponse(apiStatus, apiData)
	return apiRes	
	
}


/**
 * For Server Side request
 *
 */
export const BackendRequest = async (_path: string, data: {} = {}): Promise<ApiResponse> => {
	const headers = {
		'Content-Type': 'application/json',
	}
	const options = {
		method: "POST",
		headers: new Headers(headers),
		body: JSON.stringify(data),
	}	
	const url = process.env.BACKEND + path.join(``, _path)
	return _request(url, options)
} 

/**
 * For Client Side request
 *
 */
export const BackendProxyRequest = async (_path: string, data: {} = {}): Promise<ApiResponse> => {
	const headers = {
		'Content-Type': 'application/json',
	}
	const options = {
		method: "POST",
		headers: new Headers(headers),
		body: JSON.stringify(data),
	}	
	
	const url = path.join("/api/proxy/", _path)
	return _request(url, options)
	
}


// export const BackendRequestMultiFormData = async (url: string, data: {} = {}) => {
	
// } 



