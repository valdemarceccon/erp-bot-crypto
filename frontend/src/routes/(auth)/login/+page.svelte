<script lang="ts">
	import { fade } from 'svelte/transition';
	export let form;
	$: error_message = !form || form.ok ? "" : form.detail;
</script>

{#if error_message}
<aside transition:fade|local={{ duration: 200 }} class="alert variant-filled-error">		<!-- Icon -->
		<div class="alert-message">
				<h3>Login error</h3>
				<p>{error_message}</p>
		</div>
		<div class="alert-actions"><button on:click={() => {form = {ok: true, detail: ""}}}>x</button></div>
</aside>
{/if}

<div class="card">
<form method="POST">
	<header class="card-header flex flex-col">
    Login
  </header>
	<section class="p-4">
		<label class="label">
			<span>Username</span>
			<input required name="username" class="input" type="text" placeholder="Username" />
		</label>
		<label class="label">
			<span>Password</span>
			<input type="password" required name="password" class="input" placeholder="Password" />
		</label>
	</section>
	<footer class="card-footer flex flex-col">

		<button type="submit" class="btn variant-filled-primary flex-auto">Login</button>
	</footer>
</form>
</div>
