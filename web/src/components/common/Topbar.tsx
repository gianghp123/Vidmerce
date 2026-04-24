import { Show, UserButton, SignInButton } from "@clerk/nextjs";

export function Topbar() {
  return (
    <header style={{ display: 'flex', justifyContent: 'flex-end', padding: 20 }}>
      {/* <h1>My App</h1> */}
      <Show when="signed-in">
        <UserButton />
      </Show>
      <Show when="signed-out">
        <SignInButton />
      </Show>
    </header>
  )
}