<script lang="ts">
	import LoginLogout from '$lib/components/LoginLogout.svelte';
	import { fade } from 'svelte/transition';
	import { toastStore } from '@skeletonlabs/skeleton';
	import { enhance } from '$app/forms';
	export let form;
	$: error_message = !form || form.ok ? '' : form.message;
	$: validation_errors = form?.validation;

	let formData = {
		username: '',
		fullname: '',
		email: '',
		password: '',
		password_confirm: ''
	};

	$: password_match = formData.password == formData.password_confirm;
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
				update({ reset: false });
			};
		}}
	>
		<header class="card-header flex flex-col">
			<LoginLogout active="register" />
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
					value={formData.username}
				/>
			</label>
			{#if validation_errors?.username}
				<p class="text-error-800">{validation_errors.username}</p>
			{/if}
			<label class="label">
				<span>Name</span>
				<input
					required
					name="name"
					class="input"
					type="text"
					placeholder="Name"
					value={formData.fullname}
				/>
			</label>
			{#if validation_errors?.name}
				<p class="text-error-800">{validation_errors.name}</p>
			{/if}
			<label class="label">
				<span>Email</span>
				<input
					required
					name="email"
					value={formData.email}
					class="input"
					type="email"
					placeholder="Email"
				/>
			</label>
			{#if validation_errors?.email}
				<p class="text-error-800">{validation_errors.email}</p>
			{/if}
			<label class="label">
				<span>Password</span>
				<input
					required
					name="password"
					class="input"
					type="password"
					placeholder="Password"
					bind:value={formData.password}
				/>
			</label>
			{#if validation_errors?.password}
				<p class="text-error-800">{validation_errors.password}</p>
			{/if}
			<label class="label">
				<span>Confirm Password</span>
				<input
					required
					name="password_confirm"
					class="input"
					class:input-error={!password_match}
					type="password"
					placeholder="Password"
					bind:value={formData.password_confirm}
				/>
			</label>
			{#if !password_match}
				<p>Password should match</p>
			{/if}
		</section>
		<footer class="card-footer flex flex-col">
			<button type="submit" class="btn variant-filled-primary flex-auto">Register</button>
		</footer>
	</form>
</div>
