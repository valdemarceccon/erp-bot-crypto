<script lang="ts">
	import type { PageData } from "./$types";
	import type { UserBasicInfo } from "./+page.server";

	export let data: PageData;

	let users: UserBasicInfo[] = [];

	if (data) {
		users = data.user_list ? data.user_list : [];
	}

</script>

<div>


<div class="table-container">
	<h2>Users</h2>
	<table class="table table-hover">
		<thead>
			<tr>
				<th>Name</th>
				<th>Username</th>
				<th>Email</th>
			</tr>
		</thead>
		<tbody>
			{#each users as row, i}
				<tr>
					<td>{row.name}</td>
					<td>{row.username}</td>
					<td>{row.email}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

</div>
