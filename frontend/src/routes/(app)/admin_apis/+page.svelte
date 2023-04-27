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
		<table class="table table-hover w-full lg:w-6/12 md:w-10/12 ">
			<thead>
				<tr>
					<th>Name</th>
					<th>Exchange</th>
					<th>Username</th>
					<th class="w-52">Api Key</th>
					<th class="w-full">Api Secret</th>
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
						<td class="w-32">{row.user.username}</td>
						<td class="w-32">{row.api_key}</td>
						<td class="w-32">
							<ToggleBtn showContent={false}>
								{row.secret}
							</ToggleBtn>
						</td>
						<td class="text-center w-40">
							<form method="POST" action="?/toggleStatus" use:enhance>
								<input type="hidden" name="id" value={row.id} />
								<input type="hidden" name="client_id" value={row.user.id} />
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
