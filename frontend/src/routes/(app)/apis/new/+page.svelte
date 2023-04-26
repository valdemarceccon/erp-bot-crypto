<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { toastStore } from '@skeletonlabs/skeleton';
	import { redirect } from '@sveltejs/kit';
	import { z } from 'zod';

	export let form;

	let message: string;

	let name: string = '';
	let api_key: string = '';
	let api_secret: string = '';
	let exchange: string = '';

  let errors: {
    name?: string,
    api_key?: string,
    api_secret?: string,
    exchange?: string,
  } = {};

	let formSchema = z.object({
		name: z.string().trim().nonempty('Field is required.'),
		api_key: z.string().trim().nonempty('Field is required.'),
		api_secret: z.string().trim().nonempty('Field is required.'),
		exchange: z.string().trim().nonempty('Field is required.')
	});

	$: if (form?.ok) {
		toastStore.trigger({ message: 'Saved with success', background: 'variant-filled-success' });
		goto('/apis');
	}

	$: if (form && form.detail) {
		toastStore.trigger({ message: form.detail, background: 'variant-filled-error' });
	}

	function validateForm({form, data, action, cancel, submitter}) {
    console.log("chamou");
		let validatedForm = formSchema.safeParse({
			name: name,
			api_key: api_key,
			api_secret: api_secret,
			exchange: exchange
		});

		if (!validatedForm.success) {
      let valerrors = validatedForm.error.format();
      errors.name = valerrors.name?._errors[0]
      errors.api_key = valerrors.api_key?._errors[0]
      errors.api_secret = valerrors.api_secret?._errors[0]
      errors.exchange = valerrors.exchange?._errors[0]
			cancel();
		}
	}
</script>

<form method="POST" use:enhance={validateForm}>
	<div class="flex flex-col gap-2">
		<div class="flex gap-2">
			<label class="w-6/12">
				Key Name
				{#if errors?.name}
					<span class="badge variant-filled-error">
						{errors?.name}
					</span>
				{/if}
				<input
					bind:value={name}
					name="name"
					type="text"
					placeholder="Name"
					class="input my-2"
					class:input-error={errors?.name}
				/>
			</label>
			<label class="w-6/12">
				Exchange
				{#if errors?.exchange}
					<span class="badge variant-filled-error">
						{errors?.exchange}
					</span>
				{/if}
				<input
					bind:value={exchange}
					name="exchange"
					type="text"
					placeholder="Exchange"
					class="input my-2"
					class:input-error={errors?.exchange}
				/>
			</label>
		</div>
		<label>
			API Key
			{#if errors?.api_key}
				<span class="badge variant-filled-error">
					{errors?.api_key}
				</span>
			{/if}
			<input
				bind:value={api_key}
				name="api_key"
				type="text"
				placeholder="API key"
				class="input my-2"
				class:input-error={errors?.api_key}
			/>
		</label>
		<label>
			Secret Key
			{#if errors?.api_secret}
				<span class="badge variant-filled-error">
					{errors?.api_secret}
				</span>
			{/if}
			<input
				bind:value={api_secret}
				name="api_secret"
				type="text"
				placeholder="Secret Key"
				class="input my-2"
				class:input-error={errors?.api_secret}
			/>
		</label>

		<button class="btn variant-filled-primary mt-5" type="submit"> Create </button>
	</div>

	<!-- {message} -->
</form>
