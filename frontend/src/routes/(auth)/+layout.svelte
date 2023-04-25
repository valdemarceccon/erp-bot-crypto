<script lang="ts">
  import LogoSVG from '../../assets/logoipsum-256.svg'
</script>

<div class="sm:w-10/12 md:w-6/12 mx-auto mt-10 space-y-10">
	<div class="card p-4 bg-black items-start flex flex-row justify-center">
    <img class="bg-blue-400 bg-opacity-50 rounded-lg w-2/12 p-5" src={LogoSVG} alt="Company banner">
  </div>
  <div class="overflow-auto">
    <slot />
  </div>
</div>

<style lang="scss">
  :global(body) {
    overflow: auto;
  }
</style>
