<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData } from './$types';
	import ToggleBtn from './ToggleBtn.svelte';

	export let data: PageData;
	let hashMap = new Map<number, string>();
	hashMap.set(0, 'Inactive');
	hashMap.set(1, 'Waiting activation');
	hashMap.set(2, 'Active');
	hashMap.set(3, 'Waiting inactivation');
</script>

{#if data.api_list}
	<div class="table-container flex flex-col items-center">
		<h2 class="my-5">All apis</h2>
		<table class="table table-hover w-full">
			<thead>
				<tr>
					<th class="text-center">Name</th>
					<th class="text-center">Exchange</th>
					<th class="text-center">Username</th>
					<th class="w-52 text-center">Api Key</th>
					<th class="w-full text-center">Api Secret</th>
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
							{row.name}
						</td>
						<td class="w-32">{row.exchange}</td>
						<td class="w-32">{row.username}</td>
						<td class="w-32">{row.api_key}</td>
						<td class="w-full">
							<ToggleBtn showContent={false}>
								{row.api_secret}
							</ToggleBtn>
						</td>
						<td class="text-center w-40">
							<form method="POST" action="?/toggleStatus" use:enhance>
								<input type="hidden" name="id" value={row.id} />
								<input type="hidden" name="client_id" value={row.user_id} />
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
