<script>
let nsfwInput;
let darkmodeInput;

async function updatePreferences() {
    let endpoint = `/account/preferences/update`
    let data = {
      dark_mode: darkmodeInput.checked,
      show_nsfw_posts: nsfwInput.checked,
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

function update(){
  console.log(nsfwInput.checked)
  updatePreferences().then((res) => {
    console.log(res)
  }).then(() => {
    if(darkmodeInput?.checked) {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }
  })
}


$: nsfwActive = identity?.preferences?.show_nsfw_posts
$: darkmodeActive = identity?.preferences?.dark_mode

</script>

<div class="brd-btm pa3">
    <div class="">
        <span class="small"><strong>Preferences</strong></span>
    </div>
    <div class="mt3 flex flex-column">


      <div class="flex">
        <div class="gr-center">
          <input class="mr3" 
           type="checkbox" 
           id="darkmode" 
           checked={darkmodeActive}
           bind:this={darkmodeInput}
           on:change={update}/> 
        </div>
        <div class="gr-center">
          <label class="pointer fs-09 no-select" for="darkmode">Dark Mode </label>
        </div>
      </div>

      <div class="flex mt2">
        <div class="gr-center">
          <input class="mr3" 
           type="checkbox" 
           id="nsfw" 
           checked={nsfwActive}
           bind:this={nsfwInput}
           on:change={update}/> 
        </div>
        <div class="gr-center">
          <label class="pointer fs-09 no-select" for="nsfw">Show NSFW posts</label>
        </div>
      </div>


    </div>
</div>
