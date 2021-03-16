<script>
let active = false;

let email = identity?.email

$: emailText = email?.length > 0 ? email : 'No Email Set'


function activate() {
    active = true
}

let emailExists = false;

function reset() {
    emailExists = false
}

async function saveEmail() {
    let endpoint = `/account/email/update`
    let data = {
      email: emailInput.value,
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

let updated = false;

function save() {
    if(emailInput.value == email) {
        cancel()
        return
    }
  saveEmail().then((res) => {
    console.log(res)
    if(res?.updated) {
        email = emailInput.value
        updated = true
        cancel()
    }
    if(res?.exists) {
        emailExists = true
    }
  }).then(() => {
  })
}

function cancel() {
    active = false
    emailInput.value = null
}

let emailInput;


$: verified = identity.email_verified

async function resendEmail() {
    let endpoint = `/account/email/resend`
    let data = {
      email: email,
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

let resent = false;

function resend() {
  resendEmail().then((res) => {
    console.log(res)
    if(res?.resent) {
        resent = true
        cancel()
    }
  }).then(() => {
  })
}

</script>

<div class="brd-btm pa3">
    <div class="">
        <span class="small"><strong>Email</strong></span>
    </div>
    <div class="mt3 flex flex-column">


      {#if !active}
      <div class="flex">
          {emailText}
      </div>

          {#if !verified && !updated}
              <div class="mt2 fs-09">
                  Your email has not been verified yet. 
                  <span class="pointer primary" on:click={resend}><u>Resend verification email?</u></span>
              </div>
          {/if}

          {#if !updated}
          <div class="flex mt3">
              <button on:click={activate}>Update Email</button>
          </div>
          {:else}
          <div class="flex mt3 fs-09 success">
              We've sent your a verification email. 
          </div>
          {/if}

          {#if resent}
          <div class="flex mt3 fs-09 success">
              We've sent your a verification email. 
          </div>
          {/if}

      {/if}


      {#if active}
      <div class="flex">
          <input class="" type="email" autocomplete="off"
          bind:this={emailInput}
       on:keypress={reset}
          value={email}
          name="email" placeholder="bob@peptides.com" required autofocus>
      </div>

      {#if emailExists}
          <div class="warn fs-09 mt2">
              That email is associated with another Hummingbard account.
          </div>
      {/if}

      <div class="flex mt3">
          <button class="mr2" on:click={save}>Save</button>
          <button class="light" on:click={cancel}>Cancel</button>
      </div>
      {/if}

    </div>
</div>

<style>
input {
    width: 100%;
}
</style>
