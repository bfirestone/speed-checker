import type { ServerLoad } from '@sveltejs/kit';
import { serverFetch } from '$lib/server/api';

export const load: ServerLoad = async () => {
	try {
		const response = await serverFetch('/api/v1/dashboard');
		
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		
		const dashboardData = await response.json();
		console.log('SSR Dashboard data:', JSON.stringify(dashboardData, null, 2));
		
		return {
			dashboardData
		};
	} catch (error) {
		console.error('Failed to fetch dashboard data during SSR:', error);
		// Return empty data structure so the page can still load and fetch client-side
		return {
			dashboardData: {
				recent_iperf_tests: [],
				recent_speed_tests: [],
				statistics: {
					active_hosts: 0,
					total_iperf_tests: 0,
					total_speed_tests: 0
				}
			}
		};
	}
}; 