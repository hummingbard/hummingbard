<script>
import {tick} from 'svelte'
let active = false;
let error = false;

let success = false;

function activate() {
  success = false
  active = true
  focusPass()
}

function kill() {
  active = false
  error = false
  oldPassInput.value = null
  newPassInput.value = null
  repeatInput.value = null
}

async function focusPass() {
  await tick();
  oldPassInput.focus()
}

let oldPassInput;
let newPassInput;
let repeatInput;

async function update() {
  error = false
  if(oldPassInput.value.length < 8) {
    focusPass()
    return
  }
  if(newPassInput.value.length < 8) {
    await tick();
    newPassInput.focus();
    return
  }
  if(repeatInput.value.length < 8) {
    await tick();
    repeatInput.focus();
    return
  }
  if(repeatInput.value != newPassInput.value) {
    alert('Passwords do not match.')
    return
  }
  if(repeatInput.value == oldPassInput.value) {
    alert('New password must be different from current password.')
    return
  }
  confirmUpdate().then((res) => {
    console.log(res)
    if(res?.error) {
      error = true
    }
    if(res?.updated) {
      success = true
      kill()
    }
  }).then(() => {
  })
}

async function confirmUpdate() {
    let endpoint = `/account/password/update`
    let data = {
      old_password: oldPassInput.value,
      new_password: newPassInput.value,
      repeat_password: repeatInput.value,
    };
    let resp = await fetch(endpoint, {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Authorization': identity.access_token,
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


</script>

<div class="brd-btm pa3">
    <div class="">
        <span class="small"><strong>Password</strong></span>
    </div>

    {#if success}
      <div class="mt3">
        <span class="fs-09">Your password was changed successfully.</span>
      </div>
    {/if}

    {#if active}

    <form class="mt3 flex flex-column">

      <div class="">
          <input class="" type="password" autocomplete="off"
          bind:this={oldPassInput}
          name="password" minlength="8" placeholder="current password">
      </div>

      <div class="mt2">
          <input class="" type="password" autocomplete="off"
          bind:this={newPassInput}
          name="password" minlength="8" placeholder="new password">
      </div>

      <div class="mt2">
          <input class="" type="password" autocomplete="off"
          bind:this={repeatInput}
          name="password" minlength="8" placeholder="repeat new password">
      </div>

    </form>

    {/if}

    <div class="mt3 flex flex-column">

      {#if !active}
        <div class="flex">
            <button on:click={activate}>Change Password</button>
        </div>
      {:else}
        <div class="flex">
            <button class="mr2" on:click={update}>Change Password</button>
            <button class="light" on:click={kill}>Cancel</button>
        </div>
      {/if}

    </div>

    {#if error}
      <div class="mt3 flex flex-column">
        <span class="warn fs-09">Your password could not be changed. Are you sure your current password
          is correct?</span>
      </div>

    {/if}

</div>
