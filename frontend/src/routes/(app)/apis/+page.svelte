<script lang="ts">
	import { enhance } from '$app/forms';
	import { toastStore } from '@skeletonlabs/skeleton';
	import type { PageData } from './$types';
	export let form;

	$: if (form && form.message) {
		toastStore.trigger({ message: form.message, background: 'variant-filled-error' });
	}

	export let data: PageData;
	let hashMap = new Map<number, string>();
	hashMap.set(0, 'Inactive');
	hashMap.set(1, 'Waiting activation');
	hashMap.set(2, 'Active');
	hashMap.set(3, 'Waiting inactivation');
</script>

{#if data.api_list}
	<div class="table-container flex flex-col items-center">
		<h2 class="my-5">Your apis</h2>
		<table class="table table-hover w-6/12">
			<thead>
				<tr>
					<th>Name <a href="/apis/new" class="btn btn-sm variant-filled-success">+</a></th>
					<th>Exchange</th>
					<th class="text-center">Status</th>
					<!-- <th class="text-center "> <span> Actions </span> </th> -->
				</tr>
			</thead>
			<tbody>
				{#each data.api_list as row (row.id)}
					<tr>
						<td class="justify-center">
							<!-- <a class="btn btn-sm variant-filled-error" href="/">del</a>
							<a class="btn btn-sm variant-filled-warning" href="/">edit</a> -->
							{row.api_key_name}
						</td>
						<td class="w-32">{row.exchange}</td>
						<td class="text-center w-40">
							<form method="POST" action="?/toggleStatus" use:enhance>
								<input type="hidden" name="id" value={row.id} />
								<button
									type="submit"
									class="btn chip w-full variant-filled-primary"
									class:variant-filled-error={row.status == 0}
									class:variant-filled-success={row.status == 2}
									class:variant-filled-warning={row.status == 1}
									class:variant-filled-primary={row.status == 3}
								>
									{hashMap.get(row.status)}
								</button>
							</form>
						</td>
						<!-- <td class="w-20">
							<div class="btn-group variant-filled-warning">
								<a class="btn-sm"
									href="/apis/edit/{row.id}"
									>
									Edit
								</a>
								<a
									href="/apis/edit/{row.id}">
									Delete
								</a>
							</div>
						</td> -->
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
