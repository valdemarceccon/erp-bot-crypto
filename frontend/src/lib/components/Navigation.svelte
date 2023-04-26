<script lang="ts" context="module">
	export type NavigationItem = {label: string, href: string}
</script>

<script lang="ts">
	import { page } from "$app/stores";

	import { drawerStore } from "@skeletonlabs/skeleton";

  export let items: NavigationItem[] = [];
  function drawerClose(): void {
    drawerStore.close();
  }
  $page.url.pathname
</script>

<nav class="list-nav p-4">
  <ul>
    {#each items as item}
      <li>
        <a
          href={item.href}
          class:variant-soft={$page.url.pathname==item.href}
          on:click={drawerClose}
        >
          {item.label}
        </a>
      </li>
    {/each}
  </ul>
</nav>
