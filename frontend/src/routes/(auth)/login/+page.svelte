<script lang="ts">
	import { enhance } from '$app/forms';
	import LoginLogout from '$lib/components/LoginLogout.svelte';
	import { toastStore } from '@skeletonlabs/skeleton';
	import type { Snapshot } from './$types.js';

	let loginValues = {
		username: '',
		password: ''
	};

	export const snapshot: Snapshot = {
		capture: () => loginValues,
		restore: (value) => (loginValues = value)
	};

	export let form;
	$: error_message = !form || form.ok ? '' : form.message;
	$: {
		if (error_message) {
			toastStore.trigger({ message: error_message, background: 'variant-filled-error' });
		}
	}
</script>

<div class="card">
	<form
		method="POST"
		use:enhance={() => {
			return async ({ update }) => {
				await update({ reset: false });
			};
		}}
	>
		<header class="card-header flex flex-col">
			<LoginLogout active="login" />
		</header>
		<section class="p-4">
			<label class="label">
				<span>Username</span>
				<input
					required
					name="username"
					class="input"
					type="text"
					placeholder="Username"
					bind:value={loginValues.username}
				/>
			</label>
			<label class="label">
				<span>Password</span>
				<input
					type="password"
					required
					name="password"
					class="input"
					placeholder="Password"
					bind:value={loginValues.password}
				/>
			</label>
		</section>
		<footer class="card-footer flex flex-col">
			<button type="submit" class="btn variant-filled-primary flex-auto">Login</button>
		</footer>
	</form>
</div>
