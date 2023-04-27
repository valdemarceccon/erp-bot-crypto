<script lang="ts" context="module">
	export let userStore: any;
</script>

<script lang="ts">
	import {
		AppShell,
		AppBar,
		Avatar,
		Drawer,
		drawerStore,
		toastStore
	} from '@skeletonlabs/skeleton';
	import Navigation, { type NavigationItem } from '$lib/components/Navigation.svelte';

	import type { PageData } from './$types';
	import { enhance } from '$app/forms';
	import { page } from '$app/stores';
	import { createSessionStore, type UserToken } from './sessionStore';

	export let data: PageData;

	let navItems: NavigationItem[] = [
		{ label: 'Dashboard', href: '/' },
		{ label: 'Users', href: '/users', permission: 'ListUsers' },
		{ label: 'APIs', href: '/apis' },
		{ label: 'Comissions', href: '/comissions', permission: 'None' },
		{ label: 'Referal link', href: '/referal', permission: 'None' }
	];

	function openDrawer() {
		drawerStore.open();
	}
	let user: UserToken;
	if (data.success) {
		user = data;
		if (!userStore) {
			userStore = createSessionStore({
				email: data.email,
				name: data.name,
				token: data.access_token,
				username: data.username,
				permissions: data.permissions
			});
		}
	}

	$: if (!data.success && data.message) {
		toastStore.trigger({ message: data.message });
	}

	let permitedNavigation = navItems.filter((nav) => {
		if (!nav.permission) return true;

		if (user?.permissions) {
			for (let per of user.permissions) {
				if (nav.permission == per.name) {
					return true;
				}
			}
		}

		return false;
	});

	let noSpecialPermitionRequired: NavigationItem[] = [];
	let specialPermitionRequired: NavigationItem[] = [];

	permitedNavigation.forEach((x) =>
		(x.permission ? specialPermitionRequired : noSpecialPermitionRequired).push(x)
	);

	$: initials =
		user && user.name
			? user.name
					.split(' ')
					.map((n: any) => n[0])
					.join('')
			: '';
</script>

<Drawer>
	<Navigation items={noSpecialPermitionRequired} />
	<hr class="!border-t-2" />
	<Navigation items={specialPermitionRequired} />
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
				<form action="/logout" method="POST" use:enhance>
					<div class="flex">
						<button type="submit" class="btn !bg-transparent">Logout</button>
						<Avatar {initials} background="bg-primary-500" width="w-10" />
					</div>
				</form>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>
	<svelte:fragment slot="sidebarLeft">
		<Navigation items={noSpecialPermitionRequired} title="Menu" />
		{#if specialPermitionRequired.length > 0}
			<hr class="!border-t-2" />
			<Navigation items={specialPermitionRequired} title="Administration"/>
		{/if}
	</svelte:fragment>
	<div class="container md:p-10 mx-auto">
		<slot />
	</div>
</AppShell>

<style>
</style>
