<div class="card">
<form action="/auth/login" method="POST">
	<header class="card-header flex">
    Login
  </header>
	<section class="p-4">
		<label class="label">
			<span>Username</span>
			<input required name="username" class="input" type="text" placeholder="Username" />
		</label>
		<label class="label">
			<span>Password</span>
			<input type="password" required name="password" class="input" placeholder="Password" />
		</label>
	</section>
	<footer class="card-footer flex">
		<button type="submit" class="btn variant-filled-primary flex-auto">Login</button>
	</footer>
</form>
</div>
