<script lang="ts">
	export let showContent: boolean = true;
	export let hiddenMessage: string = 'hidden';
  const toggle = () => {
		showContent = !showContent;
	};
</script>


	{#if showContent}
  <button on:click={toggle} class="btn btn-sm variant-ghost-primary mr-10">Hide</button><slot />
	{:else}
  <button class="btn btn-sm variant-ghost-primary w-full"
	on:click={toggle}
>
		{hiddenMessage}
  </button>
	{/if}
