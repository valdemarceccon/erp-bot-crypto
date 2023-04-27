<script lang="ts">
	import { enhance, type SubmitFunction } from '$app/forms';
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

	const validateForm: SubmitFunction = ({form, data, action, cancel, submitter}) => {
    console.log("chamou");
		let validatedForm = formSchema.safeParse({
			name: name,
			api_key: api_key,
			api_secret: api_secret,
			exchange: exchange
		});

		if (!validatedForm.success) {
      let localErrors = validatedForm.error.format();
      errors.name = localErrors.name?._errors[0];
      errors.api_key = localErrors.api_key?._errors[0];
      errors.api_secret = localErrors.api_secret?._errors[0];
      errors.exchange = localErrors.exchange?._errors[0];
			cancel();
		}
	}
</script>

<form method="POST" use:enhance={validateForm}>
	<div class="flex flex-col gap-2">
		<div class="flex flex-row gap-2">
			<label class="label flex-1">
				<span class="my-5">
          Key Name
				{#if errors?.name}
					<span class="text-error-500 ml-5">
						{errors?.name}
					</span>
				{/if}
      </span>
				<input
					bind:value={name}
					name="name"
					type="text"
					placeholder="Name"
					class="input my-2"
					class:input-error={errors?.name}
				/>
			</label>
			<label class="label flex-1">
				<span>Exchange</span>
				{#if errors?.exchange}
        <span class="text-error-500 ml-5">
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
		<label class="label">
			<span>API Key</span>
			{#if errors?.api_key}
      <span class="text-error-500 ml-5">
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
		<label class="label">
			<span>Secret Key</span>
			{#if errors?.api_secret}
      <span class="text-error-500 ml-5">
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
