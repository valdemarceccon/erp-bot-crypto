<script lang="ts">
	import { toastStore } from '@skeletonlabs/skeleton';
	import type { PageData } from './$types';

	export let data: PageData;

	$: if (!data.success && data.message) {
		toastStore.trigger({ message: data.message });
	}
</script>

<div class="table-container flex flex-col items-center">
	<h2 class="my-5">Commissions</h2>
	{#if data.commissions}
		{#each data.commissions as botRun (botRun.start)}
			<p>
				Bot ran from <b>{botRun.start.toUTCString()}</b>
				{#if botRun.stop} to <b>{botRun.stop.toUTCString()}</b>{/if}
				for user <b>{botRun.username}</b>
			</p>
			<table class="table table-hover w-full">
				<thead>
					<tr>
						<th class="text-center">Date</th>
						<th class="text-center">balance</th>
						<th class="text-center">High Watermark</th>
						<th class="text-center">Fee</th>
						<th class="text-center">Profit</th>
						<!-- <th class="text-center "> <span> Actions </span> </th> -->
					</tr>
				</thead>
				<tbody>
					{#each botRun.commissions as commission (commission.date)}
						<tr>
							<td class="w-32">{commission.date.toUTCString()}</td>
							<td class="w-32">{commission.balance}</td>
							<td class="w-32">{commission.high_mark}</td>
							<td class="w-32">{commission.fee}</td>
							<td class="w-32">{commission.profit}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/each}
	{/if}
</div>
