<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	interface Host {
		id: number;
		name: string;
		hostname: string;
		port: number;
		type: string;
		active: boolean;
		description?: string;
	}

	let hosts: Host[] = [];
	let loading = true;
	let error: string | null = null;
	let showAddForm = false;
	let editingHost: Host | null = null;

	// Form data
	let formData = {
		name: '',
		hostname: '',
		port: 5201,
		type: 'lan',
		description: '',
		active: true
	};

	async function fetchHosts() {
		if (!browser) return;
		
		try {
			loading = true;
			const response = await fetch('/api/v1/hosts');
			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}
			const data = await response.json();
			hosts = data.data || [];
		} catch (e) {
			error = e instanceof Error ? e.message : 'An error occurred';
			console.error('Failed to fetch hosts:', e);
		} finally {
			loading = false;
		}
	}

	async function addHost() {
		try {
			const response = await fetch('/api/v1/hosts', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(formData),
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			// Reset form and refresh hosts
			resetForm();
			await fetchHosts();
		} catch (e) {
			console.error('Failed to add host:', e);
			alert('Failed to add host: ' + (e instanceof Error ? e.message : 'Unknown error'));
		}
	}

	async function updateHost() {
		if (!editingHost) return;
		
		try {
			const response = await fetch(`/api/v1/hosts/${editingHost.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(formData),
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			// Reset form and refresh hosts
			resetForm();
			await fetchHosts();
		} catch (e) {
			console.error('Failed to update host:', e);
			alert('Failed to update host: ' + (e instanceof Error ? e.message : 'Unknown error'));
		}
	}

	async function deleteHost(hostId: number, hostName: string) {
		if (!confirm(`Are you sure you want to delete "${hostName}"?`)) {
			return;
		}

		try {
			const response = await fetch(`/api/v1/hosts/${hostId}`, {
				method: 'DELETE',
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			await fetchHosts();
		} catch (e) {
			console.error('Failed to delete host:', e);
			alert('Failed to delete host: ' + (e instanceof Error ? e.message : 'Unknown error'));
		}
	}

	function startEdit(host: Host) {
		editingHost = host;
		formData = {
			name: host.name,
			hostname: host.hostname,
			port: host.port,
			type: host.type,
			description: host.description || '',
			active: host.active
		};
		showAddForm = true;
	}

	function resetForm() {
		formData = {
			name: '',
			hostname: '',
			port: 5201,
			type: 'lan',
			description: '',
			active: true
		};
		showAddForm = false;
		editingHost = null;
	}

	function handleSubmit() {
		if (editingHost) {
			updateHost();
		} else {
			addHost();
		}
	}

	onMount(() => {
		fetchHosts();
	});
</script>

<svelte:head>
	<title>Host Management - Speed Checker</title>
</svelte:head>

<div class="px-4 py-6 sm:px-0">
	<!-- Header -->
	<div class="mb-8 flex justify-between items-center">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Host Management</h1>
			<p class="mt-2 text-gray-600">Manage iperf3 test hosts for network performance testing</p>
		</div>
		<button
			on:click={() => { resetForm(); showAddForm = true; }}
			class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
		>
			Add Host
		</button>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
		</div>
	{:else if error}
		<div class="bg-red-50 border border-red-200 rounded-md p-4">
			<div class="flex">
				<div class="ml-3">
					<h3 class="text-sm font-medium text-red-800">Error loading hosts</h3>
					<div class="mt-2 text-sm text-red-700">
						<p>{error}</p>
					</div>
				</div>
			</div>
		</div>
	{:else}
		<!-- Add/Edit Host Form -->
		{#if showAddForm}
			<div class="mb-8 bg-white shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
						{editingHost ? 'Edit Host' : 'Add New Host'}
					</h3>
					<form on:submit|preventDefault={handleSubmit} class="space-y-4">
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div>
								<label for="name" class="block text-sm font-medium text-gray-700">Name</label>
								<input
									type="text"
									id="name"
									bind:value={formData.name}
									required
									class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
									placeholder="e.g., Home Router"
								/>
							</div>
							<div>
								<label for="hostname" class="block text-sm font-medium text-gray-700">Hostname/IP</label>
								<input
									type="text"
									id="hostname"
									bind:value={formData.hostname}
									required
									class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
									placeholder="e.g., 192.168.1.1"
								/>
							</div>
							<div>
								<label for="port" class="block text-sm font-medium text-gray-700">Port</label>
								<input
									type="number"
									id="port"
									bind:value={formData.port}
									required
									min="1"
									max="65535"
									class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
							<div>
								<label for="type" class="block text-sm font-medium text-gray-700">Type</label>
								<select
									id="type"
									bind:value={formData.type}
									required
									class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
								>
									<option value="lan">LAN</option>
									<option value="vpn">VPN</option>
									<option value="remote">Remote</option>
								</select>
							</div>
						</div>
						<div>
							<label for="description" class="block text-sm font-medium text-gray-700">Description (Optional)</label>
							<textarea
								id="description"
								bind:value={formData.description}
								rows="2"
								class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
								placeholder="Optional description of this host"
							></textarea>
						</div>
						{#if editingHost}
							<div class="flex items-center">
								<input
									type="checkbox"
									id="active"
									bind:checked={formData.active}
									class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
								/>
								<label for="active" class="ml-2 block text-sm text-gray-700">
									Active (include in testing)
								</label>
							</div>
						{/if}
						<div class="flex justify-end space-x-3">
							<button
								type="button"
								on:click={resetForm}
								class="bg-gray-300 hover:bg-gray-400 text-gray-700 font-medium py-2 px-4 rounded-md transition-colors"
							>
								Cancel
							</button>
							<button
								type="submit"
								class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
							>
								{editingHost ? 'Update Host' : 'Add Host'}
							</button>
						</div>
					</form>
				</div>
			</div>
		{/if}

		<!-- Hosts List -->
		{#if hosts.length > 0}
			<div class="bg-white shadow rounded-lg overflow-hidden">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Configured Hosts</h3>
					<div class="overflow-x-auto">
						<table class="min-w-full divide-y divide-gray-200">
							<thead class="bg-gray-50">
								<tr>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Name
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Hostname
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Port
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Type
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Status
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Description
									</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
										Actions
									</th>
								</tr>
							</thead>
							<tbody class="bg-white divide-y divide-gray-200">
								{#each hosts as host}
									<tr>
										<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
											{host.name}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
											{host.hostname}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
											{host.port}
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
												{host.type === 'lan' ? 'bg-blue-100 text-blue-800' : 
												 host.type === 'vpn' ? 'bg-purple-100 text-purple-800' : 
												 'bg-gray-100 text-gray-800'}">
												{host.type.toUpperCase()}
											</span>
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
												{host.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
												{host.active ? 'Active' : 'Inactive'}
											</span>
										</td>
										<td class="px-6 py-4 text-sm text-gray-500">
											{host.description || '-'}
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
											<div class="flex space-x-2">
												<button
													on:click={() => startEdit(host)}
													class="text-blue-600 hover:text-blue-900 bg-blue-50 hover:bg-blue-100 px-2 py-1 rounded text-xs font-medium transition-colors"
												>
													Edit
												</button>
												<button
													on:click={() => deleteHost(host.id, host.name)}
													class="text-red-600 hover:text-red-900 bg-red-50 hover:bg-red-100 px-2 py-1 rounded text-xs font-medium transition-colors"
												>
													Delete
												</button>
											</div>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		{:else}
			<div class="text-center py-12">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900">No hosts configured</h3>
				<p class="mt-1 text-sm text-gray-500">Get started by adding your first iperf3 host.</p>
				<div class="mt-6">
					<button
						on:click={() => { resetForm(); showAddForm = true; }}
						class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
					>
						Add Host
					</button>
				</div>
			</div>
		{/if}
	{/if}
</div> 