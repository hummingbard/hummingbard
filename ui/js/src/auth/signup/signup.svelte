<script>
import {tick, onMount} from 'svelte'
import { debounce } from '../../utils/utils.js'

onMount(() => {
  if(signupError) {
    passwordInput.focus()
  }
  focusPassword()
})

async function focusPassword() {
  await tick();
  usernameInput.focus()
}

let federated = false;

let emailInput;
let usernameInput;
let passwordInput;
let repeatPasswordInput;

let loggingIn;

function signUp(e) {
  e.preventDefault()
  if(usernameInput.value.length < 3) {
    alert("That username is too short.")
    usernameInput.focus()
    return
  }
  if(passwordInput.value.length == 0) {
    passwordInput.focus()
    return
  }
  if(passwordInput.value != repeatPasswordInput.value) {
    alert("Passwords do not match.")
    return
  }
  loggingIn = true
  form.submit()
}

let form;

function toggleFederated() {
  federated = !federated
  focusCaret(usernameInput)
}

$: usernamePlaceholder = federated ? `@username:server.org` : `username`
$: tooltip = !federated ? `Sign up on` : `Log in with a Hummingbard account`

$: home = window.homeServer || ""


let username = '';
let checking = false;
let available = false;
let usernameAvailable = true;

function updateUsername(e) {
  const letters = /^[0-9a-zA-Z-]+$/;
  if(!e.key.match(letters)){
    e.preventDefault()
  }
  usernameAvailable = true
  available = false
  if(username.length === 0) {
      checking = false
      return
  }
}

function resetUsername() {
  usernameInput.value = ''
  username = ''
  available = false;
  usernameAvailable = true;
  usernameInput.focus()
}


async function checkUsername() {
    let endpoint = `/username/available`
    let data = {
        username: username,
    };
    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}
function reset() {
  available = false
  debounce(() =>{
    checking = true
    if(usernameInput.value.length < 3) {
      checking = false
      usernameAvailable = true
      available = false
      return
    }
    checkUsername().then((res) => {
      console.log(res)
        if(res?.available) {
          usernameAvailable = true
          available = true
        } else if(!res?.available) {
          usernameAvailable = false
          available = false
        }
        checking = false
    }).then(() => {
    })

  }, 500, this)
}

</script>


<div class="gr-center brd w-350px pa3">

  <form class="flex flex-column" 
  action="/signup" 
  bind:this={form}
 method="POST" 
 autocomplete="off">
      <div class="flex flex-column">
      {#if signupDisabled}
        <div class="mb3 warn">
          Signing up is disabled.
        </div>
      {/if}

            <input class="" type="text"
            bind:this={emailInput} 
            name="email"
            value={signupEmail} readonly>

            <input name="federated" type=checkbox bind:checked={federated}
                   disabled={signupDisabled} hidden>
      </div>
      <div class="mt3 flex flex-column relative">

        <input
          class:oops={!usernameAvailable}
          name="username"
          placeholder="username"
          on:keypress={updateUsername}
          on:input={reset}
          bind:this={usernameInput}
          bind:value={username}
        />
        {#if checking}
          <div class="checking mh2 gr-default">
            <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
          </div>
        {/if}
        {#if !usernameAvailable}
          <div class="checking mh2 gr-default pointer" on:click={resetUsername}>
            <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
          </div>
        {/if}
        {#if available}
          <div class="checking mh2 gr-default">
            <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M21.03 5.72a.75.75 0 010 1.06l-11.5 11.5a.75.75 0 01-1.072-.012l-5.5-5.75a.75.75 0 111.084-1.036l4.97 5.195L19.97 5.72a.75.75 0 011.06 0z"></path></svg>
          </div>
        {/if}
      </div>
      <div class="mt3 flex flex-column">
          <input class="" type="password"
            bind:this={passwordInput}
          name="password" minlength="8"  disabled={signupDisabled} placeholder="password">
      </div>
      <div class="mt3 flex flex-column">
          <input class="" type="password"
            bind:this={repeatPasswordInput}
          name="repeat" minlength="8"  disabled={signupDisabled} placeholder="repeat password">
      </div>

      {#if signupError}
        <div class="mt3 warn">
          Account could not be created.
        </div>
      {/if}


    <div class="mt3 flex">
      <div class="">
        <button class="dark-button-small no-sel" 
          disabled={loggingIn || signupDisabled}
          on:click={signUp}
          type="submit" >
          Create Account
        </button>
      </div>
      {#if loggingIn}
      <div class="ml3 gr-center">
        <div class="lds-ring"><div></div><div></div><div></div><div></div></div>
      </div>
      {/if}
    </div>
  </form>
</div>

<style>
.oops {
  border: 1px solid red;
}

.checking {
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
}

</style>
