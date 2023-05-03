<script lang="ts">
	import type { PageData } from './$types';
	import type { UserBasicInfo } from './+page.server';

	export let data: PageData;

	let users: UserBasicInfo[] = [];

	if (data) {
		users = data.user_list ? data.user_list : [];
	}
</script>

<div>
	<div class="table-container flex flex-col items-center">
		<h2 class="my-5">Users</h2>
		<table class="table table-hover">
			<thead>
				<tr>
					<th class="text-center">Name</th>
					<th class="text-center">Username</th>
					<th class="text-center">Email</th>
				</tr>
			</thead>
			<tbody>
				{#each users as row, i}
					<tr>
						<td>{row.fullname}</td>
						<td>{row.username}</td>
						<td>{row.email}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
