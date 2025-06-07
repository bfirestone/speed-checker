import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';

// API configuration - use public env var for browser, fallback for development
export const API_BASE_URL = env.PUBLIC_API_BASE_URL || 'http://localhost:8080';

// Helper function to build API URLs
export function apiUrl(path: string): string {
	// In browser: use configured public URL
	// During SSR: this will be overridden by server-side fetch with internal URLs
	return `${API_BASE_URL}${path}`;
} 