<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	interface SpeedTest {
		id: number;
		timestamp: string;
		download_mbps: number;
		upload_mbps: number;
		ping_ms: number;
		jitter_ms?: number;
		server_name?: string;
		isp?: string;
	}

	interface IperfTest {
		id: number;
		timestamp: string;
		sent_mbps: number;
		received_mbps: number;
		retransmits?: number;
		mean_rtt_ms?: number;
		success: boolean;
		host?: {
			name: string;
			hostname: string;
			type: string;
		};
	}

	interface Host {
		id: number;
		name: string;
		hostname: string;
		type: string;
		active: boolean;
	}

	interface DashboardData {
		speed_tests: SpeedTest[];
		iperf_tests: IperfTest[];
		hosts: Host[];
		summary: {
			total_speed_tests: number;
			total_iperf_tests: number;
			active_hosts: number;
		};
	}

	let dashboardData: DashboardData | null = null;
	let loading = true;
	let error: string | null = null;

	async function fetchDashboardData() {
		if (!browser) return;
		
		try {
			loading = true;
			const response = await fetch('/api/v1/dashboard');
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			dashboardData = await response.json();
		} catch (e) {
			error = e instanceof Error ? e.message : 'An error occurred';
			console.error('Failed to fetch dashboard data:', e);
		} finally {
			loading = false;
		}
	}

	async function runSpeedTest() {
		try {
			const response = await fetch('/api/v1/speedtest/run', { method: 'POST' });
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			// Refresh dashboard data
			await fetchDashboardData();
		} catch (e) {
			console.error('Failed to run speed test:', e);
		}
	}

	async function runIperfTests() {
		try {
			const response = await fetch('/api/v1/iperf/run', { method: 'POST' });
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			// Refresh dashboard data
			await fetchDashboardData();
		} catch (e) {
			console.error('Failed to run iperf tests:', e);
		}
	}

	function formatTimestamp(timestamp: string): string {
		return new Date(timestamp).toLocaleString();
	}

	function formatSpeed(speed: number): string {
		return speed.toFixed(2);
	}

	onMount(() => {
		fetchDashboardData();
		// Refresh data every 30 seconds
		const interval = setInterval(fetchDashboardData, 30000);
		return () => clearInterval(interval);
	});
</script>

<svelte:head>
	<title>Speed Checker Dashboard</title>
</svelte:head>

<div class="px-4 py-6 sm:px-0">
	<!-- Header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-900">Network Performance Dashboard</h1>
		<p class="mt-2 text-gray-600">Monitor your internet speed and network performance</p>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
		</div>
	{:else if error}
		<div class="bg-red-50 border border-red-200 rounded-md p-4">
			<div class="flex">
				<div class="ml-3">
					<h3 class="text-sm font-medium text-red-800">Error loading dashboard</h3>
					<div class="mt-2 text-sm text-red-700">
						<p>{error}</p>
					</div>
				</div>
			</div>
		</div>
	{:else if dashboardData}
		<!-- Summary Cards -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
			<div class="bg-white overflow-hidden shadow rounded-lg">
				<div class="p-5">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<div class="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
								<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
								</svg>
							</div>
						</div>
						<div class="ml-5 w-0 flex-1">
							<dl>
								<dt class="text-sm font-medium text-gray-500 truncate">Speed Tests</dt>
								<dd class="text-lg font-medium text-gray-900">{dashboardData.summary.total_speed_tests}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>

			<div class="bg-white overflow-hidden shadow rounded-lg">
				<div class="p-5">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<div class="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
								<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
								</svg>
							</div>
						</div>
						<div class="ml-5 w-0 flex-1">
							<dl>
								<dt class="text-sm font-medium text-gray-500 truncate">Iperf Tests</dt>
								<dd class="text-lg font-medium text-gray-900">{dashboardData.summary.total_iperf_tests}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>

			<div class="bg-white overflow-hidden shadow rounded-lg">
				<div class="p-5">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<div class="w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center">
								<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"></path>
								</svg>
							</div>
						</div>
						<div class="ml-5 w-0 flex-1">
							<dl>
								<dt class="text-sm font-medium text-gray-500 truncate">Active Hosts</dt>
								<dd class="text-lg font-medium text-gray-900">{dashboardData.summary.active_hosts}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="mb-8 flex space-x-4">
			<button
				on:click={runSpeedTest}
				class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
			>
				Run Speed Test
			</button>
			<button
				on:click={runIperfTests}
				class="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
			>
				Run Iperf Tests
			</button>
		</div>

		<!-- Recent Tests -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Speed Tests -->
			<div class="bg-white shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Recent Speed Tests</h3>
					{#if dashboardData.speed_tests.length > 0}
						<div class="space-y-4">
							{#each dashboardData.speed_tests.slice(0, 5) as test}
								<div class="border-l-4 border-blue-400 pl-4">
									<div class="flex justify-between items-start">
										<div>
											<p class="text-sm font-medium text-gray-900">
												↓ {formatSpeed(test.download_mbps)} Mbps / ↑ {formatSpeed(test.upload_mbps)} Mbps
											</p>
											<p class="text-sm text-gray-500">
												Ping: {formatSpeed(test.ping_ms)}ms
												{#if test.server_name} • {test.server_name}{/if}
											</p>
										</div>
										<p class="text-xs text-gray-400">{formatTimestamp(test.timestamp)}</p>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-gray-500">No speed tests available</p>
					{/if}
				</div>
			</div>

			<!-- Iperf Tests -->
			<div class="bg-white shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Recent Iperf Tests</h3>
					{#if dashboardData.iperf_tests.length > 0}
						<div class="space-y-4">
							{#each dashboardData.iperf_tests.slice(0, 5) as test}
								<div class="border-l-4 {test.success ? 'border-green-400' : 'border-red-400'} pl-4">
									<div class="flex justify-between items-start">
										<div>
											<p class="text-sm font-medium text-gray-900">
												{#if test.host}{test.host.name} ({test.host.type}){/if}
											</p>
											{#if test.success}
												<p class="text-sm text-gray-500">
													↑ {formatSpeed(test.sent_mbps)} Mbps / ↓ {formatSpeed(test.received_mbps)} Mbps
													{#if test.mean_rtt_ms} • RTT: {formatSpeed(test.mean_rtt_ms)}ms{/if}
												</p>
											{:else}
												<p class="text-sm text-red-500">Test failed</p>
											{/if}
										</div>
										<p class="text-xs text-gray-400">{formatTimestamp(test.timestamp)}</p>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-gray-500">No iperf tests available</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Active Hosts -->
		{#if dashboardData.hosts.length > 0}
			<div class="mt-8 bg-white shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Active Hosts</h3>
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
						{#each dashboardData.hosts as host}
							<div class="border rounded-lg p-4">
								<div class="flex items-center justify-between">
									<h4 class="font-medium text-gray-900">{host.name}</h4>
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
										{host.type === 'lan' ? 'bg-blue-100 text-blue-800' : 
										 host.type === 'vpn' ? 'bg-purple-100 text-purple-800' : 
										 'bg-gray-100 text-gray-800'}">
										{host.type.toUpperCase()}
									</span>
								</div>
								<p class="text-sm text-gray-500 mt-1">{host.hostname}</p>
							</div>
						{/each}
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>
