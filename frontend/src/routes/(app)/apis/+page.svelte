<script lang="ts">
	import trash_icon from '../../../assets/delete_icon.svg';
	import edit_icon from '../../../assets/edit_icon.svg';
	import type { ApiListResp } from './+page.server';
	export let data;
	let hashMap = new Map<number, string>();
	hashMap.set(0, 'Inactive');
	hashMap.set(1, 'Waiting activation');
	hashMap.set(2, 'Active');
	hashMap.set(3, 'Waiting inactivation');

	function toggleStatus(status: ApiListResp, idx: number) {
		if (!data.api_list) return;
		let old=data.api_list[idx].status;
		data.api_list[idx].status = old == 0 ? 1 : (old == 1 ? 0 : (old == 2 ? 3 : 2))
	}
</script>

{#if data.api_list}
	<div class="table-container flex flex-col items-center">
		<h2 class="my-5">Your apis</h2>
		<!-- Native Table Element -->
		<table class="table table-hover w-6/12">
			<thead>
				<tr>
					<th>Name</th>
					<th>Exchange</th>
					<th class="text-center">Status</th>
					<th class="text-center">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each data.api_list as row, i}
					<tr>
						<td class="justify-center">{row.name}</td>
						<td class="table-cell-fit">{row.exchange}</td>
						<td class="table-cell-fit text-center ali">
							<button
								class="btn chip w-full"
								on:click={() => toggleStatus(row, i)}
								class:variant-filled-error={row.status == 0}
								class:variant-filled-success={row.status == 2}
								class:variant-filled-warning={row.status == 1 || row.status == 3}
							>
								{hashMap.get(row.status)}
							</button>
						</td>
						<td class="table-cell-fit">
							<div class="flex flex-row gap-x-2">
								<a
									href="/apis/edit/{row.id}"
									class="btn-icon btn-sm variant-filled-warning shrink-0"
									><img src={edit_icon} class="w-5 green-800" /></a
								>
								<a
									href="/apis/edit/{row.id}"
									class="btn-icon btn-sm variant-filled-warning shrink-0"
									><img src={trash_icon} class="w-3 green-800" /></a
								>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
