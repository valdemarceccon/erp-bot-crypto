let backend = process.env.BACKEND_PRIVATE_HOST;
if (!backend) {
  backend = "";
}

type ApiEndpoints = "/auth/token" | "/users/me"
type PrefixedApiEndpoints<T extends string> = `${T}${ApiEndpoints}`
type PublicApiEndpoints = PrefixedApiEndpoints<"/api">
