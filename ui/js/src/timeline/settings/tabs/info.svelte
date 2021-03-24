<script>
import Avatar from '../../../avatar/avatar.svelte'
import { settings,state } from '../store.js'


function updateAvatar(e) {
  $settings.info.avatar = e.detail
    setAvatar()
}

function removeAvatar(e) {
  $settings.info.avatar = null
    setAvatar()
}

async function saveState() {
    let endpoint = `/room/avatar/update`

    let data = {
      room_id: window.timeline.room_id,
      profile: window.timeline.profile,
      avatar: $settings.info.avatar,
    }
  console.log(data)

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: JSON.stringify(data),
    headers:{
      'Authorization': identity.access_token,
        'Content-Type': 'application/json'
    }
    })
    if (resp.ok) { // if HTTP-status is 200-299
    } else {
      alert("HTTP-Error: " + resp.status);
    }
    const ret = await resp.json()
    return Promise.resolve(ret)
}

function setAvatar() {
  saveState().then((res) => {
    console.log(res)
  }).then(() => {
  })
}


$: about = window.timeline?.profile ? `About you` : `About this community`

</script>


<div class="flex flex-column flex-one ">
  <div class="flex flex-column ">
    <div class="">
      <span class="small"><strong>Title</strong></span>
    </div>
    <div class="mt2">
      <input
        placeholder="Title"
        bind:value={$settings.info.title}
      />
    </div>
  </div>


  <div class="mt3 flex flex-column">
    <div class="">
      <span class="small"><strong>About</strong></span>
    </div>
    <div class="mt2">
      <textarea
        style="height:100px;"
        placeholder={about}
        bind:value={$settings.info.about}
      ></textarea>
    </div>
  </div>
</div>

<div class="ml4 flex flex-column">

  <div class="">
    <Avatar 
      on:uploaded={updateAvatar}
      on:removed={removeAvatar}
      avatar={$settings.info.avatar}/>
  </div>

  <div class="">
  </div>

</div>

<style>

input, textarea {
  width: 100%;
}

</style>
