import { env } from '$env/dynamic/private';

// Server-side API configuration using internal Docker network
const SERVER_API_BASE_URL = env.API_BASE_URL || 'http://speed-checker-api:8080';

// Server-side API URL builder
export function serverApiUrl(path: string): string {
	return `${SERVER_API_BASE_URL}${path}`;
}

// Server-side fetch wrapper
export async function serverFetch(path: string, options?: RequestInit): Promise<Response> {
	const url = serverApiUrl(path);
	return fetch(url, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...options?.headers
		}
	});
} 