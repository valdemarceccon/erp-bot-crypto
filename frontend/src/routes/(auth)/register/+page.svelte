<script lang="ts">
	import LoginLogout from '$lib/components/LoginLogout.svelte';
	import { fade } from 'svelte/transition';
	import { toastStore } from '@skeletonlabs/skeleton';
	import { enhance } from '$app/forms';
	export let form;
	$: error_message = !form || form.ok ? "" : form.detail;
	$: validation_errors = form?.validation;

	$: username = form?.values?.username ? form?.values?.username : "";
	$: email = form?.values?.email ? form?.values?.email : "";
	$: name = form?.values?.name ? form?.values?.name : "";
	$: password = form?.values?.password;
	$: password_confirm = form?.values?.password_confirm;

	$: password_match = password == password_confirm;
	$: {
		if (error_message) {
			toastStore.trigger({message: error_message, 	background: 'variant-filled-error',});
		}
	}

	function clearServerError() {
		if (form) {
			form.ok = true;
			form.detail = "";
		}
	}

</script>

<div class="card">
<form method="POST" use:enhance>
	<header class="card-header flex flex-col">
		<LoginLogout active="register" />
  </header>
	<section class="p-4">
		<label class="label">
			<span>Username</span>
			<input required name="username" class="input" type="text" placeholder="Username" value={username}/>
		</label>
		{#if validation_errors?.username}
			<p class="text-error-800">{validation_errors.username}</p>
		{/if}
		<label class="label">
			<span>Name</span>
			<input required name="name" class="input" type="text" placeholder="Name" value={name}/>
		</label>
		{#if validation_errors?.name}
			<p class="text-error-800">{validation_errors.name}</p>
		{/if}
		<label class="label">
			<span>Email</span>
			<input required name="email" value={email} class="input" type="email" placeholder="Email" />
		</label>
		{#if validation_errors?.email}
			<p class="text-error-800">{validation_errors.email}</p>
		{/if}
		<label class="label">
			<span>Password</span>
			<input bind:value={password} required name="password" class="input" type="password" placeholder="Password" />
		</label>
		{#if validation_errors?.password}
			<p class="text-error-800">{validation_errors.password}</p>
		{/if}
		<label class="label">
			<span>Confirm Password</span>
			<input bind:value={password_confirm} required name="password_confirm" class="input" class:input-error={!password_match} type="password" placeholder="Password" />
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
