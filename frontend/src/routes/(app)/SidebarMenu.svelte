<script lang="ts">
	import Navigation, { type NavigationItem } from '$lib/components/Navigation.svelte';

	export let noSpecialPermitionRequired: NavigationItem[] = [];
	export let specialPermitionRequired: NavigationItem[] = [];
</script>

{#if noSpecialPermitionRequired.length > 0}
	<Navigation items={noSpecialPermitionRequired} title="Menu" />
{/if}
{#if specialPermitionRequired.length > 0}
	<hr class="!border-t-2" />
	<Navigation items={specialPermitionRequired} title="Administration" />
{/if}
