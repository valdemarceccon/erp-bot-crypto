<script lang="ts" context="module">
	export type NavigationItem = {label: string, href: string, permission?: string}
</script>

<script lang="ts">
	import { page } from "$app/stores";
  import { Accordion, AccordionItem } from '@skeletonlabs/skeleton';

	import { drawerStore } from "@skeletonlabs/skeleton";

  export let items: NavigationItem[] = [];
  export let title: string = "";
  function drawerClose(): void {
    drawerStore.close();
  }



</script>

<nav class="list-nav p-4">
  {#if title}
  <span class="text-primary-700 font-bold">{title}</span>
  {/if}

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
