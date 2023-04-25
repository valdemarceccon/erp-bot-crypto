<script lang="ts">
  export let active: "login" | "register" | null = null;
  export let login_label: string = "Login";
  export let register_label: string = "Register";


  let buttons = [
    {
      href: "/login",
      label: login_label,
      active: active == "login"
    },
    {
      href: "/register",
      label: register_label,
      active: active == "register"
    }
  ];
</script>

<div class="flex">
  {#each buttons as button (button.href) }
  <a href={button.href}
    class="btn"
    class:!bg-transparent={!button.active}
    class:variant-soft={button.active}
  >
    {button.label}
  </a>
  {/each}
</div>
