<script>
let active = false;
let fetching = false;

let secret;
let image;

async function fetchMFA() {
    let endpoint = `/settings/mfa`
    let resp = await fetch(endpoint, {
        method: 'GET', // or 'PUT'
        headers:{
            'Authorization': identity.access_token,
            'Content-Type': 'application/json'
        }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}


function newMFA() {
    fetching = true
    fetchMFA().then(res => {
        console.log(res)
        if(res?.secret) {
            fetching = false
            image = res.image
            secret = res.secret
            active = true
        }
    }).then(()=> {
    })
}


function kill() {
    secret = null;
    image = null;
}

</script>

<div class="brd-btm pa3">
    <div class="">
        <span class="small"><strong>Two-Factor Authentication</strong></span>
    </div>
    <div class="mt3 flex flex-column">


      <div class="flex">
          Secure your Hummingbard account with an extra layer of authentication.
      </div>

      {#if !active}
      <div class="flex mt3">
          <button on:click={newMFA} disabled={fetching}>Enable Two-Factor Authentication</button>
      </div>
      {/if}

      {#if active}
          <div class="mt3 flex flex-column">
              {secret}
          </div>
      {/if}

    </div>
</div>
