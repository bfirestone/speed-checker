<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { apiUrl } from '$lib/config';
	
	export let data: { dashboardData: DashboardData };

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
		edges?: {
			host?: {
				id: number;
				name: string;
				hostname: string;
				type: string;
			};
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
		recent_speed_tests: SpeedTest[];
		recent_iperf_tests: IperfTest[];
		statistics: {
			total_speed_tests: number;
			total_iperf_tests: number;
			active_hosts: number;
		};
	}

	let dashboardData: DashboardData | null = data.dashboardData;
	let loading = false;
	let error: string | null = null;

	// Filter states
	let speedTestFilters = {
		serverName: '',
		showSlowest: false,
		limit: 10
	};

	let iperfFilters = {
		hostName: '',
		hostType: '',
		showSlowest: false,
		limit: 10
	};

	let filteredSpeedTests: SpeedTest[] = [];
	let filteredIperfTests: IperfTest[] = [];

	async function fetchDashboardData() {
		if (!browser) return;
		
		try {
			loading = true;
			const response = await fetch(apiUrl('/api/v1/dashboard'));
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			dashboardData = await response.json();
			// Debug logging
			if (dashboardData && dashboardData.recent_iperf_tests && dashboardData.recent_iperf_tests.length > 0) {
				console.log('Dashboard iperf test sample:', dashboardData.recent_iperf_tests[0]);
			}
			// Apply filters when data is loaded
			applySpeedTestFilters();
			applyIperfFilters();
		} catch (e) {
			error = e instanceof Error ? e.message : 'An error occurred';
			console.error('Failed to fetch dashboard data:', e);
		} finally {
			loading = false;
		}
	}

	async function applySpeedTestFilters() {
		if (!browser) return;
		
		try {
			const params = new URLSearchParams();
			params.append('limit', speedTestFilters.limit.toString());
			
			if (speedTestFilters.serverName.trim()) {
				params.append('server_name', speedTestFilters.serverName.trim());
			}
			
			if (speedTestFilters.showSlowest) {
				params.append('slowest', 'true');
			}

			const response = await fetch(apiUrl(`/api/v1/speedtest/results?${params}`));
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			
			const result = await response.json();
			filteredSpeedTests = result.results || result.data || [];
		} catch (e) {
			console.error('Failed to fetch filtered speed tests:', e);
			filteredSpeedTests = dashboardData?.recent_speed_tests || [];
		}
	}

	async function applyIperfFilters() {
		if (!browser) return;
		
		try {
			const params = new URLSearchParams();
			params.append('limit', iperfFilters.limit.toString());
			
			if (iperfFilters.hostName.trim()) {
				params.append('host_name', iperfFilters.hostName.trim());
			}
			
			if (iperfFilters.hostType.trim()) {
				params.append('host_type', iperfFilters.hostType.trim());
			}
			
			if (iperfFilters.showSlowest) {
				params.append('slowest', 'true');
			}

			const response = await fetch(apiUrl(`/api/v1/iperf/results?${params}`));
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			
			const result = await response.json();
			filteredIperfTests = result.results || result.data || [];
			// Debug logging
			if (filteredIperfTests.length > 0) {
				console.log('Sample iperf test:', filteredIperfTests[0]);
			}
		} catch (e) {
			console.error('Failed to fetch filtered iperf tests:', e);
			filteredIperfTests = dashboardData?.recent_iperf_tests || [];
		}
	}

	// Reactive statements to apply filters when they change
	$: if (browser && dashboardData) {
		// Watch for changes in speed test filter values
		speedTestFilters.serverName, speedTestFilters.showSlowest, speedTestFilters.limit;
		applySpeedTestFilters();
	}

	$: if (browser && dashboardData) {
		// Watch for changes in iperf filter values
		iperfFilters.hostName, iperfFilters.hostType, iperfFilters.showSlowest, iperfFilters.limit;
		applyIperfFilters();
	}

	async function runSpeedTest() {
		try {
			const response = await fetch(apiUrl('/api/v1/legacy/speedtest/run'), { method: 'POST' });
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			// Refresh data
			await fetchDashboardData();
		} catch (e) {
			console.error('Failed to run speed test:', e);
		}
	}

	async function runIperfTests() {
		try {
			const response = await fetch(apiUrl('/api/v1/legacy/iperf/run'), { method: 'POST' });
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			// Refresh data
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

	function resetSpeedTestFilters() {
		speedTestFilters = {
			serverName: '',
			showSlowest: false,
			limit: 10
		};
	}

	function resetIperfFilters() {
		iperfFilters = {
			hostName: '',
			hostType: '',
			showSlowest: false,
			limit: 10
		};
	}

	async function deleteSpeedTest(testId: number) {
		if (!confirm('Are you sure you want to delete this speed test result?')) {
			return;
		}

		try {
			const response = await fetch(apiUrl(`/api/v1/speedtest/results/${testId}`), { 
				method: 'DELETE' 
			});
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			// Refresh data
			await fetchDashboardData();
		} catch (e) {
			console.error('Failed to delete speed test:', e);
			alert('Failed to delete speed test. Please try again.');
		}
	}

	async function deleteIperfTest(testId: number) {
		if (!confirm('Are you sure you want to delete this iperf test result?')) {
			return;
		}

		try {
			const response = await fetch(apiUrl(`/api/v1/iperf/results/${testId}`), { 
				method: 'DELETE' 
			});
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			// Refresh data
			await fetchDashboardData();
		} catch (e) {
			console.error('Failed to delete iperf test:', e);
			alert('Failed to delete iperf test. Please try again.');
		}
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
	{:else if dashboardData && dashboardData.statistics}
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
								<dd class="text-lg font-medium text-gray-900">{dashboardData.statistics.total_speed_tests}</dd>
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
								<dd class="text-lg font-medium text-gray-900">{dashboardData.statistics.total_iperf_tests}</dd>
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
								<dd class="text-lg font-medium text-gray-900">{dashboardData.statistics.active_hosts}</dd>
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

		<!-- Filters and Recent Tests -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Speed Tests -->
			<div class="bg-white shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<div class="flex justify-between items-center mb-4">
						<h3 class="text-lg leading-6 font-medium text-gray-900">Speed Tests</h3>
						<button
							on:click={resetSpeedTestFilters}
							class="text-sm text-gray-500 hover:text-gray-700"
						>
							Reset Filters
						</button>
					</div>
					
					<!-- Speed Test Filters -->
					<div class="mb-4 space-y-3 p-4 bg-gray-50 rounded-lg">
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
							<div>
								<label for="speed-server-name" class="block text-sm font-medium text-gray-700 mb-1">Server Name</label>
								<input
									id="speed-server-name"
									type="text"
									bind:value={speedTestFilters.serverName}
									placeholder="Search by server name..."
									class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								/>
							</div>
							<div>
								<label for="speed-limit" class="block text-sm font-medium text-gray-700 mb-1">Limit</label>
								<select
									id="speed-limit"
									bind:value={speedTestFilters.limit}
									class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								>
									<option value={5}>5 results</option>
									<option value={10}>10 results</option>
									<option value={25}>25 results</option>
									<option value={50}>50 results</option>
								</select>
							</div>
						</div>
						<div class="flex items-center">
							<input
								type="checkbox"
								bind:checked={speedTestFilters.showSlowest}
								id="speed-slowest"
								class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
							/>
							<label for="speed-slowest" class="ml-2 block text-sm text-gray-700">
								Show slowest tests (sorted by download speed)
							</label>
						</div>
					</div>

					{#if filteredSpeedTests.length > 0}
						<div class="space-y-4">
							{#each filteredSpeedTests as test}
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
										<div class="flex items-center space-x-2">
											<p class="text-xs text-gray-400">{formatTimestamp(test.timestamp)}</p>
											<button
												on:click={() => deleteSpeedTest(test.id)}
												class="text-red-500 hover:text-red-700 p-1 rounded transition-colors"
												title="Delete this speed test result"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
												</svg>
											</button>
										</div>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-gray-500">No speed tests found</p>
					{/if}
				</div>
			</div>

			<!-- Iperf Tests -->
			<div class="bg-white shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<div class="flex justify-between items-center mb-4">
						<h3 class="text-lg leading-6 font-medium text-gray-900">Iperf Tests</h3>
						<button
							on:click={resetIperfFilters}
							class="text-sm text-gray-500 hover:text-gray-700"
						>
							Reset Filters
						</button>
					</div>
					
					<!-- Iperf Test Filters -->
					<div class="mb-4 space-y-3 p-4 bg-gray-50 rounded-lg">
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
							<div>
								<label for="iperf-host-name" class="block text-sm font-medium text-gray-700 mb-1">Host Name</label>
								<input
									id="iperf-host-name"
									type="text"
									bind:value={iperfFilters.hostName}
									placeholder="Search by host name..."
									class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								/>
							</div>
							<div>
								<label for="iperf-host-type" class="block text-sm font-medium text-gray-700 mb-1">Host Type</label>
								<select
									id="iperf-host-type"
									bind:value={iperfFilters.hostType}
									class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								>
									<option value="">All Types</option>
									<option value="lan">LAN</option>
									<option value="vpn">VPN</option>
									<option value="remote">Remote</option>
								</select>
							</div>
						</div>
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
							<div>
								<label for="iperf-limit" class="block text-sm font-medium text-gray-700 mb-1">Limit</label>
								<select
									id="iperf-limit"
									bind:value={iperfFilters.limit}
									class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
								>
									<option value={5}>5 results</option>
									<option value={10}>10 results</option>
									<option value={25}>25 results</option>
									<option value={50}>50 results</option>
								</select>
							</div>
							<div class="flex items-center pt-6">
								<input
									type="checkbox"
									bind:checked={iperfFilters.showSlowest}
									id="iperf-slowest"
									class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
								/>
								<label for="iperf-slowest" class="ml-2 block text-sm text-gray-700">
									Show slowest tests
								</label>
							</div>
						</div>
					</div>

					{#if filteredIperfTests.length > 0}
						<div class="space-y-4">
							{#each filteredIperfTests as test}
								<div class="border-l-4 {test.success ? 'border-green-400' : 'border-red-400'} pl-4">
									<div class="flex justify-between items-start">
										<div>
											{#if test.success}
												<p class="text-sm font-medium text-gray-900">
													↑ {formatSpeed(test.sent_mbps)} Mbps / ↓ {formatSpeed(test.received_mbps)} Mbps
												</p>
												<p class="text-sm text-gray-500">
													{#if test.edges?.host?.name}
														{test.edges.host.name} ({test.edges.host.type.toUpperCase()})
													{:else}
														Host: Unknown
													{/if}
													{#if test.mean_rtt_ms} • RTT: {formatSpeed(test.mean_rtt_ms)}ms{/if}
													{#if test.retransmits && test.retransmits > 0} • Retransmits: {test.retransmits}{/if}
												</p>
											{:else}
												<p class="text-sm font-medium text-red-500">Test failed</p>
												<p class="text-sm text-gray-500">
													{#if test.edges?.host?.name}
														{test.edges.host.name} ({test.edges.host.type.toUpperCase()})
													{:else}
														Host: Unknown
													{/if}
												</p>
											{/if}
										</div>
										<div class="flex items-center space-x-2">
											<p class="text-xs text-gray-400">{formatTimestamp(test.timestamp)}</p>
											<button
												on:click={() => deleteIperfTest(test.id)}
												class="text-red-500 hover:text-red-700 p-1 rounded transition-colors"
												title="Delete this iperf test result"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
												</svg>
											</button>
										</div>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-gray-500">No iperf tests found</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Active Hosts -->
		<div class="mt-8 bg-white shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Active Hosts</h3>
				<div class="text-center py-8">
					<p class="text-gray-500 mb-4">Host information is available on the dedicated hosts page</p>
					<a href="/hosts" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
						View Hosts
					</a>
				</div>
			</div>
		</div>
	{:else}
		<div class="bg-yellow-50 border border-yellow-200 rounded-md p-4">
			<div class="text-center">
				<h3 class="text-sm font-medium text-yellow-800">No dashboard data available</h3>
				<p class="mt-2 text-sm text-yellow-700">Please try refreshing the page or check that the API server is running.</p>
				<button
					on:click={() => window.location.reload()}
					class="mt-3 inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-yellow-700 bg-yellow-100 hover:bg-yellow-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
				>
					Refresh Page
				</button>
			</div>
		</div>
	{/if}
</div>
