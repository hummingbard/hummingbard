<script>
import {onMount} from 'svelte'
import { focusCaret } from '../../utils/utils.js'

onMount(() => {
  if(loginFederated) {
    federated = true
  }
  if(loginError) {
    passwordInput.focus()
  }
  focusCaret(usernameInput)
})
let federated = false;

let usernameInput;
let passwordInput;

let loggingIn;

function logIn() {
  loggingIn = true
  form.submit()
}

let form;

function toggleFederated() {
  federated = !federated
  focusCaret(usernameInput)
}

$: usernamePlaceholder = federated ? `@username:homeserver.org` : `username`
$: tooltip = !federated ? `Log in with an external matrix account` : `Log in with a Hummingbard account`

</script>


<div class="gr-center">

  <form class="flex flex-column" 
  action="/login" 
  bind:this={form}
 method="POST" 
 autocomplete="off">
      <div class="flex flex-column">
        {#if federated}
        <div class="mb3 fs-09 lh-copy">
          Log in with an external <br/> matrix account.
        </div>
        {/if}
        {#if passwordReset}
        <div class="mb3 fs-09 lh-copy">
          Your password was reset successfully. You can log in now.
        </div>
        {/if}
            <input class="" type="text" autofocus="autofocus"
            bind:this={usernameInput}
            value={loginUsername}
            name="username" placeholder={usernamePlaceholder}>
            <input name="federated" type=checkbox bind:checked={federated} hidden>
      </div>
      <div class="mt3 flex flex-column">
          <input class="" type="password"
            bind:this={passwordInput}
          name="password" minlength="8" placeholder="password">
      </div>

      {#if loginError}
              <div class="mt3 warn">
                Username or password wrong
              </div>
      {/if}

      {#if federationDisabledWarning}
        <div class="mt3 warn lh-copy">
            Federation has been disabled. Please use a local account.
        </div>
      {/if}

    <div class="mt3 flex">
      <div class="">
        <button class="dark-button-small no-sel" 
          disabled={loggingIn}
          on:click={logIn}
          type="submit" >
          Log In
        </button>
      </div>
      {#if loggingIn}
      <div class="ml3 gr-center">
        <div class="lds-ring"><div></div><div></div><div></div><div></div></div>
      </div>
      {/if}
      <div class="flex-one">
      </div>
      {#if !federationDisabled}
      <div class="gr-center">
        <div class="gr-center pointer" class:m-log={!federated}
             title={tooltip}
          on:click={toggleFederated}>
          <svg class:fill-blue={federated} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 22 22" width="22" height="22">
          <path id="Layer" fill="" d="M0.58 0.51L0.58 21.5L2.09 21.5L2.09 22L0 22L0 0L2.09 0L2.09 0.51L0.58 0.51ZM7.04 7.16L7.04 8.22L7.06 8.22C7.35 7.81 7.69 7.5 8.09 7.28C8.49 7.06 8.95 6.95 9.46 6.95C9.96 6.95 10.41 7.04 10.82 7.23C11.23 7.43 11.54 7.77 11.76 8.25C11.99 7.91 12.31 7.6 12.7 7.34C13.1 7.08 13.58 6.95 14.12 6.95C14.54 6.95 14.92 7 15.28 7.1C15.63 7.2 15.93 7.36 16.19 7.58C16.44 7.81 16.64 8.1 16.78 8.46C16.92 8.82 16.99 9.25 16.99 9.76L16.99 15.01L14.84 15.01L14.84 10.56C14.84 10.3 14.83 10.05 14.81 9.82C14.79 9.6 14.74 9.4 14.64 9.21C14.55 9.04 14.41 8.89 14.24 8.8C14.07 8.7 13.82 8.65 13.52 8.65C13.22 8.65 12.97 8.71 12.79 8.82C12.6 8.93 12.45 9.09 12.35 9.28C12.24 9.48 12.16 9.69 12.14 9.92C12.1 10.16 12.08 10.4 12.08 10.64L12.08 15.01L9.93 15.01L9.93 10.61C9.93 10.37 9.92 10.15 9.91 9.92C9.9 9.7 9.86 9.49 9.78 9.29C9.7 9.1 9.57 8.94 9.4 8.83C9.22 8.71 8.96 8.65 8.62 8.65C8.51 8.65 8.38 8.67 8.21 8.72C8.05 8.77 7.88 8.85 7.73 8.98C7.57 9.11 7.43 9.29 7.32 9.53C7.21 9.76 7.16 10.07 7.16 10.46L7.16 15.01L5.01 15.01L5.01 7.16L7.04 7.16ZM21.42 21.49L21.42 0.5L19.91 0.5L19.91 0L22 0L22 22L19.91 22L19.91 21.49L21.42 21.49Z" />
          </svg>
        </div>
      </div>
      {/if}
    </div>
  </form>

</div>

