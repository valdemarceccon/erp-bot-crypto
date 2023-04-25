<script lang="ts">
	import { AppShell, AppBar, Avatar, Drawer, drawerStore } from '@skeletonlabs/skeleton';
	import Navigation, { type NavigationItem } from '$lib/components/Navigation.svelte';
	import type { UserToken } from '$lib/stores';

	import type { PageData } from './$types';

	export let data: PageData;

	let navItems: NavigationItem[] = [
		{ label: 'Dashboard', href: '/' },
		{ label: 'Users', href: '/users' }
	];

	function openDrawer() {
		drawerStore.open();
	}

	let user: UserToken = {
		email: data.email,
		name: data.name,
		token: data.access_token,
		username: data.username
	};

	$: initials = user && user.name ? user.name.split(" ").map((n)=>n[0]).join("") : "";


</script>

	<Drawer>
		<Navigation items={navItems} />
	</Drawer>

	<AppShell slotSidebarLeft="w-0 md:w-52 bg-surface-500/10">
		<svelte:fragment slot="header">
			<AppBar>
				<svelte:fragment slot="lead">
					<button class="md:hidden btn btn-sm mr-4" on:click={openDrawer}>
						<span>
							<svg viewBox="0 0 100 80" class="fill-token w-4 h-4">
								<rect width="100" height="20" />
								<rect y="30" width="100" height="20" />
								<rect y="60" width="100" height="20" />
							</svg>
						</span>
					</button>
					<strong class="text-xl uppercase">Bot Erp</strong>
				</svelte:fragment>
				<svelte:fragment slot="trail">
					<form action="/logout" method="POST">
						<div class="flex">
							<button type="submit" class="btn !bg-transparent">Logout</button>
							<Avatar initials={initials} background="bg-primary-500" width="w-10" />
						</div>
					</form>
				</svelte:fragment>
			</AppBar>
		</svelte:fragment>
		<svelte:fragment slot="sidebarLeft">
			<Navigation items={navItems} />
		</svelte:fragment>
		<!-- <svelte:fragment slot="sidebarRight">Sidebar Right</svelte:fragment> -->
		<!-- <svelte:fragment slot="pageHeader">Page Header</svelte:fragment> -->
		<!-- Router Slot -->
		<div class="container p-10 mx-auto">
			{#if !data.success}
				<h2>{data.message}</h2>
			{/if}
			<slot />
		</div>

	</AppShell>

<style>

</style>
